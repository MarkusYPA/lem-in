#!/bin/bash

# Bad test cases
cases=("badexample00.txt" "badexample01.txt" "bad02.txt" "bad03.txt" "bad04.txt" "bad05.txt" "bad06.txt")

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
