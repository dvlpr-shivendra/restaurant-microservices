#!/bin/bash

# Create a new tmux session named 'restaurant-backend'
SESSION="restaurant-backend"
tmux new-session -d -s $SESSION

# Create a new window for Services
tmux new-window -t $SESSION -n 'services'

# Pane 1: Gateway service
tmux send-keys -t $SESSION:1 'cd gateway && air' C-m

# Split the window horizontally for the Orders service
tmux split-window -h -t $SESSION:1
tmux send-keys -t $SESSION:1.1 'cd orders && air' C-m

# Split the right pane vertically for the Payments service
tmux split-window -v -t $SESSION:1.1
tmux send-keys -t $SESSION:1.2 'cd payments && air' C-m

# Switch back to the 'services' window after starting Docker
tmux select-window -t $SESSION:1

# Select the first pane (Gateway) in the services window and attach to the tmux session
tmux select-pane -t $SESSION:1.0
tmux attach -t $SESSION
