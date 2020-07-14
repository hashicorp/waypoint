// Package logsilence exists only to silence warnings of duplicate file
// warnings from the proto package. The proto package outputs this using
// `log`. We must therefore silence the log.
//
// To use this, import this package before any other package that imports
// protobuf files. Based on the Go spec, initialization is based on import
// order.
package logsilence

import (
	"io/ioutil"
	"log"
)

func init() {
	log.SetOutput(ioutil.Discard)
}
