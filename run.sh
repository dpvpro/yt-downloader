#!/usr/bin/env bash

cd "$(dirname "$0")"

git pull --rebase origin main

go build -o yt-downloader
pkill yt-downloader
nohup ./yt-downloader &
