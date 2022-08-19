#!/bin/sh

yay -Qe | awk '{print $1}' > packages.txt
