package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	config         string
	session        string
	defaultSession string
	dryrun         bool

	tmux string
)

func main() {
	tmux = "/usr/local/bin/tmux"
	defaultSession = path.Base(os.Getenv("PWD"))
	flag.StringVar(&config, "f", ".ta", "the ta config file")
	flag.StringVar(&session, "s", defaultSession, "session name")
	flag.BoolVar(&dryrun, "d", false, "dry run mode")
	flag.Parse()

	file, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if existingSession(session) {
		attachToSession(session)
		return
	}

	for _, arguments := range Parse(session, file) {
		if dryrun {
			fmt.Printf("%s %s\n", tmux, strings.Join(arguments, " "))
		} else {
			err := exec.Command(tmux, arguments...).Run()
			if err != nil {
				log.Printf("args: %+v, err: %v\n", arguments, err)
			}
		}
	}

	if !dryrun {
		cleanup(session)
		attachToSession(session)
	}
}

func existingSession(session string) bool {
	cmd := exec.Command(tmux, Args{"ls", "-F", "\"#{session_name}\""}...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatalf("OO: %v", err)
	}

	cmd.Start()

	bytes, _ := ioutil.ReadAll(stdout)

	if err := cmd.Wait(); err != nil {
		return false
	}

	fmt.Printf("%+v\n", bytes)
	sessionName := fmt.Sprintf("\"%s\"", session)
	for _, line := range strings.Split(string(bytes), "\n") {
		if sessionName == line {
			return true
		}
	}

	return false
}

func attachToSession(session string) {
	cmd := exec.Command(tmux, Args{"attach-session", "-t", session}...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func cleanup(session string) {
	action := "tmux kill-window -t 1"
	exec.Command(tmux, sendKeys(session, "1", action)...).Run()
}
