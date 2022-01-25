package main

import (
	"fmt"
	"js_server/api"
	"net/http"
	"net/http/fcgi"
	"os"
	"strings"

	"github.com/robertkrimen/otto"
)

func openVmFile(vm *otto.Otto, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	vm.Run(f)

	return nil
}

func createVm(w http.ResponseWriter, r *http.Request) *otto.Otto {
	vm := otto.New()

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Powered-By", "Otto")
	w.Header().Set("X-Otto-Version", "1.0")
	w.Header().Set("X-Otto-Request-Method", r.Method)
	w.Header().Set("X-Otto-Request-URL", r.URL.RequestURI())
	w.Header().Set("X-Otto-Request-Proto", r.Proto)
	w.Header().Set("X-Otto-Request-Host", r.Host)
	w.Header().Set("X-Otto-Request-Path", r.URL.Path)
	w.Header().Set("X-Otto-Request-Query", r.URL.RawQuery)

	api.RegisterApi(vm, w, r)

	return vm
}

type FCGIHandler struct {
}

func (f *FCGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cgiEnv := fcgi.ProcessEnv(r)

	if strings.HasSuffix(r.URL.Path, "/stat") {
		fmt.Fprintln(w, "Hello World<br /><br />")
		fmt.Fprintln(w, "URL: ", r.URL, "<br /><br />")
		fmt.Fprintln(w, "URL.Fragment: ", r.URL.Fragment, "<br /><br />")
		fmt.Fprintln(w, "URL.Path: ", r.URL.Path, "<br /><br />")
		fmt.Fprintln(w, "URL.RequestURI: ", r.URL.RequestURI(), "<br /><br />")
		fmt.Fprintln(w, "URL.Scheme: ", r.URL.Scheme, "<br /><br />")
		fmt.Fprintln(w, "URL.User: ", r.URL.User, "<br /><br />")
		fmt.Fprintln(w, "URL.Host: ", r.URL.Host, "<br /><br />")
		fmt.Fprintln(w, "URL.Port: ", r.URL.Port(), "<br /><br />")
		fmt.Fprintln(w, "URL.RawPath: ", r.URL.RawPath, "<br /><br />")
		fmt.Fprintln(w, "URL.RawQuery: ", r.URL.RawQuery, "<br /><br />")
		fmt.Fprintln(w, "URL.Query(): ", r.URL.Query(), "<br /><br />")
		fmt.Fprintln(w, "Method: ", r.Method, "<br /><br />")
		fmt.Fprintln(w, "Header: ", r.Header, "<br /><br />")
		fmt.Fprintln(w, "RemoteAddr: ", r.RemoteAddr, "<br /><br />")
		fmt.Fprintln(w, "RequestURI: ", r.RequestURI, "<br /><br />")
		fmt.Fprintln(w, "Proto: ", r.Proto, "<br /><br />")
		fmt.Fprintln(w, "Host: ", r.Host, "<br /><br />")
		fmt.Fprintln(w, "Form: ", r.Form, "<br /><br />")
		fmt.Fprintln(w, "Environment: ", os.Environ(), "<br /><br />")

		fmt.Fprintln(w, "fcgi.ProcessEnv: ", cgiEnv, "<br />")

		return
	}

	vm := createVm(w, r)
	err := openVmFile(vm, cgiEnv["SCRIPT_FILENAME"])
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func main() {
	fcgi.Serve(nil, &FCGIHandler{})
}
