#!/bin/bash

go build -ldflags='-s -w -extldflags "-static"' -o ./cmdb_mini_agent ./main.go

chmod +x ./cmdb_mini_agent
