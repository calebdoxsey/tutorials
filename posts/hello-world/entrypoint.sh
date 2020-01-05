#!/bin/bash

# echo "===[ V1 ]==="
# for n in 1000 2000 4000 8000 16000; do
#     echo "run asm v1 $n times"
#     time bash -c "seq $n | xargs -P1 -n1 /asm/v1 > /dev/null"

#     echo "run go v1 $n times"
#     time bash -c "seq $n | xargs -P1 -n1 /go/v1 > /dev/null"
# done
# echo "============"

# echo "===[ V2 ]==="
# for n in 1024000 2048000 4096000; do
#     echo "run asm v2 $n times"
#     time bash -c "/asm/v2 $n > /dev/null"

#     echo "run go v2 $n times"
#     time bash -c "/go/v2 $n > /dev/null"
# done
# echo "============"

# echo "===[ V3 ]==="
# for n in 1024000 2048000 4096000; do
#     echo "run go v3 $n times"
#     time bash -c "/go/v3 $n > /dev/null"
# done
# echo "============"

strace -c /asm/v2 4096000 >/dev/null
strace -c /go/v3 4096000 >/dev/null
