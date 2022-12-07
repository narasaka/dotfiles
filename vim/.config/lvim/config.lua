lvim.leader = "space"

-- override
vim.list_extend(lvim.lsp.automatic_configuration.skipped_servers, { "tsserver", "denols" })
local lspconfig = require "lspconfig"
require("lvim.lsp.manager").setup("tsserver", { root_dir = lspconfig.util.root_pattern("package.json") })
require("lvim.lsp.manager").setup("denols", { root_dir = lspconfig.util.root_pattern("deno.json") })

-- general
lvim.log.level = "warn"
lvim.format_on_save = false
lvim.colorscheme = "tokyonight"

-- builtins
lvim.builtin.alpha.active = true
lvim.builtin.alpha.mode = "dashboard"
lvim.builtin.notify.active = true
lvim.builtin.terminal.active = true
lvim.builtin.nvimtree.setup.view.side = "left"
lvim.builtin.nvimtree.setup.renderer.icons.show.git = false
lvim.builtin.treesitter.ignore_install = { "haskell" }
lvim.builtin.treesitter.highlight.enabled = true

-- keybinds
lvim.keys.insert_mode["<A-j>"] = false
lvim.keys.insert_mode["<A-k>"] = false
lvim.keys.normal_mode["<A-j>"] = false
lvim.keys.normal_mode["<A-k>"] = false
lvim.keys.normal_mode["<Leader>t"] = ":TodoTelescope<cr>"
lvim.keys.normal_mode["t"] = ":TodoTrouble<cr>"
lvim.keys.normal_mode["T"] = ":TroubleToggle<cr>"
lvim.keys.normal_mode["<C-s>"] = ":Prettier<cr>"
lvim.keys.visual_block_mode["<A-j>"] = false
lvim.keys.visual_block_mode["<A-k>"] = false
lvim.keys.visual_block_mode["J"] = false
lvim.keys.visual_block_mode["K"] = false

-- vim
vim.opt.clipboard = ""
vim.opt.ignorecase = true
vim.opt.termguicolors = true
vim.opt.title = true
vim.opt.titlestring = "%<%F%=%l/%L - nvim"
vim.opt.undodir = vim.fn.stdpath "cache" .. "/undo"
vim.opt.undofile = true
vim.opt.expandtab = true
vim.opt.shiftwidth = 2
vim.opt.tabstop = 2
vim.opt.cursorline = true
vim.opt.relativenumber = true
vim.opt.wrap = false
vim.opt.spell = false
vim.opt.spelllang = "en"

-- if you don't want all the parsers change this to a table of the ones you want
-- lvim.builtin.treesitter.ensure_installed = {
--   "bash",
--   "c",
--   "javascript",
--   "json",
--   "lua",
--   "python",
--   "typescript",
--   "tsx",
--   "css",
--   "rust",
--   "java",
--   "yaml",
--   "astro"
-- }

-- lsp
local formatters = require "lvim.lsp.null-ls.formatters"
formatters.setup {
  { command = "black" },
  {
    command = "prettier",
    args = { "--trailing-comma=es5", "--single-quote" },
  }
}

-- Additional Plugins
lvim.plugins = {
  { "folke/tokyonight.nvim" },
  { "folke/trouble.nvim" },
  { "folke/todo-comments.nvim",
    event = "BufRead",
    config = function() require("todo-comments").setup() end },
  { "tpope/vim-surround" },
  { "tpope/vim-repeat" },
  { "wuelnerdotexe/vim-astro" },
  { "christoomey/vim-tmux-navigator" }
}
