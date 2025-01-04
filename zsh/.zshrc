export ZSH="/Users/narasaka/.oh-my-zsh"

# zsh theme + plugins
ZSH_THEME="dpoggi"
plugins=(git web-search vi-mode)

# lazy load nvm
export NVM_DIR="$HOME/.nvm"
_lazy_nvm() {
  [ -s "/opt/homebrew/opt/nvm/nvm.sh" ] && \. "/opt/homebrew/opt/nvm/nvm.sh"
  [ -s "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm" ] && \. "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm"
}
alias nvm="unalias nvm && _lazy_nvm && nvm"

# vi-mode
VI_MODE_RESET_PROMPT_ON_MODE_CHANGE=true
VI_MODE_SET_CURSOR=true
VI_MODE_CURSOR_NORMAL=2
VI_MODE_CURSOR_VISUAL=6
VI_MODE_CURSOR_INSERT=6
VI_MODE_CURSOR_OPPEND=0

# omz
source $ZSH/oh-my-zsh.sh

# misc
export EDITOR='vi'
export TIMEFMT=$'\nreal\t%E\nuser\t%U\nsys\t%S'

# aliases
alias gg='g++ -std=c++17 -Wshadow -Wall -O2 -Wno-unused-result'
alias gf='g++ -std=c++17 -Wshadow -Wall -g -fsanitize=address -fsanitize=undefined -D_GLIBCXX_DEGUG'
alias cptemp='cp ~/prog/templates/temp.cpp solve.cpp'
alias pytemp='cp ~/prog/templates/temp.py solve.py'
alias vim='nvim'
alias python='python3'
alias pypy='pypy3'

# paths
export PATH=/Users/narasaka/.local/bin:$PATH
export PATH=$PATH:/Applications/Postgres.app/Contents/Versions/15/bin:/Users/narasaka/go/bin
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# lazy load Google Cloud SDK
if [ -f '/Users/narasaka/Downloads/google-cloud-sdk/path.zsh.inc' ]; then
  function _lazy_gcloud() {
    . '/Users/narasaka/Downloads/google-cloud-sdk/path.zsh.inc'
    . '/Users/narasaka/Downloads/google-cloud-sdk/completion.zsh.inc'
  }
  alias gcloud="unalias gcloud && _lazy_gcloud && gcloud"
fi

# bun completions
[ -s "/Users/narasaka/.bun/_bun" ] && source "/Users/narasaka/.bun/_bun"

# pnpm setup
export PNPM_HOME="/Users/narasaka/Library/pnpm"
export PATH="$PNPM_HOME:$PATH"

# f
eval "$(thefuck --alias)"

