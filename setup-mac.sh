#!/usr/bin/env bash
# setup-mac.sh — Reproduce narasaka's dev environment on macOS
#
# Usage:
#   bash -c "$(curl -fsSL https://raw.githubusercontent.com/narasaka/dotfiles/master/setup-mac.sh)" -- dev
#   bash -c "$(curl -fsSL https://raw.githubusercontent.com/narasaka/dotfiles/master/setup-mac.sh)" -- full
#
# Modes:
#   dev   CLI tools, languages, terminal, editor, core apps (Brewfile)
#   full  dev + every other app and tool (Brewfile + Brewfile.full)
#
# With no mode given it asks (defaults to dev). SETUP_MODE=dev|full also works.
#
# Only Homebrew + cloning the repo are fatal. Every other step runs isolated:
# a failure is recorded, the rest of the setup continues, and failures are
# listed at the end. Apps already installed outside Homebrew (MDM, App Store,
# manual download) are skipped instead of fought with.
#
# Note: written for the stock macOS /bin/bash (3.2) — no associative arrays.
set -euo pipefail

# ─── Configuration ───────────────────────────────────────────────────────────

SCRIPT_VERSION="0.0.2"
DOTFILES_REPO="https://github.com/narasaka/dotfiles.git"
DOTFILES_DIR="$HOME/.dotfiles"

# ─── Helpers ─────────────────────────────────────────────────────────────────

info()  { printf '\033[1;34m▸\033[0m %s\n' "$*"; }
ok()    { printf '\033[1;32m✓\033[0m %s\n' "$*"; }
warn()  { printf '\033[1;33m!\033[0m %s\n' "$*"; }
die()   { printf '\033[1;31m✗\033[0m %s\n' "$*" >&2; exit 1; }

has() { command -v "$1" &>/dev/null; }

ensure_dir() { [[ -d "$1" ]] || mkdir -p "$1"; }

usage() {
  echo "Usage: setup-mac.sh [dev|full]"
  echo "  dev   CLI tools, languages, terminal, editor, core apps (Brewfile)"
  echo "  full  dev + every other app and tool (Brewfile + Brewfile.full)"
}

FAILED=""
fail() { warn "$1 failed"; FAILED="${FAILED}  - $1"$'\n'; }

# Run a step function in a subshell with errexit on, so the first error
# aborts that step only. Env changes inside a step don't leak out — anything
# later steps need goes in load_env.
run_step() {
  local name="$1"; shift
  local rc=0
  set +e
  ( set -e; "$@" )
  rc=$?
  set -e
  (( rc == 0 )) || fail "$name"
}

# Put every tool this script installs on PATH (safe to call repeatedly)
load_env() {
  eval "$("$BREW" shellenv)"
  # shellcheck source=/dev/null
  [[ -f "$HOME/.cargo/env" ]] && . "$HOME/.cargo/env"
  export BUN_INSTALL="$HOME/.bun"
  export PATH="$HOME/.local/bin:$BUN_INSTALL/bin:$PATH:$HOME/go/bin"
  if has fnm; then eval "$(fnm env --shell bash)"; fi
  return 0
}

# Print Brewfile casks whose app/font/pkg is already on this machine but not
# managed by Homebrew. Installing over those fails (or is blocked by MDM).
preinstalled_casks() {
  local brewfile="$1" taps casks
  taps=$(sed -n 's/^tap "\([^"]*\)".*/\1/p' "$brewfile")
  casks=$(sed -n 's/^cask "\([^"]*\)".*/\1/p' "$brewfile")
  [[ -n "$casks" ]] || return 0
  for tap in $taps; do brew tap "$tap" >/dev/null 2>&1 || true; done
  # shellcheck disable=SC2086
  brew info --json=v2 --cask $casks 2>/dev/null | /usr/bin/python3 -c '
import json, os, re, subprocess, sys
receipts = subprocess.run(["pkgutil", "--pkgs"], capture_output=True, text=True).stdout.split()
home = os.path.expanduser("~")
for cask in json.loads(sys.stdin.read(), strict=False)["casks"]:
    if cask.get("installed"):
        continue
    paths, pkg_ids = [], []
    for artifact in cask["artifacts"]:
        for kind, base in (("app", "/Applications"), ("font", home + "/Library/Fonts")):
            for item in artifact.get(kind, []):
                if isinstance(item, str):
                    paths.append(os.path.join(base, os.path.basename(item)))
                    if kind == "app":
                        paths.append(os.path.join(home, "Applications", os.path.basename(item)))
        for stanza in artifact.get("uninstall", []):
            if isinstance(stanza, dict):
                ids = stanza.get("pkgutil") or []
                pkg_ids += [ids] if isinstance(ids, str) else ids
    present = any(os.path.exists(p) for p in paths) or any(
        re.fullmatch(pid, r) for pid in pkg_ids for r in receipts)
    if present:
        print(cask["full_token"])
' || warn "Couldn't check for pre-installed apps; installing everything" >&2
}

