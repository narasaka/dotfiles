typeset -U path

# os detection
_IS_MAC=0 _IS_LINUX=0
case "$(uname -s)" in
  Darwin) _IS_MAC=1 ;;
  Linux)  _IS_LINUX=1 ;;
esac

# oh my zsh
export ZSH="$HOME/.oh-my-zsh"

if (( _IS_MAC )); then
  ZSH_THEME="dpoggi"
else
  ZSH_THEME="gentoo"
fi

plugins=(git web-search vi-mode)

# vi-mode
VI_MODE_RESET_PROMPT_ON_MODE_CHANGE=true
VI_MODE_SET_CURSOR=true
VI_MODE_CURSOR_NORMAL=2
VI_MODE_CURSOR_VISUAL=6
VI_MODE_CURSOR_INSERT=6
VI_MODE_CURSOR_OPPEND=0

source "$ZSH/oh-my-zsh.sh"

# environment
export EDITOR='nvim'
export TIMEFMT=$'\nreal\t%E\nuser\t%U\nsys\t%S'
export GPG_TTY=$(tty)

# ssh-agent (linux systemd user service)
if (( _IS_LINUX )) && [[ -n "$XDG_RUNTIME_DIR" ]]; then
  export SSH_AUTH_SOCK="$XDG_RUNTIME_DIR/ssh-agent.socket"
fi

# aliases
alias gg='g++ -std=c++17 -Wshadow -Wall -O2 -Wno-unused-result'
alias gf='g++ -std=c++17 -Wshadow -Wall -g -fsanitize=address -fsanitize=undefined -D_GLIBCXX_DEGUG'
alias cptemp='cp ~/prog/templates/temp.cpp solve.cpp'
alias pytemp='cp ~/prog/templates/temp.py solve.py'
alias vim='nvim'
alias dbui='nvim +DBUI'
alias clear='clear && clear'

# zellij: attach or create session named after current directory
zj() {
  local name="${1:-$(basename "$PWD")}"
  zellij attach "$name" -c
}

if (( _IS_MAC )); then
  alias python='python3'
  alias pypy='pypy3'
fi

# path
export PATH="$HOME/.local/bin:$PATH"

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# go
export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"

# pnpm
if (( _IS_MAC )); then
  export PNPM_HOME="$HOME/Library/pnpm"
else
  export PNPM_HOME="$HOME/.local/share/pnpm"
fi
[[ -d "$PNPM_HOME" ]] && export PATH="$PNPM_HOME:$PATH"

# fnm (homebrew on mac, ~/.local/share on linux)
FNM_PATH="$HOME/.local/share/fnm"
[[ -d "$FNM_PATH" ]] && export PATH="$FNM_PATH:$PATH"

# tool initialization (skipped if not installed)
command -v thefuck &>/dev/null && eval "$(thefuck --alias)"
command -v fnm     &>/dev/null && eval "$(fnm env --use-on-cd --shell zsh)"
command -v zoxide  &>/dev/null && eval "$(zoxide init zsh --cmd cd)"

# google cloud sdk (lazy load)
_gcloud_sdk="$HOME/Downloads/google-cloud-sdk"
if [[ -f "$_gcloud_sdk/path.zsh.inc" ]]; then
  function _lazy_gcloud() {
    source "$_gcloud_sdk/path.zsh.inc"
    source "$_gcloud_sdk/completion.zsh.inc"
  }
  alias gcloud="unalias gcloud && _lazy_gcloud && gcloud"
fi
unset _gcloud_sdk

# completions
[[ -s "$HOME/.bun/_bun" ]] && source "$HOME/.bun/_bun"

# auto venv activation/deactivation on cd
auto_venv() {
  if [[ -d ".venv" ]]; then
    if [[ "$VIRTUAL_ENV" != "$PWD/.venv" ]] || ! (( $+functions[deactivate] )); then
      # deactivate current venv first (handles stale inherited venvs from zellij/tmux)
      if [[ -n "$VIRTUAL_ENV" ]]; then
        if (( $+functions[deactivate] )); then
          deactivate
        else
          path=("${(@)path:#$VIRTUAL_ENV/bin}")
          unset VIRTUAL_ENV VIRTUAL_ENV_PROMPT
        fi
      fi
      . ".venv/bin/activate"
    fi
  elif [[ -n "$VIRTUAL_ENV" ]]; then
    local venv_root="${VIRTUAL_ENV:h}"
    if [[ "$PWD" != "$venv_root" && "$PWD" != "$venv_root/"* ]]; then
      if (( $+functions[deactivate] )); then
        deactivate
      else
        path=("${(@)path:#$VIRTUAL_ENV/bin}")
        unset VIRTUAL_ENV VIRTUAL_ENV_PROMPT
      fi
    fi
  fi
}

chpwd_functions+=(auto_venv)
auto_venv

# machine-specific overrides
[[ -f ~/.zsh_local ]] && source ~/.zsh_local
