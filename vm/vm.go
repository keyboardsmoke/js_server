package vm

import (
	"bytes"
	"net/http"
	"os"

	v8 "rogchap.com/v8go"
)

type Vm struct {
	Iso    *v8.Isolate
	Ctx    *v8.Context
	Global *v8.ObjectTemplate
}

type RegisterCallback (func(*Vm, http.ResponseWriter, *http.Request) error)

var callbacks []RegisterCallback = []RegisterCallback{}

func CreateVm() *Vm {
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)
	return &Vm{
		Iso:    iso,
		Ctx:    nil,
		Global: global,
	}
}

func (vm *Vm) AddInterface(cb RegisterCallback) {
	callbacks = append(callbacks, cb)
}

func (vm *Vm) Finalize(w http.ResponseWriter, r *http.Request) error {
	for _, cb := range callbacks {
		err := cb(vm, w, r)
		if err != nil {
			return err
		}
	}

	vm.Ctx = v8.NewContext(vm.Iso, vm.Global)

	return nil
}

func (vm *Vm) ExecuteScript(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)

	_, err = vm.Ctx.RunScript(buf.String(), filename)
	if err != nil {
		return err
	}

	return nil
}
