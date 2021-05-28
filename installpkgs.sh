#!/bin/sh

sudo pacman -S --needed base-devel git
yay -S --needed - < packages.txt
