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
	callEchoExample()
	callChannelCloseExample()
	callAnonExample()
	callMutexLockExample()
	callMulitpleChannelsExample()
	callLockingWithChannelsExample()
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

func send(ch chan string, done <-chan bool) {
	for {
		if <-done {
			close(ch)
			fmt.Println("[callChannelClosedExample] ch closed")
			return
		}

		ch <- "hello"
		time.Sleep(200 * time.Millisecond)
	}
}

func callChannelCloseExample() {
	fmt.Println("[callChannelCloseExample] start")
	msg := make(chan string)
	until := time.After(5 * time.Second)
	done := make(chan bool)

	go send(msg, done)

	for {
		select {
		case m := <-msg:
			fmt.Println(m)
		case <-until:
			done <- true
			fmt.Println("[callChannelClosedExampled] Timed out")
			return
		default:
			fmt.Println("**yawn**")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func worker(id int, lock chan bool) {
	fmt.Println("[callLockingWithChannelsExample] start worker ", id)
	lock <- true
	fmt.Println("[callLockingWithChannelsExample] locked worker ", id)
	time.Sleep(500 * time.Millisecond)
	fmt.Println(id, " is releasing its lock")
	<-lock
}

func callLockingWithChannelsExample() {
	lock := make(chan bool, 1)
	for i := 0; i < 12; i++ {
		go worker(i, lock)
		time.Sleep(10 * time.Second)
	}
}
