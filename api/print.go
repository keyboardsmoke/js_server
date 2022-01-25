package api

import (
	"fmt"
	vm "js_server/vm"
	"net/http"
	"strings"

	v8 "rogchap.com/v8go"
)

func RegisterPrintApi(v *vm.Vm, w http.ResponseWriter, r *http.Request) error {
	v.Global.Set("print", v8.NewFunctionTemplate(v.Iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		v := fmt.Sprintf("%v", info.Args())
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			v = v[1 : len(v)-1]
		}
		fmt.Fprintf(w, v)
		return nil
	}))

	return nil
}
