#!/usr/bin/env bash
# setup-ubuntu.sh — Reproduce narasaka's dev environment on Ubuntu 24.04
#
# Usage:
#   bash -c "$(curl -fsSL https://raw.githubusercontent.com/narasaka/dotfiles/master/setup-ubuntu.sh)"
#
# Or download and run:
#   curl -fsSL https://raw.githubusercontent.com/narasaka/dotfiles/master/setup-ubuntu.sh -o /tmp/setup-ubuntu.sh
#   bash /tmp/setup-ubuntu.sh
set -euo pipefail

# ─── Configuration ───────────────────────────────────────────────────────────

GO_VERSION="1.26.1"
SCRIPT_VERSION="0.0.3"
DOTFILES_REPO="https://github.com/narasaka/dotfiles.git"
DOTFILES_DIR="$HOME/.dotfiles"

# ─── Helpers ─────────────────────────────────────────────────────────────────

info()  { printf '\033[1;34m▸\033[0m %s\n' "$*"; }
ok()    { printf '\033[1;32m✓\033[0m %s\n' "$*"; }
warn()  { printf '\033[1;33m!\033[0m %s\n' "$*"; }
die()   { printf '\033[1;31m✗\033[0m %s\n' "$*" >&2; exit 1; }

has() { command -v "$1" &>/dev/null; }

ensure_dir() { [[ -d "$1" ]] || mkdir -p "$1"; }

# ─── Pre-flight ──────────────────────────────────────────────────────────────

