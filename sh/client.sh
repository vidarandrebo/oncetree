#!/bin/bash

for x in bbchain{16..27}
do
  ssh $x mkdir -p logs
done

echo bbchain16
ssh bbchain16 "nohup ./bin/benchmarkclient --known-address bbchain1:8080 --node-to-crash-address bbchain5:8080 --writer --leader 1>> /dev/null 2>> /dev/null &"

for x in bbchain{17..18}
  do echo $x
  ssh $x "nohup ./bin/benchmarkclient --known-address bbchain1:8080 --node-to-crash-address bbchain5:8080 --writer 1> /dev/null 2> /dev/null &"
done
for x in bbchain{19..27}
  do echo $x
  ssh $x "nohup ./bin/benchmarkclient --known-address bbchain1:8080 --node-to-crash-address bbchain5:8080 --reader 1> /dev/null 2> /dev/null &"
done
