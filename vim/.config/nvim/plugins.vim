call plug#begin('~/.config/nvim/autoload/plugged')

	Plug 'sheerun/vim-polyglot'
	Plug 'jiangmiao/auto-pairs'
	Plug 'sainnhe/sonokai'
	Plug 'preservim/nerdtree'
	Plug 'preservim/nerdcommenter'
	Plug 'prettier/vim-prettier', {'do': 'yarn install'}
	Plug 'sbdchd/neoformat'

call plug#end()
