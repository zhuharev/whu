#!/bin/bash

GOCACHE=/webmaster/.gocache GOPATH=/webmaster/go /usr/local/go/bin/go build -o whu
systemctl --user stop whu
mv whu /home/webmaster/bin/whu
systemctl --user start whu