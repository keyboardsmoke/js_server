package api

import (
	"fmt"
	"net/http"

	"github.com/robertkrimen/otto"
)

func RegisterConsoleApi(vm *otto.Otto, w http.ResponseWriter, r *http.Request) error {

	obj, err := vm.Object("console")
	if err != nil {
		fmt.Println(err)
		return err
	}

	obj.Set("log", func(call otto.FunctionCall) otto.Value {
		fmt.Println(call.Argument(0).String())
		return otto.Value{}
	})

	return nil
}
