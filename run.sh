#! /bin/bash

num=5
for (( i=0; i < $num; i++ )); do
    go run main.go $i &
done
