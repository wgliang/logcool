// Input-plug: fileinput
// The plug's function is real-time monitoring of the specified file, once the data is
//updated to record the data.
package fileinput

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/wgliang/logcool/utils"
	"github.com/wgliang/logcool/utils/logevent"
)

const (
	ModuleName = "file"
)

// Define fileinput' config.
type InputConfig struct {
	utils.InputConfig
	Path          string `json:"path"`
	StartPos      string `json:"start_position"`
	SincePath     string `json:"since_path"`
	SinceInterval int    `json:"since_interval"`

	hostname          string
	SinceDBInfos      map[string]*SinceDBInfo
	sinceLastInfos    []byte
	SinceLastSaveTime time.Time
}

// Init fileinput Handler.
func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeInputConfig, err error) {
	conf := InputConfig{
		InputConfig: utils.InputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
		StartPos:      "end",
		SincePath:     ".sincedb.json",
		SinceInterval: 15,

		SinceDBInfos: map[string]*SinceDBInfo{},
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}

	retconf = &conf
	return
}

// Input's start,and this is the main function of input.
func (ic *InputConfig) Start() {
	ic.Invoke(ic.start)
}

func (ic *InputConfig) start(logger *logrus.Logger, inchan utils.InChan) (err error) {
	defer func() {
		if err != nil {
			logger.Errorln(err)
		}
	}()

	var (
		matches []string
		fi      os.FileInfo
	)

	if err = ic.LoadSinceData(); err != nil {
		return
	}

	if matches, err = filepath.Glob(ic.Path); err != nil {
		return err
	}

	go ic.LoopCheckSaveSinceInfos()

	for _, fpath := range matches {
		if fpath, err = filepath.EvalSymlinks(fpath); err != nil {
			logger.Errorf("Get symlinks failed: %q\n%v", fpath, err)
			continue
		}

		if fi, err = os.Stat(fpath); err != nil {
			logger.Errorf("stat(%q) failed\n%s", ic.Path, err)
			continue
		}

		if fi.IsDir() {
			logger.Infof("Skipping directory: %q", ic.Path)
			continue
		}

		readEventChan := make(chan fsnotify.Event, 10)
		go ic.loopRead(readEventChan, fpath, logger, inchan)
		go ic.loopWatch(readEventChan, fpath, fsnotify.Create|fsnotify.Write)
	}

	return
}

func (ic *InputConfig) loopRead(
	readEventChan chan fsnotify.Event,
	fpath string,
	logger *logrus.Logger,
	inchan utils.InChan,
) (err error) {
	var (
		since     *SinceDBInfo
		fp        *os.File
		truncated bool
		ok        bool
		whence    int
		reader    *bufio.Reader
		line      string
		size      int

		buffer = &bytes.Buffer{}
	)

	if fpath, err = filepath.EvalSymlinks(fpath); err != nil {
		logger.Errorf("Get symlinks failed: %q\n%v", fpath, err)
		return
	}

	if since, ok = ic.SinceDBInfos[fpath]; !ok {
		ic.SinceDBInfos[fpath] = &SinceDBInfo{}
		since = ic.SinceDBInfos[fpath]
	}

	if since.Offset == 0 {
		if ic.StartPos == "end" {
			whence = os.SEEK_END
		} else {
			whence = os.SEEK_SET
		}
	} else {
		whence = os.SEEK_SET
	}

	if fp, reader, err = openfile(fpath, since.Offset, whence); err != nil {
		return
	}
	defer fp.Close()

	if truncated, err = isTruncated(fp, since); err != nil {
		return
	}
	if truncated {
		logger.Warnf("File truncated, seeking to beginning: %q", fpath)
		since.Offset = 0
		if _, err = fp.Seek(since.Offset, os.SEEK_SET); err != nil {
			logger.Errorf("seek file failed: %q", fpath)
			return
		}
	}

	for {
		if line, size, err = readLine(reader, buffer); err != nil {
			if err == io.EOF {
				watchev := <-readEventChan
				logger.Debug("loopRead recv:", watchev)
				if watchev.Op&fsnotify.Create == fsnotify.Create {
					logger.Warnf("File recreated, seeking to beginning: %q", fpath)
					fp.Close()
					since.Offset = 0
					if fp, reader, err = openfile(fpath, since.Offset, os.SEEK_SET); err != nil {
						return
					}
				}
				if truncated, err = isTruncated(fp, since); err != nil {
					return
				}
				if truncated {
					logger.Warnf("File truncated, seeking to beginning: %q", fpath)
					since.Offset = 0
					if _, err = fp.Seek(since.Offset, os.SEEK_SET); err != nil {
						logger.Errorf("seek file failed: %q", fpath)
						return
					}
					continue
				}
				logger.Debugf("watch %q %q %v", watchev.Name, fpath, watchev)
				continue
			} else {
				return
			}
		}

		event := logevent.LogEvent{
			Timestamp: time.Now(),
			Message:   line,
			Extra: map[string]interface{}{
				"host":   ic.hostname,
				"path":   fpath,
				"offset": since.Offset,
			},
		}

		since.Offset += int64(size)

		logger.Debugf("%q %v", event.Message, event)
		inchan <- event
		ic.CheckSaveSinceDBInfos()
	}
}

