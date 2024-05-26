#!/bin/bash
for x in bbchain{1..15}
do
  ssh $x mkdir -p logs
done

echo bbchain1
ssh bbchain1 "nohup ./bin/benchmarkreplica 1> /dev/null 2> /dev/null &"

for x in bbchain{2..15}
  do echo $x
  ssh $x "nohup ./bin/benchmarkreplica --known-address bbchain1:8080 2>> /dev/null >> /dev/null &"
done
