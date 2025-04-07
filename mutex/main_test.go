package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	// keep a copy of the os standard output
	stdOut := os.Stdout

	// get the real output
	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut
	if !strings.Contains(output, "$83200.00") {
		t.Errorf("Expected output to contain '$%d.00'", 83200)
	}
}
