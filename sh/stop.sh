#!/bin/bash
for x in bbchain{1..15}
do
  echo $x
  ssh $x killall benchmarkreplica
done

for x in bbchain{16..27}
do
  echo $x
  ssh $x killall benchmarkclient
done
