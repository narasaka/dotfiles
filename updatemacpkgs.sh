#!/bin/sh

formula=$(brew list --formula)
casks=$(brew list --cask)

echo $formula >> formula.txt
echo $casks >> casks.txt