func (self *InputConfig) loopWatch(readEventChan chan fsnotify.Event, fpath string, op fsnotify.Op) (err error) {
	var (
		event fsnotify.Event
	)
	for {
		if event, err = waitWatchEvent(fpath, op); err != nil {
			return
		}
		readEventChan <- event
	}
	return
}

func isTruncated(fp *os.File, since *SinceDBInfo) (truncated bool, err error) {
	var (
		fi os.FileInfo
	)
	if fi, err = fp.Stat(); err != nil {
		return
	}
	if fi.Size() < since.Offset {
		truncated = true
	} else {
		truncated = false
	}
	return
}

func openfile(fpath string, offset int64, whence int) (fp *os.File, reader *bufio.Reader, err error) {
	if fp, err = os.Open(fpath); err != nil {
		return
	}

	if _, err = fp.Seek(offset, whence); err != nil {
		err = errors.New("seek file failed: " + fpath)
		return
	}

	reader = bufio.NewReaderSize(fp, 16*1024)
	return
}

func readLine(reader *bufio.Reader, buffer *bytes.Buffer) (line string, size int, err error) {
	var (
		segment []byte
	)

	for {
		if segment, err = reader.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				err = errors.New("read line failed")
			}
			return
		}

		if _, err = buffer.Write(segment); err != nil {
			return
		}

		if isPartialLine(segment) {
			time.Sleep(1 * time.Second)
		} else {
			size = buffer.Len()
			line = buffer.String()
			buffer.Reset()
			line = strings.TrimRight(line, "\r\n")
			return
		}
	}

	return
}

func isPartialLine(segment []byte) bool {
	if len(segment) < 1 {
		return true
	}
	if segment[len(segment)-1] != '\n' {
		return true
	}
	return false
}

var (
	mapWatcher = map[string]*fsnotify.Watcher{}
)

func waitWatchEvent(fpath string, op fsnotify.Op) (event fsnotify.Event, err error) {
	var (
		fdir    string
		watcher *fsnotify.Watcher
		ok      bool
	)

	if fpath, err = filepath.EvalSymlinks(fpath); err != nil {
		err = errors.New("Get symlinks failed: " + fpath)
		return
	}

	fdir = filepath.Dir(fpath)

	if watcher, ok = mapWatcher[fdir]; !ok {
		//		logger.Debugf("create new watcher for %q", fdir)
		if watcher, err = fsnotify.NewWatcher(); err != nil {
			err = errors.New("create new watcher failed: " + fdir)
			return
		}
		mapWatcher[fdir] = watcher
		if err = watcher.Add(fdir); err != nil {
			err = errors.New("add new watch path failed: " + fdir)
			return
		}
	}

	for {
		select {
		case event = <-watcher.Events:
			if event.Name == fpath {
				if op > 0 {
					if event.Op&op > 0 {
						return
					}
				} else {
					return
				}
			}
		case err = <-watcher.Errors:
			err = errors.New("watcher error")
			return
		}
	}

	return
}
