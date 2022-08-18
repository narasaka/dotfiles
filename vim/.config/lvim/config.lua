lvim.leader = "space"

-- general
lvim.log.level = "warn"
lvim.format_on_save = true
lvim.colorscheme = "tokyonight"

-- builtins
lvim.builtin.alpha.active = true
lvim.builtin.alpha.mode = "dashboard"
lvim.builtin.notify.active = true
lvim.builtin.terminal.active = true
lvim.builtin.nvimtree.setup.view.side = "left"
lvim.builtin.nvimtree.setup.renderer.icons.show.git = false

-- vim
vim.opt.clipboard = "unnamedplus"
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
vim.opt.relativenumber = false
vim.opt.wrap = false
vim.opt.spell = false
vim.opt.spelllang = "en"

-- if you don't want all the parsers change this to a table of the ones you want
lvim.builtin.treesitter.ensure_installed = {
  "bash",
  "c",
  "javascript",
  "json",
  "lua",
  "python",
  "typescript",
  "tsx",
  "css",
  "rust",
  "java",
  "yaml",
}

lvim.lsp.automatic_servers_installation = true

lvim.builtin.treesitter.ignore_install = { "haskell" }
lvim.builtin.treesitter.highlight.enabled = true

-- Additional Plugins
lvim.plugins = {
    {"folke/tokyonight.nvim"},
    {"tpope/vim-surround"},
    {"tpope/vim-repeat"}
}
