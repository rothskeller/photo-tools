package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// An EXIFTool is a handle to an external exiftool invocation.  For efficiency,
// it is only invoked once and that invocation is fed multiple commands.
type EXIFTool struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
}

// Open invokes exiftool and saves a handle to it to send commands to it.
func (et *EXIFTool) Open() (err error) {
	var (
		stdout io.ReadCloser
	)
	et.cmd = exec.Command("exiftool", "-stay_open", "True", "-@", "-")
	if et.stdin, err = et.cmd.StdinPipe(); err != nil {
		return fmt.Errorf("pipe to exiftool: %s", err)
	}
	if stdout, err = et.cmd.StdoutPipe(); err != nil {
		return fmt.Errorf("pipe from exiftool: %s", err)
	}
	et.stdout = bufio.NewScanner(stdout)
	et.cmd.Stderr = os.Stderr
	if err = et.cmd.Start(); err != nil {
		return fmt.Errorf("exiftool: %s", err)
	}
	return nil
}

// Close shuts down the invocation of exiftool.
func (et *EXIFTool) Close() (err error) {
	if _, err = et.stdin.Write([]byte("-stay_open\nFalse\n")); err != nil {
		return fmt.Errorf("write exiftool: %s", err)
	}
	if err = et.stdin.Close(); err != nil {
		return fmt.Errorf("close to exiftool: %s", err)
	}
	for et.stdout.Scan() {
	}
	if err = et.stdout.Err(); err != nil {
		return fmt.Errorf("close from exiftool: %s", err)
	}
	if err = et.cmd.Wait(); err != nil {
		return fmt.Errorf("exiftool: %s", err)
	}
	return nil
}

// Run runs a single exiftool command, specified as a list of arguments.  It
// returns the lines of output generated by the command.
func (et *EXIFTool) Run(args ...string) (out []string, err error) {
	var argbuf bytes.Buffer
	for _, a := range args {
		fmt.Fprintf(&argbuf, "%s\n", a)
	}
	fmt.Fprintln(&argbuf, "-execute")
	if _, err = et.stdin.Write(argbuf.Bytes()); err != nil {
		return nil, fmt.Errorf("write exiftool: %s", err)
	}
	for et.stdout.Scan() {
		var line = et.stdout.Text()
		if line == "{ready}" {
			return
		}
		out = append(out, line)
	}
	if err = et.stdout.Err(); err != nil {
		return nil, fmt.Errorf("read exiftool: %s", err)
	}
	return nil, fmt.Errorf("read exiftool: %s", io.EOF)
}