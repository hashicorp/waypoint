package lambda

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
)

type LogPrinter struct {
	Prefix string
}

func (lg *LogPrinter) Display(r io.Reader) {
	sor, sow, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	go lg.watch(sor, false)

	ser, sew, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	go lg.watch(ser, true)

	stdcopy.StdCopy(sow, sew, r)
}

func (lg *LogPrinter) watch(r io.Reader, errors bool) {
	br := bufio.NewReader(r)

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return
		}

		fmt.Printf("%s %s", lg.Prefix, line)
	}
}
