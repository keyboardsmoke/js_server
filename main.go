package main

import (
	"fmt"
	"js_server/api"
	"js_server/vm"
	"net/http"
	"net/http/fcgi"
	"os"
	"path/filepath"
	"strings"

	v8 "rogchap.com/v8go"
)

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

func releaseStats(filename, dirfile, wd string, cgiEnv map[string]string, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Filename: ", filename, "<br />")
	fmt.Fprintln(w, "Working Directory: ", dirfile, "<br />")
	fmt.Fprintln(w, "Old Working Directory: ", wd, "<br />")
	fmt.Fprintln(w, "URL: ", r.URL, "<br />")
	fmt.Fprintln(w, "URL.Fragment: ", r.URL.Fragment, "<br />")
	fmt.Fprintln(w, "URL.Path: ", r.URL.Path, "<br />")
	fmt.Fprintln(w, "URL.RequestURI: ", r.URL.RequestURI(), "<br />")
	fmt.Fprintln(w, "URL.Scheme: ", r.URL.Scheme, "<br />")
	fmt.Fprintln(w, "URL.User: ", r.URL.User, "<br />")
	fmt.Fprintln(w, "URL.Host: ", r.URL.Host, "<br />")
	fmt.Fprintln(w, "URL.Port: ", r.URL.Port(), "<br />")
	fmt.Fprintln(w, "URL.RawPath: ", r.URL.RawPath, "<br />")
	fmt.Fprintln(w, "URL.RawQuery: ", r.URL.RawQuery, "<br />")
	fmt.Fprintln(w, "URL.Query(): ", r.URL.Query(), "<br />")
	fmt.Fprintln(w, "Method: ", r.Method, "<br />")
	fmt.Fprintln(w, "Header: ", r.Header, "<br />")
	fmt.Fprintln(w, "RemoteAddr: ", r.RemoteAddr, "<br />")
	fmt.Fprintln(w, "RequestURI: ", r.RequestURI, "<br />")
	fmt.Fprintln(w, "Proto: ", r.Proto, "<br />")
	fmt.Fprintln(w, "Host: ", r.Host, "<br />")
	fmt.Fprintln(w, "Form: ", r.Form, "<br />")
	fmt.Fprintln(w, "fcgi.ProcessEnv: ", cgiEnv, "<br />")
	fmt.Fprintln(w, "os.Environ(): ", os.Environ(), "<br />")
	fmt.Fprintln(w, "os.Args: ", os.Args, "<br />")
}

type FCGIHandler struct {
}

func (f *FCGIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	cgiEnv := fcgi.ProcessEnv(r)

	fullpath := cgiEnv["SCRIPT_FILENAME"]
	dirfile := filepath.Dir(fullpath)
	filename, aerr := filepath.Abs(fullpath)
	if aerr != nil {
		fmt.Fprintf(w, "Error: %v", aerr)
		return
	}

	if strings.HasSuffix(r.URL.Path, "/stat") {
		releaseStats(filename, dirfile, wd, cgiEnv, w, r)
		return
	}

	fillHeader(w, r)

	v := vm.CreateVm()

	v.AddInterface(api.RegisterPrintApi)
	v.AddInterface(api.RegisterConsoleApi)
	v.AddInterface(api.RegisterIncludesApi)

	err = v.Finalize(w, r)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = os.Chdir(dirfile)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = v.ExecuteScript(fullpath)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = os.Chdir(wd)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
}

func isRunningAsFcgi() bool {
	return os.Getenv("PHP_FCGI_CHILDREN") != ""
}

func main() {
	if isRunningAsFcgi() {
		fmt.Printf("%v", os.Environ())

		err := fcgi.Serve(nil, &FCGIHandler{})
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Not yet implemented.")
	}
}
