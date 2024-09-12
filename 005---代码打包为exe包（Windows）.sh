#!/bin/bash

go build -ldflags="-H windowsgui" -o ./cmdb_mini_agent.exe ./main.go

chmod +x ./cmdb_mini_agent.exe
