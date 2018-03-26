package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
)

func tee() string {
	var b1 bytes.Buffer
	var b2 bytes.Buffer

	tee := io.TeeReader(os.Stdin, &b1)
	s := bufio.NewScanner(tee)
	for s.Scan() {
		fmt.Println(s.Text())
	}

	uncolorize := colorable.NewNonColorable(&b2)
	uncolorize.Write(b1.Bytes())

	return b2.String()
}
