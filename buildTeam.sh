#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o /mnt/d/Repos/greatwhite/output/TeamServer/TeamServerLin  -ldflags="-w -s" /mnt/d/Repos/greatwhite/cmd/Teamserver/server.go
GOOS=darwin GOARCH=amd64 go build -o /mnt/d/Repos/greatwhite/output/TeamServer/TeamServerMac  -ldflags="-w -s" /mnt/d/Repos/greatwhite/cmd/Teamserver/server.go
