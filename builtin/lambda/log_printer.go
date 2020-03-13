package lambda

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/stdcopy"
)

type LogPrinter struct {
	Prefix string
}

func (lg *LogPrinter) Display(r io.Reader) {

	var wg sync.WaitGroup

	sor, sow, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go lg.watch(sor, false, &wg)

	ser, sew, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go lg.watch(ser, true, &wg)

	stdcopy.StdCopy(sow, sew, r)

	sow.Close()
	sew.Close()

	wg.Wait()
}

func (lg *LogPrinter) watch(r io.Reader, errors bool, wg *sync.WaitGroup) {
	defer wg.Done()

	br := bufio.NewReader(r)

	for {
		line, err := br.ReadString('\n')
		if line != "" {
			fmt.Printf("%s %s\n", lg.Prefix, strings.TrimRight(line, "\n "))
		}

		if err != nil {
			return
		}
	}
}
