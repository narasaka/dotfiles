#!/bin/sh

formula=$(brew list --formula --full-name)
casks=$(brew list --cask --full-name)

echo $formula > formula.txt
echo $casks > casks.txt
