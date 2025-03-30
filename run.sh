#!/usr/bin/env bash

cd "$(dirname "$0")" || exit 1

git pull --rebase origin main

go build -o yt-downloader
pkill yt-downloader
./yt-downloader & disown