main() {

[[ "$(uname -s)" == "Linux" ]] || die "This script is for Linux only."
if [[ -f /etc/os-release ]]; then
  # shellcheck source=/dev/null
  . /etc/os-release
  [[ "${ID:-}" == "ubuntu" ]] || warn "Targets Ubuntu 24.04. Detected: ${PRETTY_NAME:-$ID}"
fi

ARCH="$(dpkg --print-architecture 2>/dev/null || echo amd64)"
CODENAME="${UBUNTU_CODENAME:-${VERSION_CODENAME:-noble}}"

info "Running setup-ubuntu v${SCRIPT_VERSION}"
info "Setting up development environment on ${PRETTY_NAME:-Ubuntu}..."
echo ""

# ─── 1. Core prerequisites ──────────────────────────────────────────────────

info "Installing core prerequisites..."
sudo apt-get update -qq
sudo apt-get install -y ca-certificates curl wget gnupg gpg apt-transport-https \
  software-properties-common >/dev/null
ok "Core prerequisites"

# ─── 2. External APT repositories ───────────────────────────────────────────

info "Adding external APT repositories..."

# 1Password CLI
if [[ ! -f /etc/apt/sources.list.d/1password.list ]]; then
  info "  → 1Password CLI"
  curl -sS https://downloads.1password.com/linux/keys/1password.asc \
    | sudo gpg --dearmor --yes --output /usr/share/keyrings/1password-archive-keyring.gpg
  echo "deb [arch=${ARCH} signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/${ARCH} stable main" \
    | sudo tee /etc/apt/sources.list.d/1password.list >/dev/null
  sudo mkdir -p /etc/debsig/policies/AC2D62742012EA22/
  curl -sS https://downloads.1password.com/linux/debian/debsig/1password.pol \
    | sudo tee /etc/debsig/policies/AC2D62742012EA22/1password.pol >/dev/null
  sudo mkdir -p /usr/share/debsig/keyrings/AC2D62742012EA22
  curl -sS https://downloads.1password.com/linux/keys/1password.asc \
    | sudo gpg --dearmor --yes --output /usr/share/debsig/keyrings/AC2D62742012EA22/debsig.gpg
fi

# Docker CE
sudo rm -f /etc/apt/sources.list.d/docker.list
if [[ ! -f /etc/apt/sources.list.d/docker.sources ]]; then
  info "  → Docker CE"
  sudo install -m 0755 -d /etc/apt/keyrings
  sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
  sudo chmod a+r /etc/apt/keyrings/docker.asc
  cat <<-DOCKER | sudo tee /etc/apt/sources.list.d/docker.sources >/dev/null
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: ${CODENAME}
Components: stable
Architectures: ${ARCH}
Signed-By: /etc/apt/keyrings/docker.asc
DOCKER
fi

# Doppler CLI
if [[ ! -f /etc/apt/sources.list.d/doppler-cli.list ]]; then
  info "  → Doppler CLI"
  curl -sLf --retry 3 --tlsv1.2 --proto "=https" \
    'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' \
    | sudo gpg --dearmor --yes -o /usr/share/keyrings/doppler-archive-keyring.gpg
  echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" \
    | sudo tee /etc/apt/sources.list.d/doppler-cli.list >/dev/null
fi

# GitHub CLI
if [[ ! -f /etc/apt/sources.list.d/github-cli.list ]]; then
  info "  → GitHub CLI"
  sudo mkdir -p -m 755 /etc/apt/keyrings
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
    | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg >/dev/null
  sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
  echo "deb [arch=${ARCH} signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
    | sudo tee /etc/apt/sources.list.d/github-cli.list >/dev/null
fi

# Google Cloud SDK
if [[ ! -f /etc/apt/sources.list.d/google-cloud-sdk.list ]]; then
  info "  → Google Cloud SDK"
  curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
    | sudo gpg --dearmor --yes -o /usr/share/keyrings/cloud.google.gpg
  echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" \
    | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list >/dev/null
fi

# Helm
if [[ ! -f /etc/apt/sources.list.d/helm-stable-debian.list ]]; then
  info "  → Helm"
  curl -fsSL https://packages.buildkite.com/helm-linux/helm-debian/gpgkey \
    | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg >/dev/null
  echo "deb [signed-by=/usr/share/keyrings/helm.gpg] https://packages.buildkite.com/helm-linux/helm-debian/any/ any main" \
    | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list >/dev/null
fi

# Kubernetes (kubectl)
if [[ ! -f /etc/apt/sources.list.d/kubernetes.list ]]; then
  info "  → Kubernetes"
  sudo mkdir -p -m 755 /etc/apt/keyrings
  curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.35/deb/Release.key \
    | sudo gpg --dearmor --yes -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
  sudo chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg
  echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.35/deb/ /" \
    | sudo tee /etc/apt/sources.list.d/kubernetes.list >/dev/null
  sudo chmod 644 /etc/apt/sources.list.d/kubernetes.list
fi

# Tailscale
if [[ ! -f /etc/apt/sources.list.d/tailscale.list ]]; then
  info "  → Tailscale"
  curl -fsSL "https://pkgs.tailscale.com/stable/ubuntu/${CODENAME}.noarmor.gpg" \
    | sudo tee /usr/share/keyrings/tailscale-archive-keyring.gpg >/dev/null
  curl -fsSL "https://pkgs.tailscale.com/stable/ubuntu/${CODENAME}.tailscale-keyring.list" \
    | sudo tee /etc/apt/sources.list.d/tailscale.list >/dev/null
fi

ok "APT repositories configured"

# ─── 3. APT packages ────────────────────────────────────────────────────────

info "Installing APT packages (this may take a while)..."
sudo apt-get update -qq

APT_PACKAGES=(
  # build tools
  build-essential pkg-config zlib1g-dev libexpat1-dev make

  # version control
  git git-lfs

  # shell & terminal
  zsh zoxide stow locales

  # system utilities
  btop lsof openssh-server fuse3 unzip xclip sysstat

  # terminal fun
  neofetch toilet toilet-fonts

  # terminal image viewers
  chafa jp2a w3m w3m-img imagemagick

  # databases
  postgresql postgresql-client

  # python
  python3 python3-dev python3-pip python3-venv

  # container & cloud
  docker-ce docker-ce-cli containerd.io
  docker-buildx-plugin docker-compose-plugin
  1password-cli
  caddy
  doppler
  gh
  google-cloud-cli google-cloud-cli-gke-gcloud-auth-plugin
  helm
  kubectl
  tailscale
)

sudo apt-get install -y "${APT_PACKAGES[@]}"
ok "APT packages installed"

# ─── 4. Locale ───────────────────────────────────────────────────────────────

info "Configuring locale..."
sudo locale-gen en_US.UTF-8 >/dev/null 2>&1
sudo update-locale LANG=en_US.UTF-8 >/dev/null 2>&1
ok "Locale set to en_US.UTF-8"

# ─── 4b. Ghostty terminfo ───────────────────────────────────────────────────

if ! infocmp xterm-ghostty &>/dev/null; then
  warn "Ghostty terminfo not found. Run this from your Ghostty client:"
  warn "  infocmp -x xterm-ghostty | ssh $(hostname) -- tic -x -"
else
  ok "Ghostty terminfo (already installed)"
fi

# ─── 5. Go ───────────────────────────────────────────────────────────────────

if ! has go; then
  info "Installing Go ${GO_VERSION}..."
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf /tmp/go.tar.gz
  rm -f /tmp/go.tar.gz
fi
export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"
ok "Go $(go version | awk '{print $3}')"

# ─── 6. Rust ─────────────────────────────────────────────────────────────────

if ! has rustc; then
  info "Installing Rust..."
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --quiet
fi
# shellcheck source=/dev/null
[[ -f "$HOME/.cargo/env" ]] && . "$HOME/.cargo/env"
ok "Rust $(rustc --version | awk '{print $2}')"

# ─── 7. Neovim (AppImage) ───────────────────────────────────────────────────

if ! has nvim; then
  info "Installing Neovim..."
  ensure_dir "$HOME/.local/bin"
  curl -fsSL https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.appimage \
    -o "$HOME/.local/bin/nvim"
  chmod u+x "$HOME/.local/bin/nvim"
fi
export PATH="$HOME/.local/bin:$PATH"
ok "Neovim $(nvim --version 2>/dev/null | head -1 | awk '{print $2}')"

# ─── 8. fnm (Fast Node Manager) ──────────────────────────────────────────────

if ! has fnm; then
  info "Installing fnm..."
  curl -fsSL https://fnm.vercel.app/install | bash -s -- --skip-shell
fi
export PATH="$HOME/.local/share/fnm:$PATH"
eval "$(fnm env --use-on-cd --shell bash)"
ok "fnm $(fnm --version | awk '{print $2}')"

# Install Node.js LTS
info "Installing Node.js LTS via fnm..."
fnm install --lts
fnm default lts-latest
eval "$(fnm env --use-on-cd --shell bash)"
ok "Node.js $(node -v)"

# ─── 9. Bun ──────────────────────────────────────────────────────────────────

if ! has bun; then
  info "Installing Bun..."
  curl -fsSL https://bun.sh/install | bash
fi
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"
ok "Bun $(bun -v)"

# ─── 10. uv ──────────────────────────────────────────────────────────────────

if ! has uv; then
  info "Installing uv..."
  curl -LsSf https://astral.sh/uv/install.sh | sh
fi
ok "uv $(uv --version | awk '{print $2}')"

# ─── 11. mise ────────────────────────────────────────────────────────────────

if ! has mise; then
  info "Installing mise..."
  curl https://mise.run | sh
fi
ok "mise $(mise --version | awk '{print $1}')"

# ─── 12. k3d ─────────────────────────────────────────────────────────────────

if ! has k3d; then
  info "Installing k3d..."
  curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
fi
ok "k3d $(k3d version 2>/dev/null | head -1 | awk '{print $3}')"

# ─── 13. Tilt ────────────────────────────────────────────────────────────────

if ! has tilt; then
  info "Installing Tilt..."
  curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
fi
ok "Tilt $(tilt version 2>/dev/null | head -1)"

# ─── 14. op-ssh-load ────────────────────────────────────────────────────────

if ! has op-ssh-load; then
  info "Installing op-ssh-load..."
  ensure_dir "$HOME/.local/bin"
  curl -sSfL https://raw.githubusercontent.com/narasaka/op-ssh-load/main/op-ssh-load \
    -o "$HOME/.local/bin/op-ssh-load"
  chmod +x "$HOME/.local/bin/op-ssh-load"
  ensure_dir "$HOME/.config/op"
  chmod 700 "$HOME/.config/op"
fi
ok "op-ssh-load"

# ─── 15. Go tools ───────────────────────────────────────────────────────────

info "Installing Go tools..."
declare -A GO_TOOLS=(
  [lazygit]="github.com/jesseduffield/lazygit@latest"
  [lazydocker]="github.com/jesseduffield/lazydocker@latest"
  [ctlptl]="github.com/tilt-dev/ctlptl/cmd/ctlptl@latest"
  [goose]="github.com/pressly/goose/v3/cmd/goose@latest"
  [weasel]="github.com/narasaka/weasel@latest"
  [gopls]="golang.org/x/tools/gopls@latest"
  [cloud-sql-proxy]="github.com/GoogleCloudPlatform/cloud-sql-proxy/v2@latest"
)
for name in "${!GO_TOOLS[@]}"; do
  if ! has "$name"; then
    info "  → $name"
    go install "${GO_TOOLS[$name]}"
  else
    ok "  $name (already installed)"
  fi
done
ok "Go tools installed"

# ─── 16. Cargo tools ────────────────────────────────────────────────────────

info "Installing Cargo tools (this takes a while)..."
if ! has zellij; then
  info "  → zellij"
  cargo install --locked zellij
else
  ok "  zellij (already installed)"
fi
if ! has gws; then
  info "  → google-workspace-cli"
  cargo install --git https://github.com/googleworkspace/cli --locked
else
  ok "  gws (already installed)"
fi
if ! has rg; then
  info "  → ripgrep"
  cargo install ripgrep
else
  ok "  ripgrep (already installed)"
fi
ok "Cargo tools installed"

# ─── 17. Bun global packages ────────────────────────────────────────────────

info "Installing Bun global packages..."
bun install -g opencode-ai tree-sitter-cli
ok "Bun globals: opencode-ai, tree-sitter-cli"

# ─── 18. npm global packages ────────────────────────────────────────────────

info "Installing npm global packages..."
npm install -g typescript typescript-language-server
ok "npm globals: typescript, typescript-language-server"

# ─── 19. pnpm ────────────────────────────────────────────────────────────────

info "Enabling pnpm via corepack..."
npm install --global corepack@latest
corepack enable pnpm
ok "pnpm"

# ─── 20. Python tools (thefuck via uv) ──────────────────────────────────────

if ! has thefuck; then
  info "Installing thefuck via uv..."
  uv tool install thefuck --python 3.11
fi
ok "thefuck"

# ─── 21. Oh My Zsh ──────────────────────────────────────────────────────────

if [[ ! -d "$HOME/.oh-my-zsh" ]]; then
  info "Installing Oh My Zsh..."
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
  # Remove the generated .zshrc — our dotfiles provide it
  rm -f "$HOME/.zshrc"
fi
ok "Oh My Zsh"

# ─── 22. Dotfiles ───────────────────────────────────────────────────────────

if [[ ! -d "$DOTFILES_DIR" ]]; then
  info "Cloning dotfiles..."
  git clone "$DOTFILES_REPO" "$DOTFILES_DIR"
else
  info "Dotfiles already present, pulling latest..."
  git -C "$DOTFILES_DIR" pull --ff-only 2>/dev/null || true
fi

info "Applying dotfiles via stow..."
# Remove files that would conflict with stow symlinks
for f in .zshrc .gitconfig .vimrc; do
  [[ -f "$HOME/$f" && ! -L "$HOME/$f" ]] && rm -f "$HOME/$f"
done

cd "$DOTFILES_DIR"
LC_ALL=C stow */ --no-folding
ok "Dotfiles applied"

# ─── 23. Cargo env in .zshenv ───────────────────────────────────────────────

if ! grep -q 'cargo/env' "$HOME/.zshenv" 2>/dev/null; then
  echo '. "$HOME/.cargo/env"' >> "$HOME/.zshenv"
  ok "Cargo env added to .zshenv"
fi

# ─── 24. Systemd ssh-agent service ──────────────────────────────────────────

info "Setting up ssh-agent systemd user service..."
ensure_dir "$HOME/.config/systemd/user"
cat > "$HOME/.config/systemd/user/ssh-agent.service" <<'EOF'
[Unit]
Description=SSH Agent

[Service]
Type=simple
Environment=SSH_AUTH_SOCK=%t/ssh-agent.socket
ExecStart=/usr/bin/ssh-agent -D -a %t/ssh-agent.socket

[Install]
WantedBy=default.target
EOF

systemctl --user daemon-reload
systemctl --user enable ssh-agent.service
systemctl --user start ssh-agent.service 2>/dev/null || true
ok "ssh-agent service enabled"

# ─── 25. Docker group ───────────────────────────────────────────────────────

info "Adding user to docker group..."
sudo groupadd docker 2>/dev/null || true
sudo usermod -aG docker "$USER"
ok "User added to docker group"

# ─── 26. Default shell → zsh ────────────────────────────────────────────────

if [[ "$SHELL" != *"zsh"* ]]; then
  info "Setting zsh as default shell..."
  chsh -s "$(which zsh)"
  ok "Default shell → zsh"
else
  ok "zsh is already the default shell"
fi

# ─── Done ────────────────────────────────────────────────────────────────────

echo ""
echo "┌─────────────────────────────────────────────────────────────────┐"
echo "│  Setup complete!                                                │"
echo "│                                                                 │"
echo "│  Log out and back in, then:                                     │"
echo "│                                                                 │"
echo "│    tailscale up             # connect to tailnet                │"
echo "│    op signin                # sign in to 1Password              │"
echo "│    op-ssh-load              # load SSH keys                     │"
echo "│    gh auth login            # authenticate GitHub CLI           │"
echo "│    gcloud auth login        # authenticate Google Cloud         │"
echo "│    doppler login            # authenticate Doppler              │"
echo "│                                                                 │"
echo "│  Open nvim once to bootstrap plugins (lazy.nvim auto-installs) │"
echo "│                                                                 │"
echo "└─────────────────────────────────────────────────────────────────┘"
}

main "$@"
