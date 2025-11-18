# 📋 New Mac Setup Checklist

Use this checklist when setting up a new MacBook with these dotfiles.

## Before You Start

- [ ] Backup any existing configurations from the new Mac
- [ ] Ensure you have admin access
- [ ] Connect to reliable internet

## Quick Setup (Recommended Path)

### Step 1: Clone Repository
```bash
git clone https://github.com/YOUR_USERNAME/dotfiles.git ~/dotfiles
cd ~/dotfiles
```

- [ ] Repository cloned successfully

### Step 2: Configure Personal Information
```bash
vim git/.gitconfig
```

Update these fields:
- [ ] `name = YOUR_NAME_HERE` → Your actual name
- [ ] `email = YOUR_EMAIL_HERE` → Your actual email

### Step 3: Run Fresh Install Script
```bash
./fresh-install.sh
```

This will take 15-30 minutes depending on your internet speed.

- [ ] Script completed without errors
- [ ] Homebrew installed
- [ ] All packages installed
- [ ] Configurations symlinked

### Step 4: Install Oh-My-Zsh
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

- [ ] Oh-My-Zsh installed
- [ ] Shell switched to zsh (if not already)

### Step 5: Reload Shell
```bash
source ~/.zshrc
```

- [ ] Shell reloaded without errors
- [ ] Zsh prompt showing correctly

### Step 6: Set Up Tmux Plugins
```bash
tmux
# Inside tmux, press: Ctrl+a then I (capital i)
```

- [ ] Tmux launched successfully
- [ ] Plugins installed (Ctrl+a I)
- [ ] Catppuccin theme applied

## Optional Setup

### Neovim/LunarVim (if using)
```bash
# For LunarVim
LV_BRANCH='release-1.4/neovim-0.9' bash <(curl -s https://raw.githubusercontent.com/LunarVim/LunarVim/release-1.4/neovim-0.9/utils/installer/install.sh)

# For regular Neovim with Packer
git clone --depth 1 https://github.com/wbthomason/packer.nvim \
  ~/.local/share/nvim/site/pack/packer/start/packer.nvim
```

- [ ] LunarVim installed (if using)
- [ ] Packer installed (if using regular nvim)
- [ ] Plugins synced (`:PackerSync` in nvim)

### Secrets Configuration (if needed)
```bash
# For Advent of Code scripts
echo 'export AOC_SESSION_ID="your_session_id"' > ~/.scripts/secret
```

- [ ] Secret file created and configured

### Google Cloud SDK (if needed)
```bash
# Download from: https://cloud.google.com/sdk/docs/install
# Or use: brew install --cask google-cloud-sdk
```

- [ ] Google Cloud SDK installed (if needed)
- [ ] Path configured in .zshrc

## Verification Steps

### Test Basic Tools
```bash
# Test zsh
echo $SHELL  # Should show /bin/zsh or /usr/local/bin/zsh

# Test git config
git config --get user.name
git config --get user.email

# Test tmux
tmux -V

# Test neovim
nvim --version

# Test common CLI tools
fzf --version
rg --version
bat --version
```

- [ ] All basic tools working

### Test Key Bindings

#### Tmux
- [ ] `Ctrl+a` works as prefix
- [ ] `Ctrl+a |` splits vertically
- [ ] `Ctrl+a -` splits horizontally
- [ ] `Ctrl+a r` reloads config

#### Vim
- [ ] `Ctrl+\` toggles NERDTree (in vim)
- [ ] Relative line numbers working
- [ ] System clipboard integration working

### Test Application Launchers
- [ ] Ghostty terminal installed and configured
- [ ] Cursor IDE installed
- [ ] Ghostty quick terminal: `Ctrl+\`` works

## Customization (Optional)

### macOS Settings
```bash
# View all system preferences
defaults read > ~/macos_settings_before.txt
```

Consider customizing:
- [ ] Dock position and size
- [ ] Mission Control settings
- [ ] Trackpad settings
- [ ] Keyboard repeat rate

### Additional Apps
Install any apps not in casks.txt:
- [ ] Additional browsers
- [ ] Design tools
- [ ] Entertainment apps

## Post-Setup

### Update Package Lists
After installing additional packages:
```bash
./updatemacpkgs.sh
git add formula.txt casks.txt
git commit -m "Update package lists"
git push
```

- [ ] Package lists updated and committed

### Test Everything
- [ ] Open a new terminal window
- [ ] Verify all environment variables
- [ ] Test project setups (Node, Python, Go, etc.)
- [ ] Verify all key applications launch

### Backup
- [ ] Time Machine configured
- [ ] Important files backed up to cloud
- [ ] SSH keys generated/copied
- [ ] GPG keys configured (if using)

## Troubleshooting Reference

### If Stow Fails
```bash
# Check for conflicts
stow -n -v */ --no-folding

# Backup and remove conflicts
mkdir ~/dotfiles_backup
mv ~/.zshrc ~/dotfiles_backup/
# ... repeat for other files

# Try again
./install.sh
```

### If Homebrew Command Not Found
```bash
# Apple Silicon
eval "$(/opt/homebrew/bin/brew shellenv)"

# Intel
eval "$(/usr/local/bin/brew shellenv)"
```

### If Tmux Plugins Don't Load
```bash
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
tmux source ~/.tmux.conf
# Then: Ctrl+a I
```

## Completion

When all boxes are checked:
- [ ] All configurations working correctly
- [ ] All development tools functional
- [ ] Personal information updated
- [ ] No errors in shell startup

### 🎉 Setup Complete!

Your Mac is now fully configured and ready for development.

**Next Steps:**
1. Star this repo for easy access
2. Set up your development projects
3. Install project-specific tools as needed

**Keep This Repo Updated:**
- Commit new configurations when you make changes
- Update package lists regularly with `./updatemacpkgs.sh`
- Document any manual setup steps you discover

---

**Setup Time Estimate:** 45-60 minutes
**Last Updated:** 2025-11-18


