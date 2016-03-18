package repoclone

import (
	"bufio"
    "io"
	"log"
	"os"
	"os/exec"
)

// todo: move to shared space
//func run(cmd *exec.Cmd, debug bool) (int, error) {
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
		log.Println(err)
		return 1, err
	}
	return 0, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
