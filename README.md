# dotfiles

my configs managed by [stow](https://www.gnu.org/software/stow/manual/stow.html).

## what's in here

- zsh: [oh-my-zsh](https://ohmyz.sh/) with vi-mode
- vim: neovim config with lazy.nvim
- git: global gitconfig
- zellij: terminal multiplexer config
- opencode: opencode editor config + custom commands
- tmux: tmux config

## ubuntu fresh install

sets up a full dev environment on ubuntu 24.04 from scratch.

```
bash -c "$(curl -fsSL https://raw.githubusercontent.com/narasaka/dotfiles/master/setup-ubuntu.sh)"
```

this installs everything: apt packages, external repos (docker, tailscale, gcloud, etc.), go, rust, node (fnm), bun, neovim, cargo/go tools, oh-my-zsh, and applies all dotfiles via stow.

after it finishes, log out and back in, then run:

```
tailscale up
op signin
op-ssh-load
gh auth login
gcloud auth login
doppler login
```

the script is idempotent -- safe to re-run.

## mac fresh install

sets up a mac from scratch. pick `dev` (what i need to work) or `full` (everything):

```
bash -c "$(curl -fsSL https://raw.githubusercontent.com/narasaka/dotfiles/master/setup-mac.sh)" -- dev
bash -c "$(curl -fsSL https://raw.githubusercontent.com/narasaka/dotfiles/master/setup-mac.sh)" -- full
```

leave off the mode and it asks.

- `dev`: homebrew, everything in `Brewfile` (cli tools, languages, ghostty, 1password, docker, tailscale, gcloud, fonts), rust (rustup), node lts (fnm), go/cargo/bun/npm globals, claude code, oh-my-zsh, tmux + nvim plugins, macos dock/dark mode defaults, and applies all dotfiles via stow
- `full`: dev + everything in `Brewfile.full` (other languages, databases, gui apps, misc tools)

if a file would block a stow symlink, it gets moved to `~/.dotfiles-backup/<timestamp>/` first. to add a package, edit `Brewfile` or `Brewfile.full`.

apps already installed outside homebrew (mdm, app store, manual download) are skipped. a failing step doesn't stop the script -- it keeps going and lists what failed at the end.

after it finishes, open 1password, tailscale, docker and karabiner once to sign in and grant permissions, then run:

```
gh auth login
gcloud auth login
doppler login
```

the script is idempotent -- safe to re-run.

## applying dotfiles only

```
./install.sh
```

this runs `stow */ --no-folding` to symlink everything into `~`.

