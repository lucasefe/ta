package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"ta"
)

var (
	config         string
	session        string
	defaultSession string
	dryrun         bool
	tmux           string
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

	for _, arguments := range ta.Parse(session, file) {
		err := runTmuxCommand(arguments, dryrun)
		if err != nil {
			log.Printf("args: %+v, err: %v\n", arguments, err)
		}
	}
	cleanup(session, dryrun)
	attachToSession(session, dryrun)
}

func runTmuxCommand(arguments ta.Args, dryrun bool) error {
	if dryrun {
		fmt.Printf("%s %s\n", tmux, strings.Join(arguments, " "))
		return nil
	}

	return exec.Command(tmux, arguments...).Run()
}

func attachToSession(session string, dryrun bool) {
	if dryrun {
		return
	}

	cmd := exec.Command(tmux, ta.Args{"attach-session", "-t", session}...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func cleanup(session string, dryrun bool) {
	if dryrun {
		return
	}

	action := "tmux kill-window -t 1"
	target := fmt.Sprintf("%s:1", session)
	args := ta.Args{"send-keys", "-t", target, action, "C-m"}
	exec.Command(tmux, args...).Run()
}
