#!/bin/bash
# This script runs as root when the container starts
echo "**** Setting up custom desktop shortcuts ****"

# Ensure the Desktop folder exists for the default 'abc' user
mkdir -p /config/Desktop

cat <<EOF > /config/Desktop/LazyPM.desktop
[Desktop Entry]
Version=1.0
Type=Application
Name=PM Survey
Comment=Launch the Prompt Manageragment Survey
Exec=/usr/local/bin/pm survey start
Icon=utilities-terminal
Terminal=true
Categories=Utility;
EOF

chmod +x /config/Desktop/LazyPM.desktop
chown abc:abc /config/Desktop/LazyPM.desktop

pm # initialize the service to ensure the database is set up

echo "**** Setting default browser to Chromium... ****"

xdg-settings set default-web-browser chromium.desktop

echo "**** Setting default handlers for HTTP and HTTPS... ****"
xdg-mime default chromium.desktop x-scheme-handler/http
xdg-mime default chromium.desktop x-scheme-handler/https

echo "**** Setting up bash completion... ****"
source /etc/bash/bash_completion.sh

echo "Setting Alacritty as default terminal..."
mkdir -p /config/.config/xfce4

cat > /config/.config/xfce4/helpers.rc <<EOF
TerminalEmulator=alacritty.desktop
TerminalEmulatorDismissed=true
EOF

chown -R abc:abc /config/.config

echo "**** Desktop setup complete ****"