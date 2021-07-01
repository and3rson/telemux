#!/bin/bash

echo "# Changelog"
echo

tags=`git tag | grep -v gormpersistence | sort -Vr`
tag=
for prev in $tags
do
    if [[ "$tag" == "" ]]
    then
        tag=$prev
        continue
    fi
    echo "## $tag"
    echo
    lines=`git log $prev..$tag --pretty="format:%b" | sed '/^$/d'`
    IFS=$'\n'
    for line in $lines
    do
        line=$(echo "$line" | sed -re "s/^([a-z_]+): ([a-z-]+)\b/\`\1\`: *\2*/g")
        echo "- $line"
    done
    echo
    tag=$prev
done
