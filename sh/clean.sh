for x in bbchain{1..27}; do echo $x && ssh $x rm -r bin main start.sh; done
