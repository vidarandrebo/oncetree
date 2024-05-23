#!/bin/bash
for x in bbchain{19..27}
do
  ssh $x mkdir -p logs
done

echo bbchain19
ssh bbchain19 "nohup ./bin/benchmarkreplica 1> /dev/null 2> /dev/null &"

for x in bbchain{20..27}
  do echo $x
  ssh $x "nohup ./bin/benchmarkreplica --known-address bbchain19:8080 2>> /dev/null >> /dev/null &"
done
