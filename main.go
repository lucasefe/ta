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
	defaultSession = path.Base(os.Getenv("PWD"))
	tmux = "/usr/local/bin/tmux"
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
		err := runTmuxCommand(arguments)
		if err != nil {
			log.Printf("args: %+v, err: %v\n", arguments, err)
		}
	}

	if !dryrun {
		cmd := exec.Command(tmux, ta.Args{"attach-session", "-t", session}...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func runTmuxCommand(arguments ta.Args) error {
	if dryrun {
		fmt.Printf("%s %s\n", tmux, strings.Join(arguments, " "))
		return nil
	}

	return exec.Command(tmux, arguments...).Run()
}
