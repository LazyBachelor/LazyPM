#!/bin/bash
# This script runs as root when the container starts
echo "**** Setting up custom desktop shortcuts ****"

# Ensure the Desktop folder exists for the default 'abc' user
mkdir -p /config/Desktop

cat <<EOF > /config/Desktop/MyGoApp.desktop
[Desktop Entry]
Version=1.0
Type=Application
Name=PM Survey
Comment=Launch my custom Go binary
Exec=/usr/local/bin/survey start
Icon=utilities-terminal
Terminal=true
Categories=Utility;
EOF

chmod +x /config/Desktop/MyGoApp.desktop
chown abc:abc /config/Desktop/MyGoApp.desktop

echo "**** Desktop setup complete ****"