# ─── Steps ───────────────────────────────────────────────────────────────────

step_rust() {
  if ! has rustc; then
    info "Installing Rust..."
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --quiet --no-modify-path
    . "$HOME/.cargo/env"
  fi
  if [[ -f "$HOME/.cargo/env" ]] && ! grep -q 'cargo/env' "$HOME/.zshenv" 2>/dev/null; then
    echo '. "$HOME/.cargo/env"' >> "$HOME/.zshenv"
    ok "Cargo env added to .zshenv"
  fi
  ok "Rust $(rustc --version | awk '{print $2}')"
}

step_node() {
  info "Installing Node.js LTS via fnm..."
  fnm install --lts
  fnm default lts-latest
  eval "$(fnm env --shell bash)"
  ok "Node.js $(node -v)"
}

step_go_tools() {
  info "Installing Go tools..."
  local entry name pkg failed=0
  for entry in \
    weasel:github.com/narasaka/weasel@latest \
    gopls:golang.org/x/tools/gopls@latest
  do
    name="${entry%%:*}"
    pkg="${entry#*:}"
    if has "$name"; then
      ok "  $name (already installed)"
    else
      info "  → $name"
      go install "$pkg" || { warn "  $name failed"; failed=1; }
    fi
  done
  return "$failed"
}

step_cargo_tools() {
  if has gws; then
    ok "gws (already installed)"
    return 0
  fi
  info "Installing google-workspace-cli (cargo, takes a while)..."
  cargo install --git https://github.com/googleworkspace/cli --locked
  ok "gws"
}

step_bun_globals() {
  info "Installing Bun global packages..."
  bun install -g opencode-ai
  ok "Bun globals: opencode-ai"
}

step_npm_globals() {
  info "Installing npm global packages..."
  npm install -g typescript typescript-language-server corepack@latest
  corepack enable pnpm
  ok "npm globals: typescript, typescript-language-server, pnpm (corepack)"
}

step_claude() {
  if ! has claude; then
    info "Installing Claude Code..."
    curl -fsSL https://claude.ai/install.sh | bash
  fi
  ok "Claude Code"
}

step_gcloud() {
  has gcloud || return 0
  if ! has gke-gcloud-auth-plugin; then
    info "Installing gke-gcloud-auth-plugin..."
    gcloud components install gke-gcloud-auth-plugin --quiet >/dev/null 2>&1
  fi
  ok "gcloud + gke-gcloud-auth-plugin"
}

step_ohmyzsh() {
  if [[ ! -d "$HOME/.oh-my-zsh" ]]; then
    info "Installing Oh My Zsh..."
    local had_zshrc=0
    [[ -e "$HOME/.zshrc" || -L "$HOME/.zshrc" ]] && had_zshrc=1
    RUNZSH=no CHSH=no KEEP_ZSHRC=yes \
      sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
    # Remove the generated .zshrc — our dotfiles provide it
    if (( ! had_zshrc )); then rm -f "$HOME/.zshrc"; fi
  fi
  ok "Oh My Zsh"
}

step_dotfiles() {
  info "Applying dotfiles via stow..."
  cd "$DOTFILES_DIR"

  # Move anything that would block a stow symlink into a backup dir
  local backup_dir pkg src rel target
  backup_dir="$HOME/.dotfiles-backup/$(date +%Y%m%d-%H%M%S)"
  for pkg in */; do
    pkg="${pkg%/}"
    while IFS= read -r -d '' src; do
      rel="${src#"$pkg"/}"
      target="$HOME/$rel"
      if [[ -L "$target" ]]; then
        [[ "$(readlink "$target")" == *".dotfiles/"* ]] && continue
      elif [[ ! -e "$target" ]]; then
        continue
      fi
      ensure_dir "$(dirname "$backup_dir/$rel")"
      mv "$target" "$backup_dir/$rel"
      warn "  backed up ~/$rel → $backup_dir/$rel"
    done < <(find "$pkg" \( -type f -o -type l \) -print0)
  done

  LC_ALL=C stow --target="$HOME" --no-folding --restow */
  ok "Dotfiles applied"
}

