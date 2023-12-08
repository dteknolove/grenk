#!/bin/bash

echo "grenk build..."
go build -o ./grenk ./cmd/grenk/main.go
wait
echo "move grenk binary to usr/go/bin or make your own path"



