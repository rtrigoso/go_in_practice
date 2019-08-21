package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	// callEchoExample()
	callAnonExample()
	callMutexLockExample()
	callMulitpleChannelsExample()
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

type words struct {
	sync.Mutex
	found map[string]int
}

func newWords() *words {
	return &words{found: map[string]int{}}
}

func (w *words) add(word string, n int) {
	w.Lock()
	defer w.Unlock()
	count, ok := w.found[word]
	if !ok {
		w.found[word] = n
		return
	}

	w.found[word] = count + n
}

func tallyWords(filename string, dict *words) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words := strings.ToLower(scanner.Text())
		dict.add(words, 1)
	}

	return scanner.Err()
}

func callMutexLockExample() {
	fmt.Println("[callMutexLockExample] start")
	var wg sync.WaitGroup
	files := []string{"./data/tallytest.txt", "./data/tallytest1.txt"}

	w := newWords()
	for _, f := range files {
		wg.Add(1)
		go func(file string) {
			fmt.Println("[callMutexLockExample] testing file ", file)
			if err := tallyWords(file, w); err != nil {
				fmt.Printf("[callMutexLockExample] error: %s", err.Error())
				os.Exit(1)
			}
			wg.Done()
		}(f)
	}
	wg.Wait()

	fmt.Println("[callMutexLockExample] Words that appear more than once:")

	for word, count := range w.found {
		if count > 1 {
			fmt.Println(word, " - ", count)
		}
	}
}

func readStdin(out chan<- []byte) {
	for {
		data := make([]byte, 1024)
		l, _ := os.Stdin.Read(data)
		if l > 0 {
			out <- data
		}
	}
}

func callMulitpleChannelsExample() {
	fmt.Println("[callMulitpleChannelsExample] start")
	done := time.After(10 * time.Second)
	echo := make(chan []byte)
	go readStdin(echo)
	for {
		select {
		case buf := <-echo:
			os.Stdout.Write(buf)
		case <-done:
			fmt.Println("Timed out")
			os.Exit(0)
		}
	}
}
