package repoclone

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

var debug bool

// todo: move to shared space
//func run(cmd *exec.Cmd, debug bool) (int, error) {
func run(dir string, command string, arg ...string) (int, error) {
	// app, err := exec.LookPath(command)
	// check(err)
	// cmd := exec.Command(app, arg)
	// cmd.Dir = "./"
	cmd := exec.Command(command, arg...)

	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	check(err)
	err = cmd.Start()
	check(err)
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		if debug {
			log.Printf(in.Text())
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
		return 1, err
	}
	return 0, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
