package api

import (
	"fmt"
	vm "js_server/vm"
	"net/http"
	"os"
	"strings"

	v8 "rogchap.com/v8go"
)

func RegisterConsoleApi(v *vm.Vm, w http.ResponseWriter, r *http.Request) error {

	console := v8.NewObjectTemplate(v.Iso)

	console.Set("log", v8.NewFunctionTemplate(v.Iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		v := fmt.Sprintf("%v", info.Args())
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			v = v[1 : len(v)-1]
		}

		fmt.Fprintf(os.Stdout, v)
		return nil
	}))

	v.Global.Set("console", console)

	return nil
}
