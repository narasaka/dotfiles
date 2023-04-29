lvim.leader = "space"

-- overrides
vim.list_extend(lvim.lsp.automatic_configuration.skipped_servers, { "tsserver", "denols", "clangd" })
local lspconfig = require "lspconfig"
local capabilities = vim.lsp.protocol.make_client_capabilities()
capabilities.offsetEncoding = 'utf-8'
require("lvim.lsp.manager").setup("tsserver", { root_dir = lspconfig.util.root_pattern("package.json") })
require("lvim.lsp.manager").setup("denols", { root_dir = lspconfig.util.root_pattern("deno.json") })
require("lvim.lsp.manager").setup("clangd", { capabilities = capabilities })
require("lvim.lsp.manager").setup("tailwindcss", { 
  settings = { 
    tailwindCSS = { 
      experimental = {
        classRegex = {
          {"cva\\(([^)]*)\\)",
             "[\"'`]([^\"'`]*).*?[\"'`]"},
        }
      }
    }
  }
})
require("luasnip").filetype_extend("typescriptreact", { "html" })
require("luasnip").filetype_extend("typescript", { "javascript" })

-- general
lvim.log.level = "warn"
lvim.format_on_save = false
lvim.colorscheme = "tokyonight-night"

-- builtins
lvim.builtin.alpha.active = true
lvim.builtin.alpha.mode = "dashboard"
lvim.builtin.terminal.active = true
lvim.builtin.nvimtree.setup.view.side = "left"
lvim.builtin.nvimtree.setup.renderer.icons.show.git = false
lvim.builtin.nvimtree.setup.view.adaptive_size = false
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
    filetypes = { "typescript", "typescriptreact", "javascript", "astro", "css", "scss" },
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
  { "christoomey/vim-tmux-navigator" },
  {
    "zbirenbaum/copilot.lua",
    event = "VimEnter",
    config = function()
      vim.defer_fn(function()
        require("copilot").setup(
          {
            suggestion = {
              enabled = true,
              auto_trigger = true,
              debounce = 75,
              keymap = {
                accept = "<C-f>",
                accept_word = false,
                accept_line = false,
                next = "<M-]>",
                prev = "<M-[>",
                dismiss = "<C-]>",
              },
            },
          }
        )
      end, 100)
    end,
  },
  { "zbirenbaum/copilot-cmp",
    after = { "copilot.lua", "nvim-cmp" },
  },
  { "lukas-reineke/indent-blankline.nvim" },
  { 'olivercederborg/poimandres.nvim' }
}

table.insert(lvim.builtin.cmp.sources, 1, { name = "copilot" })
