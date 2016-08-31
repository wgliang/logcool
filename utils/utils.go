package utils

import (
	"fmt"
)

func CheckErrorValues(refvs []reflect.Value) (err error) {
	for _, refv := range refvs {
		if refv.IsValid() {
			refvi := refv.Interface()
			switch refvi.(type) {
			case error:
				return refvi.(error)
			}
		}
	}
	return
}
