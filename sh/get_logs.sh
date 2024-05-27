#!/bin/bash

rm -r oncetree/replicas
mkdir oncetree/replicas
for x in bbchain{1..15}
do
  rsync $x:logs/* oncetree/replicas/
done

rm -r oncetree/writes
mkdir oncetree/writes
for x in bbchain{16..18}
do
  rsync $x:logs/* oncetree/writes/
done

rm -r oncetree/reads
mkdir oncetree/reads
for x in bbchain{19..27}
do
  rsync $x:logs/* oncetree/reads/
done
