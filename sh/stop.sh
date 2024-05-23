#!/bin/bash
for x in bbchain{19..27}; do ssh $x killall benchmarkreplica; done
