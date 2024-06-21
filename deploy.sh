#!/usr/bin/env bash

git pull --rebase origin main
go build .
pkill main
nohup ./main &
