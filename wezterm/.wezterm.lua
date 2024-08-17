local wezterm = require 'wezterm'

local config = wezterm.config_builder()

config.color_scheme = 'Azu (Gogh)'
config.font = wezterm.font 'FiraCode Nerd Font'
config.font_size = 14.0

return config

