my configs managed by [stow](https://www.gnu.org/software/stow/manual/stow.html)

# dotfiles

Personal dotfiles for macOS development environment.

## Quick Start

```bash
git clone https://github.com/YOUR_USERNAME/dotfiles.git ~/dotfiles
cd ~/dotfiles
./fresh-install.sh
```

**📖 For detailed setup instructions, see [SETUP_GUIDE.md](SETUP_GUIDE.md)**

## What's Included

- **vim**: [lunarvim](https://www.lunarvim.org/) + neovim configs
- **shell**: [oh-my-zsh](https://ohmyz.sh/) (zsh) with custom aliases and functions
- **tmux**: custom configuration with Catppuccin theme
- **terminals**: Ghostty, Wezterm
- **window management**: yabai, sketchybar, karabiner

## Components

- `git/` - Git configuration (⚠️ Update with your name/email)
- `zsh/` - Zsh and Oh-My-Zsh configuration
- `tmux/` - Tmux configuration and plugins
- `vim/` - Vim, Neovim, and LunarVim configs
- `ghostty/` - Ghostty terminal emulator config
- `wezterm/` - Wezterm terminal emulator config
- `karabiner/` - Keyboard customization (Karabiner-Elements)
- `sketchybar/` - macOS menu bar replacement

## Installation Scripts

- `fresh-install.sh` - Complete setup for new Mac (all-in-one)
- `mac_primer.sh` - macOS system preferences
- `install_mac_pkgs.sh` - Install all packages via Homebrew
- `install.sh` - Symlink dotfiles using Stow

## Management Scripts

- `updatemacpkgs.sh` - Update package lists from current system
- `update-fresh-install.sh` - Regenerate fresh-install.sh

