#!/bin/sh

# homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

# dock: small, on the left, hides instantly
defaults write com.apple.dock orientation -string left
defaults write com.apple.dock tilesize -int 26
defaults write com.apple.dock autohide -bool true
defaults write com.apple.dock autohide-delay -float 0
defaults write com.apple.dock autohide-time-modifier -int 0
killall Dock
