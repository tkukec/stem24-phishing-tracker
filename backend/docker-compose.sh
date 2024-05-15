#!/bin/bash

if [[ $1 == "up" ]]; then
    eval `ssh-agent`
    ssh-add
    docker build --ssh default . -t hrzg1do20.hr.asseco-see.local/voice/agent-management:dev
    kill $SSH_AGENT_PID
    docker-compose up -d
elif [[ $1 == "down" ]]; then
    docker-compose down
    docker rmi hrzg1do20.hr.asseco-see.local/voice/agent-management:dev
else
    echo "Please use up or down argument"
fi