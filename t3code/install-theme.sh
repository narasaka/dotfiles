#!/bin/sh
# Install the Mono B&W T3 Code theme (black/white, monospace, dithered planet
# background) and make it this machine's theme.
#
# Safe to run before T3 Code's first launch and safe to re-run. A running app
# picks the change up live.
#
# Files are copied, not stowed: T3 Code refuses to load theme files that are
# symlinks.

set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
theme_id=mono-bw
userdata="${T3CODE_HOME:-$HOME/.t3}/userdata"
themes_dir="$userdata/themes"

mkdir -p "$themes_dir"
# The background image sits beside the theme file under the same id.
cp "$here/themes/$theme_id.png" "$themes_dir/$theme_id.png"

if command -v t3 >/dev/null 2>&1; then
  # Publishes the JSON into the themes dir and selects it.
  t3 theme set --base-dir "${T3CODE_HOME:-$HOME/.t3}" "$here/themes/$theme_id.json"
else
  # T3 Code's CLI isn't on PATH yet (it's installed on the app's first
  # launch), so do what `t3 theme set` does: copy the file and set
  # defaultTheme in settings.json, keeping every other key.
  cp "$here/themes/$theme_id.json" "$themes_dir/$theme_id.json"
  python3 - "$userdata/settings.json" "$theme_id" <<'EOF'
import datetime, json, os, sys
path, theme_id = sys.argv[1], sys.argv[2]
settings = {}
if os.path.exists(path):
    with open(path) as f:
        settings = json.load(f)
settings["defaultTheme"] = theme_id
settings["defaultThemeSetAt"] = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"
tmp = path + ".tmp"
with open(tmp, "w") as f:
    json.dump(settings, f, indent=2)
    f.write("\n")
os.replace(tmp, path)
EOF
  echo "Installed \"$theme_id\" and set it as the T3 Code theme."
fi
