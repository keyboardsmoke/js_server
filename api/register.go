package api

import (
	"net/http"

	"github.com/robertkrimen/otto"
)

func RegisterApi(vm *otto.Otto, w http.ResponseWriter, r *http.Request) error {
	err := RegisterConsoleApi(vm, w, r)
	if err != nil {
		return err
	}

	err = RegisterPrintApi(vm, w, r)
	if err != nil {
		return err
	}

	return nil
}
