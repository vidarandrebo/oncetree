for x in bbchain{19..27}; do echo $x && rsync -rP bin $x:; done