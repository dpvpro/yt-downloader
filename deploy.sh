#!/usr/bin/env bash

git pull --rebase origin main
go build -o yt-downloader
pkill yt-downloader
nohup ./yt-downloader &
