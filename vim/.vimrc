set nu rnu
set sw=2
set ts=2
set autoindent
set nowrap
set smartcase
set noswapfile
set showcmd
set incsearch
set clipboard=unnamedplus
set splitright
set splitbelow

augroup numbertoggle
  autocmd!
  autocmd BufEnter,FocusGained,InsertLeave,WinEnter * if &nu && mode() != "i" | set rnu   | endif
  autocmd BufLeave,FocusLost,InsertEnter,WinLeave   * if &nu                  | set nornu | endif
augroup END

packloadall

"Vimplug
call plug#begin('~/.vim/plugged')


Plug 'sainnhe/sonokai'
Plug 'preservim/nerdtree'
Plug 'preservim/nerdcommenter'
"Plug 'christoomey/vim-tmux-navigator'
Plug 'tpope/vim-fugitive'
Plug 'prettier/vim-prettier', {'do': 'yarn install'}
Plug 'sbdchd/neoformat'

call plug#end()

if has('termguicolors')
	set termguicolors
endif

"Sonokai
let g:sonokai_style = 'shusia'
let g:sonokai_enable_italic = 1
let g:sonokai_disable_italic_comment = 1

colorscheme sonokai
filetype plugin on

"paths
let g:python3_host_prog = expand("/usr/bin/python")
let g:node_host_prog = expand("/home/narasaka/nvm/versions/node/v15.14.0/bin/node")

"NERDTree
nnoremap <C-\> :NERDTreeToggle<CR>

"Key bindings
nnoremap <C-h> <C-w>h
nnoremap <C-j> <C-w>j
nnoremap <C-k> <C-w>k
nnoremap <C-l> <C-w>l

"Statusline
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
