syntax enable
set hidden
set pumheight=10
set iskeyword+=-
set splitbelow
set splitright
set t_Co=256
set sw=2
set ts=2
set smarttab
set smartindent
set autoindent
set nu rnu
set cursorline
set updatetime=300
set timeoutlen=500
set formatoptions-=cro
set clipboard=unnamedplus

function! GitBranch()
  return system("git rev-parse --abbrev-ref HEAD 2>/dev/null | tr -d '\n'")
endfunction

function! StatuslineGit()
  let l:branchname = GitBranch()
  return strlen(l:branchname) > 0?'  '.l:branchname.' ':''
endfunction

set statusline=
set statusline+=%#PmenuSel#
set statusline+=%{StatuslineGit()}
set statusline+=%#LineNr#
set statusline+=\ %f
set statusline+=%m

augroup numbertoggle
	autocmd!
	autocmd BufEnter,FocusGained,InsertLeave,WinEnter * if &nu && mode() != "i" | set rnu   | endif
	autocmd BufLeave,FocusLost,InsertEnter,WinLeave   * if &nu                  | set nornu | endif
augroup END



