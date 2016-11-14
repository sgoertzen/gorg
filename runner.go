package gorg

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

var waitTime, _ = time.ParseDuration("10s")

func runWithRetries(dir string, command string, arg ...string) (int, error) {
	var result int
	var err error
	retries := 0
	for {
		result, err = run(dir, command, arg...)
		if err == nil || retries > 5 {
			break
		}
		retries++
		log.Println("Retrying...")
		time.Sleep(waitTime)
	}
	return result, err
}

func run(dir string, command string, arg ...string) (int, error) {
	if debug {
		log.Printf("Running: %s/%s %s", dir, command, arg)
	}
	cmd := exec.Command(command, arg...)

	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	check(err)
	stderr, err := cmd.StderrPipe()
	check(err)
	err = cmd.Start()
	check(err)

	in := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for in.Scan() {
		if debug {
			log.Printf(in.Text())
		}
	}

	err = cmd.Wait()
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
