package api

import (
	"net/http"

	"github.com/robertkrimen/otto"
)

type registerCallback (func(*otto.Otto, http.ResponseWriter, *http.Request) error)

var callbacks []registerCallback = []registerCallback{
	RegisterConsoleApi,
	RegisterPrintApi,
}

func RegisterApi(vm *otto.Otto, w http.ResponseWriter, r *http.Request) error {
	for _, cb := range callbacks {
		err := cb(vm, w, r)
		if err != nil {
			return err
		}
	}
	return nil
}
