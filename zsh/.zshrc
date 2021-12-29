# Path to your oh-my-zsh installation.
export ZSH="/home/narasaka/.oh-my-zsh"

ZSH_THEME="dpoggi"

# Which plugins would you like to load?
# Standard plugins can be found in $ZSH/plugins/
# Custom plugins may be added to $ZSH_CUSTOM/plugins/
# Example format: plugins=(rails git textmate ruby lighthouse)
# Add wisely, as too many plugins slow down shell startup.
plugins=(git web-search nvm)

# sources
source $ZSH/oh-my-zsh.sh

# User configuration

export DISPLAY=127.0.0.1:0.0
export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
export EDITOR='vim'
export TIMEFMT=$'\nreal\t%E\nuser\t%U\nsys\t%S'
# export MANPATH="/usr/local/man:$MANPATH"

 #Preferred editor for local and remote sessions
 if [[ -n $SSH_CONNECTION ]]; then
   export EDITOR='vim'
 else
   export EDITOR='mvim'
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
#alias vim='nvim'
alias league='sudo sysctl -w abi.vsyscall32=0'
alias python='python3'
alias wh='/mnt/c/Users/given'
alias clip='clip.exe'
alias pypy='pypy3'

# stop cursor blink (print DECSCUSR)
echo -e "\e[2 q"
