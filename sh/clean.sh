#!/bin/bash
for x in bbchain{1..27}
do
  echo $x
  ssh $x rm -r logs
done
