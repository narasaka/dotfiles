# Path to your oh-my-zsh installation.
export ZSH="/Users/narasaka/.oh-my-zsh"

# ZSH theme + plugins
ZSH_THEME="garyblessington"
plugins=(git web-search nvm)

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

# The next line updates PATH for the Google Cloud SDK.
if [ -f '/Users/narasaka/google-cloud-sdk/path.zsh.inc' ]; then . '/Users/narasaka/google-cloud-sdk/path.zsh.inc'; fi

# The next line enables shell command completion for gcloud.
if [ -f '/Users/narasaka/google-cloud-sdk/completion.zsh.inc' ]; then . '/Users/narasaka/google-cloud-sdk/completion.zsh.inc'; fi

export PATH=$PATH:/Applications/Postgres.app/Contents/Versions/15/bin
