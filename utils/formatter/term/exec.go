package term

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Exec(cmd string, dir string) (out chan string, err error) {
	out = make(chan string)
	var outb, errb bytes.Buffer
	toks := strings.Split(cmd, " ")
	c := exec.Command(toks[0], toks[1:]...)
	c.Dir = dir
	c.Stdout = &outb
	c.Stderr = &errb

	go func() {
		for output := range out {
			fmt.Println(output)
		}
	}()

	go func() (e error) {
		if e = c.Run(); e != nil {
			return
		}
	loop:
		for {
			tmp := make([]byte, 1024)
			if _, e = outb.Read(tmp); e != nil {
				out <- e.Error()
				break loop
			} else {
				out <- outb.String()
			}
		}
		return
	}()
	c.Wait()
	return
}
