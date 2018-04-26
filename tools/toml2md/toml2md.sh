#!/bin/sh

while read l; do
	case "$l" in
		\#*) echo ${l#?};;
		'') echo "$l";;
		*) echo "\`$l\`";;
	esac; 
done < config.commented.toml
