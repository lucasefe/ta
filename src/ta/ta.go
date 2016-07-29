package ta

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	create     = "n"
	horizontal = "h"
	vertical   = "v"
)

func Parse(session string, file *os.File) (cmds []string) {
	windows := []string{}
	scanner := bufio.NewScanner(file)

	cmds = append(cmds, setupCommands(session))
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.SplitN(line, " ", 3)

		if len(arr) != 3 {
			continue
		}

		window := arr[0]
		action := arr[2]
		operation := arr[1]
		target := ""

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
			cmds = append(cmds, sendKeys(session, window, action))

		case horizontal, vertical:
			if !contains(windows, window) {
				log.Fatalf("need to create the window first before splitting it: %v", line)
			}

			cmds = append(cmds, splitWindow(session, window, tmuxSplit(operation), target))
			cmds = append(cmds, sendKeys(session, window, action))
		}

	}

	cmds = append(cmds, cleanupCommands(session))
	return
}

func contains(arr []string, v1 string) bool {
	for _, v2 := range arr {
		if v1 == v2 {
			return true
		}
	}

	return false
}

func newWindow(session, window string) string {
	return fmt.Sprintf("tmux new-window -a -t %s -n %s -c $PWD", session, window)
}

func splitWindow(session, window, split, target string) string {
	if target != "" {
		target = fmt.Sprintf(".%s", target)
	}

	return fmt.Sprintf("tmux split-window -t %s:%s%s -%s", session, window, target, split)
}

func sendKeys(session, window, action string) string {
	return fmt.Sprintf("tmux send-keys -t %s:%s \"%s\" C-m", session, window, action)
}

func tmuxSplit(operation string) string {
	if operation == horizontal {
		return "v"
	}

	return "h"
}

func setupCommands(session string) string {
	return fmt.Sprintf("tmux new-session -s %s -d", session)
}

func killCommands(session string) string {
	return fmt.Sprintf("tmux kill-session -t %s ", session)
}

func cleanupCommands(session string) string {
	return fmt.Sprintf("tmux kill-window -t %s:$( tmux list-windows -t %s -F \"1\" | head -n 1 );  tmux attach-session -t %s", session, session, session)
}
