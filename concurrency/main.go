package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

func main() {
	callEchoExample()
	callAnonExample()
	os.Exit(0)
}

func echo(in io.Reader, out io.Writer) {
	io.Copy(out, in)
}

func callEchoExample() {
	go echo(os.Stdin, os.Stdout)
	time.Sleep(10 * time.Second)
	fmt.Println("[callEchoExample] Timed out.")
}

func callAnonExample() {
	fmt.Println("[callAnonExample] outside goroutine")
	go func() {
		fmt.Println("[callAnonExample] inside a goroutine")
	}()
	fmt.Println("[callAnonExample] outside goroutine again")
	runtime.Gosched()
}
