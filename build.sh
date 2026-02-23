#!/usr/bin/env bash

APP=poecampain

case "$OSTYPE" in
    msys*|cygwin*)
        APP+=.exe
        ;;
esac

go build -ldflags='-s -w' -o dist/$APP src/*.go

cp -r data dist/
