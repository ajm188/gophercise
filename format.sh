#!/bin/sh
find . -name "*.go" | while read fname; do go fmt "$fname"; done
unset $fname
