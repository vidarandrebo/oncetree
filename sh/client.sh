#!/bin/bash
while getopts n:l:w:r: flag
do
    case "${flag}" in
        n) num_nodes=${OPTARG};;
        l) leader=${OPTARG};;
        r) reader=${OPTARG};;
        w) writer=${OPTARG};;
    esac
done
echo "num_nodes: $num_nodes";
echo "leader: $leader";

if [[ $leader == "true" ]]; then
    echo "create leader writer"
  nohup ./bin/benchmarkclient --known-address bbchain1:8080 --node-to-crash-address bbchain5:8080 --writer --leader 1>> /dev/null 2>> /dev/null &
  for ((i = 1; i < $num_nodes; i++))
  do
    echo "create writer"
    nohup ./bin/benchmarkclient --known-address bbchain1:8080 --node-to-crash-address bbchain5:8080 --writer 1>> /dev/null 2>> /dev/null &
  done
else
  if [[ $reader == "true" ]]; then
    for ((i = 0; i < $num_nodes; i++))
    do
      echo "create reader"
      nohup ./bin/benchmarkclient --known-address bbchain1:8080 --node-to-crash-address bbchain5:8080 --reader 1>> /dev/null 2>> /dev/null &
    done
  elif [[ $writer == "true" ]]; then
    for ((i = 0; i < $num_nodes; i++))
    do
      echo "create writer"
      nohup ./bin/benchmarkclient --known-address bbchain1:8080 --node-to-crash-address bbchain5:8080 --writer 1>> /dev/null 2>> /dev/null &
    done
  fi

fi

