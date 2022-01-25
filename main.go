package main

import (
	"bytes"
	"fmt"
	"js_server/api"
	"net/http"
	"net/http/fcgi"
	"os"
	"strings"

	v8 "rogchap.com/v8go"
)

func openVmFile(iso *v8.Isolate, ctx *v8.Context, global *v8.ObjectTemplate, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	contents := buf.String()

	_, err = ctx.RunScript(contents, filename)
	if err != nil {
		return err
	}

	return nil
}

func createVm(w http.ResponseWriter, r *http.Request) (*v8.Isolate, *v8.ObjectTemplate) {
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)
	return iso, global
}

func createCtx(iso *v8.Isolate, global *v8.ObjectTemplate) *v8.Context {
	return v8.NewContext(iso, global)
}

func fillHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Powered-By", "v8go")
	w.Header().Set("X-v8-Version", v8.Version())
	w.Header().Set("X-v8-Request-Method", r.Method)
	w.Header().Set("X-v8-Request-URL", r.URL.RequestURI())
	w.Header().Set("X-v8-Request-Proto", r.Proto)
	w.Header().Set("X-v8-Request-Host", r.Host)
	w.Header().Set("X-v8-Request-Path", r.URL.Path)
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
		fmt.Fprintln(w, "fcgi.ProcessEnv: ", cgiEnv, "<br />")

		return
	}

	fillHeader(w, r)

	iso, global := createVm(w, r)

	ae := api.RegisterApi(iso, global, w, r)
	if ae != nil {
		fmt.Fprintln(w, ae)
		return
	}

	ctx := createCtx(iso, global)

	ve := openVmFile(iso, ctx, global, cgiEnv["SCRIPT_FILENAME"])
	if ve != nil {
		fmt.Fprintln(w, ve)
	}
}

func main() {
	err := fcgi.Serve(nil, &FCGIHandler{})
	if err != nil {
		panic(err)
	}
}
