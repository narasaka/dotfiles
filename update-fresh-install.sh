#!/bin/bash

COMMIT_HASH=$(git rev-parse --short HEAD)

# Combine existing scripts into fresh-install.sh
echo "#!/bin/bash" > fresh-install.sh
echo "" >> fresh-install.sh
tail -n +2 mac_primer.sh >> fresh-install.sh
echo "" >> fresh-install.sh
echo "# Commit: $COMMIT_HASH" >> fresh-install.sh
echo "" >> fresh-install.sh
tail -n +2 install_mac_pkgs.sh >> fresh-install.sh
echo "" >> fresh-install.sh
tail -n +2 install.sh >> fresh-install.sh

# Make executable
chmod +x fresh-install.sh