#!/bin/bash

tags=`git tag | sort -Vr`
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
        echo "- $line"
    done
    echo
    tag=$prev
done
