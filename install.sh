#!/bin/sh

set -eu

dotfiles_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$dotfiles_dir"

LC_ALL=C stow --target="$HOME" --no-folding */
