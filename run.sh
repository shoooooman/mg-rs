#! /bin/bash
if [ $1 = "start" ]; then
    go build main.go

    num=50
    for (( i=0; i < $num; i++ )); do
        ./main $i --tag=mg-rs &
    done
    rm ./main
elif [ $1 = "kill" ]; then
    pgrep -f mg-rs | xargs kill
else
    echo please input valid command
fi
