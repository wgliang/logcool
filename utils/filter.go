package utils

import (
	"errors"
	"fmt"
	"github.com/codegangsta/inject"
	"logcool/utils/logevent"
)

// Filter base type interface.
type TypeFilterConfig interface {
	TypeConfig
	Event(logevent.LogEvent) logevent.LogEvent
}

// Filter base type struct.
type FilterConfig struct {
	CommonConfig
}

// FilterHandler type interface.
type FilterHandler interface{}

var (
	mapFilterHandler = map[string]FilterHandler{}
)

// Registe FilterHandler.
func RegistFilterHandler(name string, handler FilterHandler) {
	mapFilterHandler[name] = handler
}

// Run Filters
func (c *Config) RunFilters() (err error) {
	fmt.Println("Filter start...")
	_, err = c.Injector.Invoke(c.runFilters)
	fmt.Println("Filter end...")
	return
}

// run Filetrs.
func (c *Config) runFilters(inchan InChan, outchan OutChan) (err error) {
	fmt.Println("running filter...")
	filters, err := c.getFilters()
	fmt.Println(filters)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case event := <-inchan:
				for _, filter := range filters {
					event = filter.Event(event)
				}
				outchan <- event
			}
		}
	}()
	return
}

// get Filters.
func (c *Config) getFilters() (filters []TypeFilterConfig, err error) {
	fmt.Println("--------")
	for _, confraw := range c.FilterRaw {
		handler, ok := mapFilterHandler[confraw["type"].(string)]
		if !ok {
			err = errors.New(confraw["type"].(string))
			return
		}

		inj := inject.New()
		inj.SetParent(c)
		inj.Map(&confraw)

		refvs, err := inj.Invoke(handler)
		if err != nil {
			return []TypeFilterConfig{}, err
		}

		for _, refv := range refvs {
			if !refv.CanInterface() {
				continue
			}
			if conf, ok := refv.Interface().(TypeFilterConfig); ok {
				conf.SetInjector(inj)
				filters = append(filters, conf)
			}
		}
	}
	return
}
