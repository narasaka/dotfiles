# Path to your oh-my-zsh installation.
export ZSH="/Users/narasaka/.oh-my-zsh"

# ZSH theme + plugins
ZSH_THEME="dpoggi"
plugins=(git web-search nvm vi-mode)

# vi-mode settings
VI_MODE_RESET_PROMPT_ON_MODE_CHANGE=true
VI_MODE_SET_CURSOR=true
VI_MODE_CURSOR_NORMAL=2
VI_MODE_CURSOR_VISUAL=6
VI_MODE_CURSOR_INSERT=6
VI_MODE_CURSOR_OPPEND=0

# sources
source $ZSH/oh-my-zsh.sh

# User configuration
export EDITOR='vi'
export TIMEFMT=$'\nreal\t%E\nuser\t%U\nsys\t%S'
# export MANPATH="/usr/local/man:$MANPATH"

# Preferred editor for local and remote sessions
 if [[ -n $SSH_CONNECTION ]]; then
   export EDITOR='vi'
 else
   export EDITOR='vi'
 fi

# Compilation flags
# export ARCHFLAGS="-arch x86_64"

# aliases
alias gg='g++ -std=c++17 -Wshadow -Wall -O2 -Wno-unused-result'
alias gf='g++ -std=c++17 -Wshadow -Wall -g -fsanitize=address -fsanitize=undefined -D_GLIBCXX_DEGUG'
alias cptemp='cp ~/prog/templates/temp.cpp solve.cpp'
alias pytemp='cp ~/prog/templates/temp.py solve.py'
alias vim='lvim'
alias python='python3'
alias pypy='pypy3'

# refer(s)
[ -s "/opt/homebrew/opt/nvm/nvm.sh" ] && \. "/opt/homebrew/opt/nvm/nvm.sh"  # This loads nvm
[ -s "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm" ] && \. "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm"  # This loads nvm bash_completion

export PATH=/Users/narasaka/.local/bin:$PATH
export PATH=$PATH:/Applications/Postgres.app/Contents/Versions/15/bin

# The next line updates PATH for the Google Cloud SDK.
if [ -f '/Users/narasaka/Downloads/google-cloud-sdk/path.zsh.inc' ]; then . '/Users/narasaka/Downloads/google-cloud-sdk/path.zsh.inc'; fi

# The next line enables shell command completion for gcloud.
if [ -f '/Users/narasaka/Downloads/google-cloud-sdk/completion.zsh.inc' ]; then . '/Users/narasaka/Downloads/google-cloud-sdk/completion.zsh.inc'; fi

# bun completions
[ -s "/Users/narasaka/.bun/_bun" ] && source "/Users/narasaka/.bun/_bun"

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# go
export PATH=/Users/narasaka/go/bin:$PATH

# pnpm
export PNPM_HOME="/Users/narasaka/Library/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac
