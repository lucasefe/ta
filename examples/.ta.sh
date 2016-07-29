#!/usr/bin/env bash
if [[ -z $(tmux ls -F test 2>/dev/null | grep "^test$") ]]; then
  tmux $TMUX_OPTS new-session -s test -d
fi

tmux kill-session -t test
tmux new-session -s test -d 

# UNO
tmux new-window -a -t test -n uno -c $PWD
tmux send-keys -t test:uno "echo izq solo" C-m

tmux split-window -t test:uno -h
tmux send-keys -t test:uno "echo der arriba" C-m

tmux split-window -t test:uno -v
tmux send-keys -t test:uno "echo der abajo" C-m

# TRES
tmux new-window -a -t test -n tres -c $PWD
tmux send-keys -t test:tres "echo izq arriba" C-m

tmux split-window -t test:tres -h
tmux send-keys -t test:tres "echo der arriba" C-m

tmux split-window -t test:tres -v
tmux send-keys -t test:tres "echo der abajo" C-m

tmux split-window -t test:tres.1 -v
tmux send-keys -t test:tres "echo izq abajo" C-m

# cleanup
tmux kill-window -t test:$( tmux list-windows -t test -F "1" | head -n 1 )
tmux attach-session -t test
