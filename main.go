package main

import (
	"fmt"
	"log"
	"os"

	"ta"
)

func main() {
	session := "test"
	config := "examples/.ta"
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
