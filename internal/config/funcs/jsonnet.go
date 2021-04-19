package funcs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func Jsonnet() map[string]function.Function {
	return map[string]function.Function{
		"jsonnetdir":  JsonnetDirFunc,
		"jsonnetfile": JsonnetFileFunc,
	}
}

func jsonnetVM(opts cty.Value) *jsonnet.VM {
	vm := jsonnet.MakeVM()

	// Get our options. If it isn't a valid type that we're looking for
	// then just return the VM as-is.
	optsT := opts.Type()
	if opts.IsNull() {
		return vm
	}
	if !optsT.IsObjectType() && !optsT.IsMapType() {
		return vm
	}

	// Jsonnet has FOUR ways to parameterize a file, all with string key/values.
	// This function abstracts setting it so we can set all four methods.
	setter := func(k string, cb func(string, string)) {
		val := attr(opts, k)
		if !val.IsNull() && val.CanIterateElements() {
			for k, v := range val.AsValueMap() {
				cb(k, v.AsString())
			}

		}
	}
	setter("ext_vars", vm.ExtVar)
	setter("ext_code", vm.ExtCode)
	setter("tla_vars", vm.TLAVar)
	setter("tla_code", vm.TLACode)

	return vm
}

func attr(v cty.Value, k string) cty.Value {
	t := v.Type()
	if t.IsMapType() {
		return v.Index(cty.StringVal(k))
	} else if t.HasAttribute(k) {
		return v.GetAttr(k)
	}

	return cty.NullVal(cty.String)
}

// JsonnetDirFunc constructs a function that converts a directory of
// jsonnet files into standard JSON files with the same name but a "json"
// extension instead of the "jsonnet" extension. The converted files are
// stored in a temporary directory that is returned; the original files are
// untouched.
var JsonnetDirFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "dir",
			Type: cty.String,
		},
		{
			Name: "options",
			Type: cty.DynamicPseudoType,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		td, err := ioutil.TempDir("", "waypoint")
		if err != nil {
			return cty.DynamicVal, err
		}

		root := args[0].AsString()
		vm := jsonnetVM(args[1])
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			dir := td

			// Determine if we have any directory
			stripped := strings.TrimPrefix(path, root)
			if len(stripped) == 0 {
				panic("empty path") // should never happen
			}
			if stripped[0] == '/' || stripped[0] == '\\' {
				// Get rid of any prefix '/' which could happen if tpl.Path doesn't
				// end in a directory sep.
				stripped = stripped[1:]
			}
			if v := filepath.Dir(stripped); v != "." {
				dir = filepath.Join(dir, v)
				if err := os.MkdirAll(dir, 0700); err != nil {
					return err
				}
			}

			// Render
			jsonStr, err := vm.EvaluateFile(path)
			if err != nil {
				return err
			}

			// We ensure the filename ends with ".json"
			filename := filepath.Base(path)
			filename = strings.TrimSuffix(filename, ".jsonnet")
			filename += ".json"

			// We'll copy the file into the temporary directory
			path = filepath.Join(dir, filename)
			return ioutil.WriteFile(path, []byte(jsonStr), 0600)
		})
		if err != nil {
			return cty.DynamicVal, err
		}

		return cty.StringVal(td), nil
	},
})

// JsonnetFileFunc constructs a function that converts a single Jsonnet file
// to JSON. This returns the path to a new file with a "json" extension.
var JsonnetFileFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "path",
			Type: cty.String,
		},
		{
			Name: "options",
			Type: cty.DynamicPseudoType,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		td, err := ioutil.TempDir("", "waypoint")
		if err != nil {
			return cty.DynamicVal, err
		}

		path := args[0].AsString()

		// Render
		vm := jsonnetVM(args[1])
		jsonStr, err := vm.EvaluateFile(path)
		if err != nil {
			return cty.DynamicVal, err
		}

		// We ensure the filename ends with ".json"
		filename := filepath.Base(path)
		filename = strings.TrimSuffix(filename, ".jsonnet")
		filename += ".json"

		// Write the file
		path = filepath.Join(td, filename)
		if err := ioutil.WriteFile(path, []byte(jsonStr), 0600); err != nil {
			return cty.DynamicVal, err
		}
		return cty.StringVal(path), nil
	},
})

// JsonnetFile converts a single Jsonnet file to JSON.
func JsonnetFile(dir cty.Value) (cty.Value, error) {
	return JsonnetFileFunc.Call([]cty.Value{dir})
}

// JsonnetDir converts a directory of Jsonnet files into JSON.
func JsonnetDir(dir cty.Value) (cty.Value, error) {
	return JsonnetDirFunc.Call([]cty.Value{dir})
}
