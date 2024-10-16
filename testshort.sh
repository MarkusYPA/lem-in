#!/bin/bash

# Shorter good test cases
cases=("example00.txt" "example01.txt" "example02.txt" "example03.txt" "example04.txt" "example05.txt")

# Looping through them
for case in "${cases[@]}"
do
    path="testcases/$case"

    echo
    echo "Test case: $case"
    echo
    go run . "$path"
    echo "-----------------------------"
done
