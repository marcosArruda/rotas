#!/usr/bin/env bash
OOS=windows GOARCH=amd64 go build -o rotas-win.exe
GOOS=linux GOARCH=amd64 go build -o rotas-debian
GOOS=darwin GOARCH=amd64 go build -o rotas-macos
