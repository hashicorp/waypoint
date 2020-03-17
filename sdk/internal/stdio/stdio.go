package stdio

import (
	"os"
)

// NOTE(mitchellh): this won't work on Windows. We need to do something like
// this: https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilea?redirectedfrom=MSDN

// Stdout returns the stdout file that was passed as an extra file descriptor
// to the plugin. We do this so that we can get access to a real TTY if
// possible for subprocess output.
func Stdout() *os.File {
	return os.NewFile(uintptr(3), "stdout")
}

// Stderr. See stdout for details.
func Stderr() *os.File {
	return os.NewFile(uintptr(4), "stderr")
}
