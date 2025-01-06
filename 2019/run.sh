#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
FILE=$(realpath "$1")

if [[ ! -f "$FILE" ]]; then
    echo "File not exist"
    exit 1
fi

cd $SCRIPT_DIR/utils/machine/run
go run main.go -run "$FILE" "${@:2}"
cd $SCRIPT_DIR