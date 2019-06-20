#!/bin/bash

GOCACHE=/root/.gocache GOPATH=/root/go /usr/local/go/bin/go build -o whu
systemctl stop whu
mv whu /home/webmaster/bin/whu
systemctl start whu