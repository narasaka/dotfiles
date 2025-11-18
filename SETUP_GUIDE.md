# 🚀 Dotfiles Setup Guide for New MacBook

This guide will walk you through setting up your new MacBook with all your configurations using these dotfiles.

## 📋 Prerequisites

Before starting, ensure you have:
- A fresh macOS installation
- Admin access to the machine
- Internet connection

## 🔧 Quick Setup (Recommended)

For a completely fresh Mac, run the all-in-one setup script:

```bash
# Clone this repository
git clone https://github.com/YOUR_USERNAME/dotfiles.git ~/dotfiles
cd ~/dotfiles

# Run the complete fresh install script
./fresh-install.sh
```

This will:
1. Install Homebrew
2. Configure macOS dock settings (faster animations)
3. Install all Homebrew formulas and casks
4. Install tmux plugin manager (tpm)
5. Symlink all configuration files using GNU Stow

## 🎯 Step-by-Step Setup (Manual)

If you prefer to install components individually:

### 1. Install Homebrew

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

### 2. Clone This Repository

```bash
git clone https://github.com/YOUR_USERNAME/dotfiles.git ~/code_dev/dotfiles
cd ~/code_dev/dotfiles
```

### 3. Configure macOS Settings

```bash
./mac_primer.sh
```

This sets up:
- Faster dock auto-hide animations
- Other macOS preferences

### 4. Install Packages

```bash
./install_mac_pkgs.sh
```

This installs:
- All Homebrew formulas (see `formula.txt`)
- All Homebrew casks (see `casks.txt`)
- Tmux Plugin Manager (tpm)

### 5. Symlink Configuration Files

```bash
./install.sh
```

