#!/bin/bash

for x in bbchain{16..27}
do
  ssh $x mkdir -p logs
done

#1
echo bbchain16
ssh bbchain16 "./sh/client.sh -n 3 -l true"

#2
for x in bbchain{17..18}
  do echo $x
  ssh $x "./sh/client.sh -n 3 -w true"
done

#9
for x in bbchain{19..27}
  do echo $x
  ssh $x "./sh/client.sh -n 3 -r true"
done
