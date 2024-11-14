#!/bin/bash

# Check if an argument was provided
if [ -z "$1" ]; then
  echo "Specify a text file in the testcases folder"
  echo "Usage: $0 <filename>"
  echo "For instance: $0 example05"
  exit 1
fi

# Run a container with port 8080 and pipe input:
go run . testcases/$1.txt | docker run -i -p 8080:8080 leminvisualizer
