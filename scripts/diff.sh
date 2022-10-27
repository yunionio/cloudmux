#!/bin/bash

target_dir=$1
src=$2

#diff -u "$src" "$target_dir/$src" | nvim -R -
colordiff -u "$target_dir/$src" "$src" | less
