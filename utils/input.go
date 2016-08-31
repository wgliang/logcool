package utils

import (
	"errors"
	"fmt"
	"github.com/codegangsta/inject"
)

// Input base type interface.
type TypeInputConfig interface {
	TypeConfig
	Start()
}

// Input base type struct.
type InputConfig struct {
	CommonConfig
}

// InputHandler type interface.
type InputHandler interface{}

var (
	mapInputHandler = map[string]InputHandler{}
)

// Registe InputHandler.
func RegistInputHandler(name string, handler InputHandler) {
	mapInputHandler[name] = handler
}

// Run Inputs.
func (t *Config) RunInputs() (err error) {
	fmt.Println("Input start...")
	_, err = t.Injector.Invoke(t.runInputs)
	fmt.Println("Input end...")
	return
}

// run Inputs.
func (t *Config) runInputs(inchan InChan) (err error) {
	fmt.Println(t)
	inputs, err := t.getInputs(inchan)

	if err != nil {
		fmt.Println("============")
		return
	}
	for _, input := range inputs {
		fmt.Println(input)
		go input.Start()
	}
	return
}

// get Inputs.
func (t *Config) getInputs(inchan InChan) (inputs []TypeInputConfig, err error) {

	for _, confraw := range t.InputRaw {
		fmt.Println(confraw)
		handler, ok := mapInputHandler[confraw["type"].(string)]
		fmt.Println("??????")
		fmt.Println(handler)
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
