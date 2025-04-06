package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
    msg = "Hello, world!"
    updateMessage(msg)
    if msg != "Hello, world!" {
        t.Errorf("updateMessage() = %v, want %v", msg, "Hello, world!")
    }
}

func Test_printMessage(t *testing.T) {
    stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
    msg = "Hello, universe!"
    var wg sync.WaitGroup
    wg.Add(1)
	go printMessage(&wg)

	wg.Wait()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut
	if !strings.Contains(output, msg) {
		t.Errorf("printMessage() = %v, want %v", output, msg)
	}
}

func Test_main(t *testing.T) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	main()
	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedMessages := []string{
		"Hello, universe!",
		"Hello, cosmos!",
		"Hello, world!",
	}
    if len(lines) != len(expectedMessages) {
        t.Errorf("main() output length = %d, want %d", len(lines), len(expectedMessages))
    }
    for i, line := range lines {
        if line != expectedMessages[i] {
            t.Errorf("main() output %d = %v, want %v", i, line, expectedMessages[i])
        }
    }
}
