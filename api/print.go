package api

import (
	"net/http"

	"github.com/robertkrimen/otto"
)

func RegisterPrintApi(vm *otto.Otto, w http.ResponseWriter, r *http.Request) error {
	vm.Set("echo", func(str string) {
		w.Write([]byte(str))
	})

	return nil
}
