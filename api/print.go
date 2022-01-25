package api

import (
	"fmt"
	"net/http"
	"strings"

	v8 "rogchap.com/v8go"
)

func RegisterPrintApi(iso *v8.Isolate, global *v8.ObjectTemplate, w http.ResponseWriter, r *http.Request) error {
	global.Set("print", v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		v := fmt.Sprintf("%v", info.Args())
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			v = v[1 : len(v)-1]
		}
		fmt.Fprintf(w, v)
		return nil
	}))

	return nil
}
