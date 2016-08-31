package utils

import (
	"errors"
	"github.com/codegangsta/inject"
)

type TypeInputConfig interface {
	TypeConfig
	Start()
}

type InputConfig struct {
	CommonConfig
}

type InputHandler interface{}

var (
	mapInputHandler = map[string]InputHandler{}
)

func RegistInputHandler(name string, handler InputHandler) {
	mapInputHandler[name] = handler
}

func (t *Config) RunInputs() (err error) {
	_, err = t.Injector.Invoke(t.runInputs)
	return
}

func (t *Config) runInputs(inchan InChan) (err error) {
	inputs, err := t.getInputs(inchan)
	if err != nil {
		return
	}
	for _, input := range inputs {
		go input.Start()
	}
	return
}

func (t *Config) getInputs(inchan InChan) (inputs []TypeInputConfig, err error) {
	for _, confraw := range t.InputRaw {
		handler, ok := mapInputHandler[confraw["type"].(string)]
		if !ok {
			err = errors.New(confraw["type"].(string))
			return
		}

		inj := inject.New()
		inj.SetParent(t)
		inj.Map(&confraw)
		inj.Map(inchan)

		refvs, err := inj.Invoke(handler)
		if err != nil {
			return []TypeInputConfig{}, err
		}

		for _, refv := range refvs {
			if !refv.CanInterface() {
				continue
			}
			if conf, ok := refv.Interface().(TypeInputConfig); ok {
				conf.SetInjector(inj)
				inputs = append(inputs, conf)
			}
		}
	}
	return
}
