#! /bin/bash

num=2
for (( i=0; i < $num; i++ )); do
    go run main.go $i &
done