step_tmux_plugins() {
  if [[ ! -d "$HOME/.tmux/plugins/tpm" ]]; then
    info "Installing tmux plugin manager..."
    git clone -q https://github.com/tmux-plugins/tpm "$HOME/.tmux/plugins/tpm"
  fi
  # install_plugins needs a running server; only spin one up if none exists
  if ! tmux has-session 2>/dev/null; then
    tmux new-session -d -s tpm-setup
    local rc=0
    "$HOME/.tmux/plugins/tpm/bin/install_plugins" >/dev/null || rc=$?
    tmux kill-server 2>/dev/null || true
    (( rc == 0 )) || return "$rc"
  fi
  ok "tmux plugins"
}

step_nvim_plugins() {
  info "Installing Neovim plugins (pinned to lazy-lock.json)..."
  nvim --headless "+Lazy! restore" +qa >/dev/null 2>&1
  ok "Neovim plugins"
}

step_macos_defaults() {
  info "Applying macOS defaults..."
  defaults write -g AppleInterfaceStyle -string Dark
  defaults write com.apple.dock autohide -bool true
  defaults write com.apple.dock autohide-delay -float 0
  defaults write com.apple.dock autohide-time-modifier -float 0
  defaults write com.apple.dock tilesize -int 26
  defaults write com.apple.dock orientation -string left
  killall Dock 2>/dev/null || true
  ok "macOS defaults (dark mode, dock: left, small, instant autohide)"
}

step_shell() {
  if [[ "$(dscl . -read "/Users/$USER" UserShell 2>/dev/null | awk '{print $2}')" == */zsh ]]; then
    ok "zsh is already the default shell"
    return 0
  fi
  info "Setting zsh as default shell..."
  chsh -s /bin/zsh
  ok "Default shell → zsh"
}

# ─── Main ────────────────────────────────────────────────────────────────────

main() {

[[ "$(uname -s)" == "Darwin" ]] || die "This script is for macOS only."

MODE="${SETUP_MODE:-}"
for arg in "$@"; do
  case "$arg" in
    dev|--dev)   MODE="dev" ;;
    full|--full) MODE="full" ;;
    -h|--help)   usage; exit 0 ;;
    *)           die "Unknown argument: $arg (expected dev or full)" ;;
  esac
done

if [[ -z "$MODE" ]]; then
  if [[ -r /dev/tty ]]; then
    printf 'Install type? [d]ev / [f]ull (default: dev): '
    read -r reply </dev/tty || reply=""
    case "$reply" in
      f|F|full) MODE="full" ;;
      *)        MODE="dev" ;;
    esac
  else
    MODE="dev"
  fi
fi
[[ "$MODE" == "dev" || "$MODE" == "full" ]] || die "Invalid mode: $MODE"

ARCH="$(uname -m)"

info "Running setup-mac v${SCRIPT_VERSION}"
info "Setting up ${MODE} environment on macOS $(sw_vers -productVersion) (${ARCH})..."
echo ""

# Ask for sudo once and keep it alive (Homebrew + some casks need it)
info "Requesting sudo (needed by Homebrew and some casks)..."
sudo -v
while true; do sudo -n true; sleep 50; kill -0 "$$" || exit; done 2>/dev/null &
SUDO_KEEPALIVE_PID=$!
trap 'kill "$SUDO_KEEPALIVE_PID" 2>/dev/null || true' EXIT

# ─── 1. Homebrew (also installs Xcode Command Line Tools) ───────────────────

if [[ -x /opt/homebrew/bin/brew ]]; then
  BREW=/opt/homebrew/bin/brew
elif [[ -x /usr/local/bin/brew ]]; then
  BREW=/usr/local/bin/brew
else
  info "Installing Homebrew..."
  NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  if [[ "$ARCH" == "arm64" ]]; then BREW=/opt/homebrew/bin/brew; else BREW=/usr/local/bin/brew; fi
