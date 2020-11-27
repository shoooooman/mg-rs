#! /bin/bash
i=0
for fn in $@; do
    python3 ./run/design.py $i
    bash ./run.sh start 2> ${fn}
    echo ${fn} done
    let i=$i+1
done
