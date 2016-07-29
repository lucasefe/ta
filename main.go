package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"

	"ta"
)

var tmux, config, session, defaultSession string

func main() {
	tmux = "/usr/local/bin/tmux"
	defaultSession = path.Base(os.Getenv("PWD"))
	flag.StringVar(&config, "f", ".ta", "the ta config file")
	flag.StringVar(&session, "s", defaultSession, "the ta config file")
	flag.Parse()

	file, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, arguments := range ta.Parse(session, file) {
		err := exec.Command(tmux, arguments...).Run()
		if err != nil {
			log.Println("args: %+v, err: %v", arguments, err)
		}
	}

	cmd := exec.Command(tmux, ta.Args{"attach-session", "-t", session}...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}
