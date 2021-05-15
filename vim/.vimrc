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

call plug#begin('~/.vim/plugged')


Plug 'sainnhe/sonokai'
Plug 'preservim/nerdtree'
Plug 'preservim/nerdcommenter'

call plug#end()

if has('termguicolors')
	set termguicolors
endif

let g:sonokai_style = 'shusia'
let g:sonokai_enable_italic = 1
let g:sonokai_disable_italic_comment = 1

colorscheme sonokai
filetype plugin on

"NERDTree
nnoremap <C-\> :NERDTreeToggle<CR>
nnoremap <C-h> <C-w>h
nnoremap <C-j> <C-w>j
nnoremap <C-k> <C-w>k
nnoremap <C-l> <C-w>l

