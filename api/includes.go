package api

import (
	"fmt"
	vm "js_server/vm"
	"net/http"
	"os"

	v8 "rogchap.com/v8go"
)

func RegisterIncludesApi(v *vm.Vm, w http.ResponseWriter, r *http.Request) error {
	v.Global.Set("require", v8.NewFunctionTemplate(v.Iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		scriptname := info.Args()[0].String()
		err := v.ExecuteScript(scriptname)
		if err != nil {
			wd, _ := os.Getwd()
			fmt.Fprintln(w, err, "<br />")
			fmt.Fprintln(w, "Error loading script:", scriptname, "<br />")
			fmt.Fprintln(w, "Please check the file exists and is in the correct directory.<br />")
			fmt.Fprintln(w, "Working Directory: ", wd, "<br />")
		}

		return nil
	}))

	return nil
}
