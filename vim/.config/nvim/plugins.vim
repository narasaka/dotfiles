call plug#begin('~/.config/nvim/autoload/plugged')

	Plug 'jiangmiao/auto-pairs'
	Plug 'preservim/nerdtree'
	Plug 'preservim/nerdcommenter'
	Plug 'prettier/vim-prettier', {'do': 'yarn install'}
	Plug 'sainnhe/sonokai'
	Plug 'sbdchd/neoformat'
	Plug 'sheerun/vim-polyglot'
	Plug 'tpope/vim-surround'
	Plug 'tpope/vim-repeat'

call plug#end()