fi
[[ -x "$BREW" ]] || die "Homebrew install failed"
eval "$("$BREW" shellenv)"
ok "Homebrew $(brew --version | head -1 | awk '{print $2}')"

# Login shells need brew on PATH (.zshrc doesn't set it up)
if ! grep -q 'brew shellenv' "$HOME/.zprofile" 2>/dev/null; then
  echo "eval \"\$(${BREW} shellenv)\"" >> "$HOME/.zprofile"
  ok "brew shellenv added to .zprofile"
fi

xcode-select -p &>/dev/null || die "Xcode Command Line Tools missing. Run: xcode-select --install"

if [[ "$MODE" == "full" && "$ARCH" == "arm64" ]] && ! /usr/bin/pgrep -q oahd; then
  info "Installing Rosetta 2 (some casks need it)..."
  sudo softwareupdate --install-rosetta --agree-to-license >/dev/null || fail "Rosetta 2"
fi

# ─── 2. Dotfiles repo ────────────────────────────────────────────────────────

if [[ ! -d "$DOTFILES_DIR" ]]; then
  info "Cloning dotfiles..."
  git clone "$DOTFILES_REPO" "$DOTFILES_DIR"
else
  info "Dotfiles already present, pulling latest..."
  git -C "$DOTFILES_DIR" pull --ff-only 2>/dev/null || warn "Couldn't fast-forward dotfiles, using local copy"
fi
ok "Dotfiles at $DOTFILES_DIR"

# ─── 3. Homebrew packages ────────────────────────────────────────────────────

BREWFILES="$DOTFILES_DIR/Brewfile"
[[ "$MODE" == "full" ]] && BREWFILES="$BREWFILES $DOTFILES_DIR/Brewfile.full"

for brewfile in $BREWFILES; do
  name="$(basename "$brewfile")"

  skip="$(preinstalled_casks "$brewfile" | tr '\n' ' ')"
  for cask in $skip; do
    warn "  skipping $cask (already installed outside Homebrew)"
  done

  info "Installing $name (this may take a while)..."
  if HOMEBREW_BUNDLE_CASK_SKIP="$skip" brew bundle --file="$brewfile"; then
    ok "$name installed"
  else
    fail "Some $name entries (see brew output above)"
  fi
done

load_env

# ─── 4. Languages & tools ───────────────────────────────────────────────────

run_step "Rust"                 step_rust;        load_env
run_step "Node.js"              step_node;        load_env
run_step "Go tools"             step_go_tools
run_step "Cargo tools (gws)"    step_cargo_tools
run_step "Bun globals"          step_bun_globals
run_step "npm globals + pnpm"   step_npm_globals
run_step "Claude Code"          step_claude
run_step "gcloud components"    step_gcloud

# ─── 5. Shell & dotfiles ────────────────────────────────────────────────────

run_step "Oh My Zsh"            step_ohmyzsh
run_step "Dotfiles (stow)"      step_dotfiles
run_step "tmux plugins"         step_tmux_plugins
run_step "Neovim plugins"       step_nvim_plugins

# ─── 6. macOS ───────────────────────────────────────────────────────────────

run_step "macOS defaults"       step_macos_defaults
run_step "Default shell"        step_shell

# ─── Done ────────────────────────────────────────────────────────────────────

echo ""
if [[ -n "$FAILED" ]]; then
  warn "Finished with failures (re-run the script to retry):"
  printf '%s' "$FAILED"
  echo ""
fi
echo "┌─────────────────────────────────────────────────────────────────┐"
printf '│  Setup complete! (%-4s)                                         │\n' "$MODE"
echo "│                                                                 │"
echo "│  Open a new terminal (Ghostty), then:                           │"
echo "│                                                                 │"
echo "│    open -a 1Password        # sign in, enable SSH agent + CLI   │"
echo "│    open -a Tailscale        # connect to tailnet                │"
echo "│    open -a Docker           # finish Docker Desktop setup       │"
echo "│    open -a Karabiner-Elements  # grant input permissions        │"
echo "│    gh auth login            # authenticate GitHub CLI           │"
echo "│    gcloud auth login        # authenticate Google Cloud         │"
echo "│    doppler login            # authenticate Doppler              │"
echo "│                                                                 │"
echo "│  Log out and back in for dark mode to fully apply.              │"
echo "│                                                                 │"
echo "└─────────────────────────────────────────────────────────────────┘"

[[ -z "$FAILED" ]]
}

main "$@"
