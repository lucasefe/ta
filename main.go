package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"ta"
)

var config, session, defaultSession string

func main() {
	defaultSession = path.Base(os.Getenv("PWD"))
	flag.StringVar(&config, "f", ".ta", "the ta config file")
	flag.StringVar(&session, "s", defaultSession, "the ta config file")
	flag.Parse()

	file, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	commands := ta.Parse(session, file)
	for _, command := range commands {
		fmt.Println(command)
	}
}
