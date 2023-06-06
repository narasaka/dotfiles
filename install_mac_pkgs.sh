#!/bin/sh

brew install $(cat formula.txt)
brew install --cask $(cat casks.txt)
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
