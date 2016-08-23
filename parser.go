package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
)

const (
	active     = "a"
	create     = "c"
	horizontal = "h"
	vertical   = "v"
)

type Args []string

func Parse(session string, file *os.File) (cmds []Args) {
	windows := Args{}
	scanner := bufio.NewScanner(file)

	cmds = append(cmds, setupCommands(session))
	for scanner.Scan() {
		line := scanner.Text()
		captures, err := ParseLine(line)

		if err != nil {
			log.Printf("Skipped: '%v'. Err: %v. Got: %+v", line, err, captures)
			continue
		}

		window := captures["window"]
		operation := captures["operation"]
		target := captures["target"]
		action := captures["action"]

		if len(operation) > 1 {
			target = operation[0:(len(operation) - 1)]
			operation = operation[len(operation)-1:]
		}

		switch operation {
		case create:
			if contains(windows, window) {
				log.Fatal("Already created window %s", window)
			}

			windows = append(windows, window)
			cmds = append(cmds, newWindow(session, window))
		case horizontal, vertical:
			if !contains(windows, window) {
				log.Fatalf("need to create the window first before splitting it: %v", line)
			}

			cmds = append(cmds, splitWindow(session, window, tmuxSplit(operation), target))
		case active:
			cmds = append(cmds, selectWindow(session, window))

			if target != "" {
				cmds = append(cmds, selectPane(session, target))
			}
		}

		if action != "" {
			cmds = append(cmds, sendKeys(session, window, action))
		}
	}

	return
}

func contains(arr Args, v1 string) bool {
	for _, v2 := range arr {
		if v1 == v2 {
			return true
		}
	}

	return false
}

func newWindow(session, window string) Args {
	return Args{"new-window", "-a", "-t", session, "-n", window, "-c", os.Getenv("PWD")}
}

func selectWindow(session, window string) Args {
	target := fmt.Sprintf("%s:%s", session, window)
	return Args{"select-window", "-t", target}
}

func selectPane(session, target string) Args {
	pane := fmt.Sprintf("%s.%s", session, target)
	return Args{"select-pane", "-t", pane}
}

func splitWindow(session, window, split, target string) Args {
	if target != "" {
		target = fmt.Sprintf(".%s", target)
	}

	target = fmt.Sprintf("%s:%s%s", session, window, target)
	return Args{"split-window", "-t", target, fmt.Sprintf("-%s", split)}
}

func sendKeys(session, window, action string) Args {
	target := fmt.Sprintf("%s:%s", session, window)
	return Args{"send-keys", "-t", target, action, "C-m"}
}

func tmuxSplit(operation string) string {
	if operation == horizontal {
		return "v"
	}

	return "h"
}

func setupCommands(session string) Args {
	return Args{"new-session", "-d", "-s", session}
}

func killCommands(session string) Args {
	return Args{"kill-session", "-t", session}
}

func ParseLine(line string) (map[string]string, error) {
	var re = regexp.MustCompile(`(?P<window>[a-z]*)\s(?P<target>\d?)(?P<operation>[cvha])\s?(?P<action>.*)`)
	match := re.FindStringSubmatch(line)
	captures := make(map[string]string)

	if match == nil {
		return captures, errors.New("Can't parse line")
	}

	for i, name := range re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]
	}

	return captures, nil
}
