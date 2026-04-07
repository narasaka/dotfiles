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

## mac

```
./fresh-install.sh
```

## applying dotfiles only

```
./install.sh
```

this runs `stow */ --no-folding` to symlink everything into `~`.

