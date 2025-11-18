#!/bin/bash


# homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

# faster dock
defaults write com.apple.dock autohide-delay -float 0; defaults write com.apple.dock autohide-time-modifier -int 0;killall Dock

# Commit: c164b73


brew install $(cat formula.txt)
brew install --cask $(cat casks.txt)
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm


stow */ --no-folding
