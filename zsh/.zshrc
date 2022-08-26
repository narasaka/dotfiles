# Path to your oh-my-zsh installation.
export ZSH="/home/narasaka/.oh-my-zsh"

ZSH_THEME="dpoggi"

plugins=(git web-search nvm)

# sources
source $ZSH/oh-my-zsh.sh

# User configuration

export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
export EDITOR='vim'
export TIMEFMT=$'\nreal\t%E\nuser\t%U\nsys\t%S'
# export MANPATH="/usr/local/man:$MANPATH"

 #Preferred editor for local and remote sessions
 if [[ -n $SSH_CONNECTION ]]; then
   export EDITOR='vim'
 else
   export EDITOR='vim'
 fi

# vi bindings
#set -o vi

# Compilation flags
# export ARCHFLAGS="-arch x86_64"

# aliases
alias gg='g++ -std=c++17 -Wshadow -Wall -O2 -Wno-unused-result'
alias gf='g++ -std=c++17 -Wshadow -Wall -g -fsanitize=address -fsanitize=undefined -D_GLIBCXX_DEGUG'
alias cptemp='cp ~/prog/templates/temp.cpp solve.cpp'
alias pytemp='cp ~/prog/templates/temp.py solve.py'
alias tmux='TERM=screen-256color-bce tmux'
alias vim='lvim'
alias python='python3'
alias clip='clip.exe'
alias pypy='pypy3'

# refer(s)

[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
