#!/bin/bash
for x in bbchain{1..27}
do
  echo $x
  rsync -rP oncetree/bin $x:
done
