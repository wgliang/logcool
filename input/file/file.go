// Input-plug: fileinput
// The plug's function is real-time monitoring of the specified file, once the data is
//updated to record the data.
package fileinput

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"

	"github.com/wgliang/logcool/utils"
)

const (
	ModuleName = "file"
)

var (
	mapWatcher = map[string]*fsnotify.Watcher{}
)

type SinceDBInfo struct {
	Offset int64 `json:"offset"`
}

// Define fileinput' config.
type InputConfig struct {
	utils.InputConfig
	DirsPath  []string `json:"dirspath"`
	FileType  string   `json:"filetype"`
	StartPos  string   `json:"startpos"`
	SincePath string   `json:"sincepath"`
	Intervals int      `json:"intervals"`

	hostname          string
	SinceDBInfos      map[string]*SinceDBInfo
	sinceLastInfos    []byte
	SinceLastSaveTime time.Time
}

func init() {
	utils.RegistInputHandler(ModuleName, InitHandler)
}

// Init fileinput Handler.
func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeInputConfig, err error) {
	conf := InputConfig{
		InputConfig: utils.InputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},

		SinceDBInfos: map[string]*SinceDBInfo{},
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}
	// get hostname.
	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}

	retconf = &conf
	return
}

// Input's start,and this is the main function of input.
func (ic *InputConfig) Start() {
	ic.Invoke(ic.monitor)
}

// load current since data.
func (ic *InputConfig) LoadSinceData() (err error) {
	var (
		raw []byte
	)
	log.Debug("LoadSinceDBInfos")
	ic.SinceDBInfos = map[string]*SinceDBInfo{}

	if ic.SincePath == "" || ic.SincePath == "/dev/null" {
		log.Warnf("No valid sincedb path")
		return
	}
	if _, err = os.Stat(ic.SincePath); err != nil {
		log.Debugf("sincedb not found: %q", ic.SincePath)
		return
	}

	if raw, err = ioutil.ReadFile(ic.SincePath); err != nil {
		log.Errorf("Read sincedb failed: %q\n%s", ic.SincePath, err)
		return
	}

	if err = json.Unmarshal(raw, &ic.SinceDBInfos); err != nil {
		log.Errorf("Unmarshal sincedb failed: %q\n%s", ic.SincePath, err)
		return
	}

	return
}

// save since data info.
func (ic *InputConfig) SaveSinceDBInfos() (err error) {
	var (
		raw []byte
	)
	log.Debug("SaveSinceDBInfos")
	ic.SinceLastSaveTime = time.Now()

	if ic.SincePath == "" || ic.SincePath == "/dev/null" {
		log.Warnf("No valid sincedb path")
		return
	}

	if raw, err = json.Marshal(ic.SinceDBInfos); err != nil {
		log.Errorf("Marshal sincedb failed: %s", err)
		return
	}
	ic.sinceLastInfos = raw

	if err = ioutil.WriteFile(ic.SincePath, raw, 0664); err != nil {
		log.Errorf("Write sincedb failed: %q\n%s", ic.SincePath, err)
		return
	}

	return
}

// check since data info.
func (ic *InputConfig) CheckSaveSinceDBInfos() (err error) {
	var (
		raw []byte
	)
	if time.Since(ic.SinceLastSaveTime) > time.Duration(ic.Intervals)*time.Second {
		if raw, err = json.Marshal(ic.SinceDBInfos); err != nil {
			log.Errorf("Marshal sincedb failed: %s", err)
			return
		}
		if bytes.Compare(raw, ic.sinceLastInfos) != 0 {
			err = ic.SaveSinceDBInfos()
		}
	}
	return
}

// load check save since data
func (ic *InputConfig) LoopCheckSaveSinceInfos() (err error) {
	for {
		time.Sleep(time.Duration(ic.Intervals) * time.Second)
		if err = ic.CheckSaveSinceDBInfos(); err != nil {
			return
		}
	}
}

func (ic *InputConfig) monitor(logger *logrus.Logger, inchan utils.InChan) (err error) {
	defer func() {
		if err != nil {
			logger.Errorln(err)
		}
	}()

	var (
		matches = make([]string, 0)
		fi      os.FileInfo
	)

	if err = ic.LoadSinceData(); err != nil {
		return
	}
	// load all files.
	if len(ic.DirsPath) < 1 {
		return
	}
	for _, dir := range ic.DirsPath {
		var matche []string
		matche, err = fileList(dir, ic.FileType)
		if err != nil {
			logger.Errorln(err)
		}
		matches = append(matches, matche...)
	}

	go ic.LoopCheckSaveSinceInfos()

	for _, fpath := range matches {
		// get all sysmlinks.
		if fpath, err = filepath.EvalSymlinks(fpath); err != nil {
			logger.Errorf("Get symlinks failed: %q\n%v", fpath, err)
			continue
		}
		// check file status.
		if fi, err = os.Stat(fpath); err != nil {
			logger.Errorf("Stat(%q) failed\n%s", ic.DirsPath, err)
			continue
		}
		// check file isDir?
		if fi.IsDir() {
			logger.Infof("Skipping directory: %q", ic.DirsPath)
			continue
		}
		// monitor file.
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
	// if no path's since, add a null SinceDB.
	if since, ok = ic.SinceDBInfos[fpath]; !ok {
		ic.SinceDBInfos[fpath] = &SinceDBInfo{}
		since = ic.SinceDBInfos[fpath]
	}
	// set or get offset index.
	if since.Offset == 0 {
		if ic.StartPos == "end" {
			whence = os.SEEK_END // SEEK_END = 2
		} else {
			whence = os.SEEK_SET // SEEK_SET = 0
		}
	} else {
		whence = os.SEEK_SET // SEEK_SET = 0
	}
	// open the file.
	if fp, reader, err = openFile(fpath, since.Offset, whence); err != nil {
		return
	}
	defer fp.Close()
	// seek beginning.
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
	// load all data.
	for {
		if line, size, err = readLine(reader, buffer); err != nil {
			if err == io.EOF {
				watchev := <-readEventChan
				logger.Debug("loopRead recv:", watchev)
				if watchev.Op&fsnotify.Create == fsnotify.Create {
					logger.Warnf("File recreated, seeking to beginning: %q", fpath)
					fp.Close()
					since.Offset = 0
					if fp, reader, err = openFile(fpath, since.Offset, os.SEEK_SET); err != nil {
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

		event := utils.LogEvent{
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

func (ic *InputConfig) loopWatch(readEventChan chan fsnotify.Event, fpath string, op fsnotify.Op) (err error) {
	var (
		event fsnotify.Event
	)
	for {
		if event, err = waitWatchEvent(fpath, op); err != nil {
			return
		}
		readEventChan <- event
	}
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

func openFile(fpath string, offset int64, whence int) (fp *os.File, reader *bufio.Reader, err error) {
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

		if len(segment) < 1 {
			time.Sleep(1 * time.Second)
		} else {
			if segment[len(segment)-1] != '\n' {
				time.Sleep(1 * time.Second)
			} else {
				size = buffer.Len()
				line = buffer.String()
				buffer.Reset()
				line = strings.TrimRight(line, "\r\n")
				return
			}

		}
	}
}

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
}

// list all config templates.
func fileList(dirPath string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPath+PthSep+fi.Name())
		}
	}

	return files, nil
}
