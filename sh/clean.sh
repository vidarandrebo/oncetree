#!/bin/bash
for x in bbchain{19..27}; do echo $x && ssh $x rm -r bin logs ; done
