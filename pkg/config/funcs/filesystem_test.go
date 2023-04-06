// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/zclconf/go-cty/cty"
)

func TestFile(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("testdata/filesystem/hello.txt"),
			cty.StringVal("Hello World"),
			false,
		},
		{
			cty.StringVal("testdata/filesystem/icon.png"),
			cty.NilVal,
			true, // Not valid UTF-8
		},
		{
			cty.StringVal("testdata/filesystem/missing"),
			cty.NilVal,
			true, // no file exists
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("File(\".\", %#v)", test.Path), func(t *testing.T) {
			abs, err := filepath.Abs(test.Path.AsString())
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			test.Path = cty.StringVal(abs)

			got, err := File(test.Path)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("testdata/filesystem/hello.txt"),
			cty.BoolVal(true),
			false,
		},
		{
			cty.StringVal(""), // empty path
			cty.BoolVal(false),
			true,
		},
		{
			cty.StringVal("testdata/filesystem/missing"),
			cty.BoolVal(false),
			false, // no file exists
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("FileExists(\".\", %#v)", test.Path), func(t *testing.T) {
			abs, err := filepath.Abs(test.Path.AsString())
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			test.Path = cty.StringVal(abs)

			got, err := FileExists(test.Path)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestFileSet(t *testing.T) {
	tests := []struct {
		Path    cty.Value
		Pattern cty.Value
		Want    cty.Value
		Err     bool
	}{
		{
			cty.StringVal("."),
			cty.StringVal("testdata*"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("{testdata,missing}"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/missing"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/missing*"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("*/missing"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("**/missing"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/filesystem/*.txt"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/filesystem/hello.txt"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/filesystem/hello.???"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/filesystem/hello*"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/hello.tmpl"),
				cty.StringVal("testdata/filesystem/hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/filesystem/hello.{tmpl,txt}"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/hello.tmpl"),
				cty.StringVal("testdata/filesystem/hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("*/*/hello.txt"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("testdata/**/list*"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/list.tmpl"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("**/hello.{tmpl,txt}"),
			cty.SetVal([]cty.Value{
				cty.StringVal("testdata/filesystem/hello.tmpl"),
				cty.StringVal("testdata/filesystem/hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("."),
			cty.StringVal("["),
			cty.SetValEmpty(cty.String),
			true,
		},
		{
			cty.StringVal("."),
			cty.StringVal("\\"),
			cty.SetValEmpty(cty.String),
			true,
		},
		{
			cty.StringVal("testdata/filesystem"),
			cty.StringVal("missing"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("testdata/filesystem"),
			cty.StringVal("missing*"),
			cty.SetValEmpty(cty.String),
			false,
		},
		{
			cty.StringVal("testdata/filesystem"),
			cty.StringVal("*.txt"),
			cty.SetVal([]cty.Value{
				cty.StringVal("hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("testdata/filesystem"),
			cty.StringVal("hello.txt"),
			cty.SetVal([]cty.Value{
				cty.StringVal("hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("testdata/filesystem"),
			cty.StringVal("hello.???"),
			cty.SetVal([]cty.Value{
				cty.StringVal("hello.txt"),
			}),
			false,
		},
		{
			cty.StringVal("testdata/filesystem"),
			cty.StringVal("hello*"),
			cty.SetVal([]cty.Value{
				cty.StringVal("hello.tmpl"),
				cty.StringVal("hello.txt"),
			}),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("FileSet(\".\", %#v, %#v)", test.Path, test.Pattern), func(t *testing.T) {
			abs, err := filepath.Abs(test.Path.AsString())
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			test.Path = cty.StringVal(abs)

			got, err := FileSet(test.Path, test.Pattern)

			if test.Err {
				if err == nil {
					t.Fatalf("succeeded; want error\ngot:  %#v", got)
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestFileBase64(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("testdata/filesystem/hello.txt"),
			cty.StringVal("SGVsbG8gV29ybGQ="),
			false,
		},
		{
			cty.StringVal("testdata/filesystem/icon.png"),
			cty.StringVal("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAq1BMVEX///9cTuVeUeRcTuZcTuZcT+VbSe1cTuVdT+MAAP9JSbZcT+VcTuZAQLFAQLJcTuVcTuZcUuBBQbA/P7JAQLJaTuRcT+RcTuVGQ7xAQLJVVf9cTuVcTuVGRMFeUeRbTeJcTuU/P7JeTeZbTOVcTeZAQLJBQbNAQLNaUORcTeZbT+VcTuRAQLNAQLRdTuRHR8xgUOdgUN9cTuVdTeRdT+VZTulcTuVAQLL///8+GmETAAAANnRSTlMApibw+osO6DcBB3fIX87+oRk3yehB0/Nj/gNs7nsTRv3dHmu//JYUMLVr3bssjxkgEK5CaxeK03nIAAAAAWJLR0QAiAUdSAAAAAlwSFlzAAADoQAAA6EBvJf9gwAAAAd0SU1FB+EEBRIQDxZNTKsAAACCSURBVBjTfc7JFsFQEATQQpCYxyBEzJ55rvf/f0ZHcyQLvelTd1GngEwWycs5+UISyKLraSi9geWKK9Gr1j7AeqOJVtt2XtD1Bchef2BjQDAcCTC0CsA4mihMtXw2XwgsV2sFw812F+4P3y2GdI6nn3FGSs//4HJNAXDzU4Dg/oj/E+bsEbhf5cMsAAAAJXRFWHRkYXRlOmNyZWF0ZQAyMDE3LTA0LTA1VDE4OjE2OjE1KzAyOjAws5bLVQAAACV0RVh0ZGF0ZTptb2RpZnkAMjAxNy0wNC0wNVQxODoxNjoxNSswMjowMMLLc+kAAAAZdEVYdFNvZnR3YXJlAHd3dy5pbmtzY2FwZS5vcmeb7jwaAAAAC3RFWHRUaXRsZQBHcm91cJYfIowAAABXelRYdFJhdyBwcm9maWxlIHR5cGUgaXB0YwAAeJzj8gwIcVYoKMpPy8xJ5VIAAyMLLmMLEyMTS5MUAxMgRIA0w2QDI7NUIMvY1MjEzMQcxAfLgEigSi4A6hcRdPJCNZUAAAAASUVORK5CYII="),
			false,
		},
		{
			cty.StringVal("testdata/filesystem/missing"),
			cty.NilVal,
			true, // no file exists
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("FileBase64(\".\", %#v)", test.Path), func(t *testing.T) {
			abs, err := filepath.Abs(test.Path.AsString())
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			test.Path = cty.StringVal(abs)

			got, err := FileBase64(test.Path)
			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestBasename(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("testdata/filesystem/hello.txt"),
			cty.StringVal("hello.txt"),
			false,
		},
		{
			cty.StringVal("hello.txt"),
			cty.StringVal("hello.txt"),
			false,
		},
		{
			cty.StringVal(""),
			cty.StringVal("."),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Basename(%#v)", test.Path), func(t *testing.T) {
			got, err := Basename(test.Path)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestDirname(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("testdata/filesystem/hello.txt"),
			cty.StringVal("testdata/filesystem"),
			false,
		},
		{
			cty.StringVal("testdata/filesystem/foo/hello.txt"),
			cty.StringVal("testdata/filesystem/foo"),
			false,
		},
		{
			cty.StringVal("hello.txt"),
			cty.StringVal("."),
			false,
		},
		{
			cty.StringVal(""),
			cty.StringVal("."),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Dirname(%#v)", test.Path), func(t *testing.T) {
			got, err := Dirname(test.Path)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestPathExpand(t *testing.T) {
	homePath, err := homedir.Dir()
	if err != nil {
		t.Fatalf("Error getting home directory: %v", err)
	}

	tests := []struct {
		Path cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("~/test-file"),
			cty.StringVal(filepath.Join(homePath, "test-file")),
			false,
		},
		{
			cty.StringVal("~/another/test/file"),
			cty.StringVal(filepath.Join(homePath, "another/test/file")),
			false,
		},
		{
			cty.StringVal("/root/file"),
			cty.StringVal("/root/file"),
			false,
		},
		{
			cty.StringVal("/"),
			cty.StringVal("/"),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Dirname(%#v)", test.Path), func(t *testing.T) {
			got, err := Pathexpand(test.Path)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