This uses [GNU Stow](https://www.gnu.org/software/stow/) to symlink all dotfiles to your home directory.

## ⚙️ Required Manual Configuration

After running the setup scripts, you **MUST** configure these personal settings:

### 1. Git Configuration

Edit `git/.gitconfig` and update with your information:

```bash
vim git/.gitconfig
```

Replace:
- `YOUR_NAME_HERE` with your full name
- `YOUR_EMAIL_HERE` with your email address

### 2. Oh-My-Zsh Installation

The `.zshrc` expects Oh-My-Zsh to be installed:

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

### 3. Reload Shell Configuration

After symlinking, reload your shell:

```bash
source ~/.zshrc
```

### 4. Optional: Configure Secrets

If you use Advent of Code scripts, create a secrets file:

```bash
echo 'export AOC_SESSION_ID="your_session_id_here"' > ~/.scripts/secret
```

### 5. Tmux Plugins

After first launching tmux, install plugins:

```bash
tmux
# Inside tmux, press: Ctrl+a then I (capital i)
```

### 6. Neovim/LunarVim Setup

If using LunarVim, install it:

```bash
LV_BRANCH='release-1.4/neovim-0.9' bash <(curl -s https://raw.githubusercontent.com/LunarVim/LunarVim/release-1.4/neovim-0.9/utils/installer/install.sh)
```

For regular Neovim, install Packer (plugin manager):

```bash
git clone --depth 1 https://github.com/wbthomason/packer.nvim \
  ~/.local/share/nvim/site/pack/packer/start/packer.nvim
```

Then open nvim and run:
```
:PackerSync
```

## 📦 What Gets Installed

### Development Tools
- **Languages**: Go, Rust, Erlang, Gleam, Zig, Python, Node.js
- **Version Managers**: fnm (Node), nvm (backup)
- **Databases**: PostgreSQL, Redis
- **CLI Tools**: fzf, ripgrep, bat, fd, zoxide, lazygit, btop, htop
- **Cloud**: AWS SDK, Google Cloud SDK, Doppler, Railway, Vercel

### Applications (via Cask)
- **Terminals**: Ghostty, Wezterm
- **Editors**: Cursor, VSCodium, Neovim
- **Dev Tools**: Docker, Beekeeper Studio, Postman
- **Communication**: Zoom, Discord (Legcord)
- **Utilities**: Tailscale, CloudFlare WARP, AnyDesk, RustDesk

### Window Management
- **yabai**: Tiling window manager
- **sketchybar**: Custom status bar (configs included)
- **karabiner**: Keyboard customization

## 🔄 Keeping Things Updated

### Update Package Lists

After installing new packages, update your package lists:

```bash
# For macOS
./updatemacpkgs.sh

# For Arch Linux (if dual-booting)
./updatepkgs.sh
```

### Regenerate Fresh Install Script

After making changes to installation scripts:

```bash
./update-fresh-install.sh
```

## 📂 Directory Structure

```
dotfiles/
├── git/                  # Git configuration
│   └── .gitconfig
├── zsh/                  # Zsh configuration
│   ├── .zshrc
│   └── .scripts/         # Custom shell scripts
├── tmux/                 # Tmux configuration
│   └── .tmux.conf
├── vim/                  # Vim/Neovim configuration
│   ├── .vimrc
│   └── .config/
│       ├── nvim/
│       └── lvim/
├── ghostty/              # Ghostty terminal config
│   └── .config/ghostty/
├── wezterm/              # Wezterm terminal config
│   └── .wezterm.lua
├── karabiner/            # Keyboard customization
│   └── .config/karabiner/
├── sketchybar/           # macOS status bar
│   └── .config/sketchybar/
└── wallpaper/            # Desktop wallpapers
```

## 🔑 Key Features

### ZSH Configuration
- **Theme**: dpoggi
- **Plugins**: git, web-search, vi-mode
- **Enhancements**: 
  - zoxide (smart cd)
  - thefuck (command correction)
  - fnm (fast Node.js version manager)
  - Auto-activate Python virtual environments

### Tmux Configuration
- **Prefix**: `Ctrl+a` (instead of default `Ctrl+b`)
- **Mouse**: Enabled
- **Splitting**: `Ctrl+a |` (vertical), `Ctrl+a -` (horizontal)
- **Navigation**: Vim-style with tmux-navigator
- **Theme**: Catppuccin Mocha

### Vim/Neovim
- **Line numbers**: Relative line numbers
- **Clipboard**: System clipboard integration
- **Splits**: Open right and below
- **Theme**: Sonokai (Shusia variant)
- **Plugins**: NERDTree, vim-fugitive, vim-surround, prettier

## 🎨 Terminal Configuration

### Ghostty
- **Font**: FiraCode Nerd Font Mono (12pt)
- **Theme**: Jellybeans
- **Features**: Copy on select, hidden titlebar proxy icon
- **Global hotkey**: `Ctrl+\`` for quick terminal

### Wezterm
- **Font**: FiraCode Nerd Font (14pt)
- **Color Scheme**: Azu (Gogh)

## 🐛 Troubleshooting

### Stow Conflicts

If stow reports conflicts (existing files), you can:

```bash
# Backup existing configs
mkdir ~/dotfiles_backup
mv ~/.zshrc ~/dotfiles_backup/
mv ~/.tmux.conf ~/dotfiles_backup/
# ... etc

# Then run stow again
./install.sh
```

### Oh-My-Zsh Not Found

Make sure to install Oh-My-Zsh before using the .zshrc:

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

### Homebrew Not in PATH

Add to your shell profile:

```bash
# For Apple Silicon Macs
eval "$(/opt/homebrew/bin/brew shellenv)"

# For Intel Macs
eval "$(/usr/local/bin/brew shellenv)"
```

### Tmux Plugin Manager Not Working

Ensure tpm is installed:

```bash
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
```

Then reload tmux config:
```bash
tmux source ~/.tmux.conf
```

## 📝 Notes

- **Linux Support**: This repo includes Arch Linux packages (`packages.txt`) but the main focus is macOS
- **Auto-generated Files**: `packer_compiled.lua` has been removed as it contains machine-specific paths and will be auto-generated
- **Customization**: Feel free to modify any configs to match your preferences
- **Backups**: Always backup your existing dotfiles before running these scripts

## 🔐 Security

Remember to:
- Never commit secrets or API keys
- Use `~/.scripts/secret` for sensitive environment variables (this file is not tracked)
- Review all scripts before running them

## 📚 Resources

- [GNU Stow Manual](https://www.gnu.org/software/stow/manual/stow.html)
- [Oh My Zsh](https://ohmyz.sh/)
- [Tmux Plugin Manager](https://github.com/tmux-plugins/tpm)
- [LunarVim](https://www.lunarvim.org/)

## 🤝 Contributing

These are personal dotfiles, but feel free to fork and adapt for your own use!

---

**Last Updated**: 2025-11-18

Made with ❤️ for seamless macOS development environment setup


