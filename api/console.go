package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	v8 "rogchap.com/v8go"
)

func RegisterConsoleApi(iso *v8.Isolate, global *v8.ObjectTemplate, w http.ResponseWriter, r *http.Request) error {

	console := v8.NewObjectTemplate(iso)

	console.Set("log", v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		v := fmt.Sprintf("%v", info.Args())
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			v = v[1 : len(v)-1]
		}

		fmt.Fprintf(os.Stdout, v)
		return nil
	}))

	global.Set("console", console)

	return nil
}
