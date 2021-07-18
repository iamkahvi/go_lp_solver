#!/bin/bash

cd ./test_LPs/input/

# go --version

for file in vanderbei*.txt
do
    printf "TESTING $file \n"
    go run ../../main.go < $file 2> /dev/null | colordiff ../output/$file -
    printf "\n"
done

# cd ../../test_LPs_volume2/input/

# # go --version

# for file in *.txt
# do
#     printf "TESTING $file \n"
#     go run ../../main.go < $file 2> /dev/null | colordiff ../output/$file -
#     printf "\n"
# done