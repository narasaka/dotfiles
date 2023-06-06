#!/bin/sh

brew install $(cat formula.txt)
brew install --cask $(cat casks.txt)
