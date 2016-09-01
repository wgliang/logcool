package fileinput

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type SinceDBInfo struct {
	Offset int64 `json:"offset"`
}

func (self *InputConfig) LoadSinceData() (err error) {
	var (
		raw []byte
	)
	log.Debug("LoadSinceDBInfos")
	self.SinceDBInfos = map[string]*SinceDBInfo{}

	if self.SincePath == "" || self.SincePath == "/dev/null" {
		log.Warnf("No valid sincedb path")
		return
	}
	if _, err = os.Stat(self.SincePath); err != nil {
		log.Debugf("sincedb not found: %q", self.SincePath)
		return
	}

	if raw, err = ioutil.ReadFile(self.SincePath); err != nil {
		log.Errorf("Read sincedb failed: %q\n%s", self.SincePath, err)
		return
	}

	if err = json.Unmarshal(raw, &self.SinceDBInfos); err != nil {
		log.Errorf("Unmarshal sincedb failed: %q\n%s", self.SincePath, err)
		return
	}

	return
}

func (self *InputConfig) SaveSinceDBInfos() (err error) {
	var (
		raw []byte
	)
	log.Debug("SaveSinceDBInfos")
	self.SinceLastSaveTime = time.Now()

	if self.SincePath == "" || self.SincePath == "/dev/null" {
		log.Warnf("No valid sincedb path")
		return
	}

	if raw, err = json.Marshal(self.SinceDBInfos); err != nil {
		log.Errorf("Marshal sincedb failed: %s", err)
		return
	}
	self.sinceLastInfos = raw

	if err = ioutil.WriteFile(self.SincePath, raw, 0664); err != nil {
		log.Errorf("Write sincedb failed: %q\n%s", self.SincePath, err)
		return
	}

	return
}

func (self *InputConfig) CheckSaveSinceDBInfos() (err error) {
	var (
		raw []byte
	)
	if time.Since(self.SinceLastSaveTime) > time.Duration(self.SinceInterval)*time.Second {
		if raw, err = json.Marshal(self.SinceDBInfos); err != nil {
			log.Errorf("Marshal sincedb failed: %s", err)
			return
		}
		if bytes.Compare(raw, self.sinceLastInfos) != 0 {
			err = self.SaveSinceDBInfos()
		}
	}
	return
}

func (self *InputConfig) LoopCheckSaveSinceInfos() (err error) {
	for {
		time.Sleep(time.Duration(self.SinceInterval) * time.Second)
		if err = self.CheckSaveSinceDBInfos(); err != nil {
			return
		}
	}
	return
}
