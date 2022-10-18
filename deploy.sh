#!/usr/bin/env bash

git pull --rebase origin main
go build main.go
pkill main
nohup ./main &