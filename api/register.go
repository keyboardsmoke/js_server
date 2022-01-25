package api

import (
	"net/http"

	v8 "rogchap.com/v8go"
)

type registerCallback (func(*v8.Isolate, *v8.ObjectTemplate, http.ResponseWriter, *http.Request) error)

var callbacks []registerCallback = []registerCallback{
	RegisterConsoleApi,
	RegisterPrintApi,
}

func RegisterApi(iso *v8.Isolate, global *v8.ObjectTemplate, w http.ResponseWriter, r *http.Request) error {
	for _, cb := range callbacks {
		err := cb(iso, global, w, r)
		if err != nil {
			return err
		}
	}
	return nil
}
