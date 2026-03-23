#!/bin/bash
# This script runs as root when the container starts
echo "**** Setting up custom desktop shortcuts ****"

mkdir -p /config/Desktop
mkdir -p /config/.config/xfce4/xfconf/xfce-perchannel-xml

cat <<EOF > /config/Desktop/LazyPM.desktop
[Desktop Entry]
Version=1.0
Type=Application
Name=PM Survey
Comment=Launch the Project Management Survey
Exec=/usr/local/bin/pm start
Icon=utilities-terminal
Terminal=true
Categories=Utility;
EOF

cat <<EOF > /config/.config/xfce4/xfconf/xfce-perchannel-xml/xfce4-terminal.xml
<?xml version="1.1" encoding="UTF-8"?>

<channel name="xfce4-terminal" version="1.0">
  <property name="background-mode" type="string" value="TERMINAL_BACKGROUND_SOLID"/>
  <property name="cell-height-scale" type="double" value="1.2"/>
  <property name="cell-width-scale" type="double" value="1"/>
  <property name="color-background" type="string" value="#000000"/>
  <property name="color-background-vary" type="bool" value="false"/>
  <property name="color-bold" type="string" value=""/>
  <property name="color-bold-is-bright" type="bool" value="true"/>
  <property name="color-bold-use-default" type="bool" value="true"/>
  <property name="color-cursor" type="string" value=""/>
  <property name="color-cursor-foreground" type="string" value=""/>
  <property name="color-cursor-use-default" type="bool" value="true"/>
  <property name="color-foreground" type="string" value="#ffffff"/>
  <property name="color-palette" type="string" value="#000000;#aa0000;#00aa00;#aa5500;#0000aa;#aa00aa;#00aaaa;#aaaaaa;#555555;#ff5555;#55ff55;#ffff55;#5555ff;#ff55ff;#55ffff;#ffffff"/>
  <property name="color-selection" type="string" value=""/>
  <property name="color-selection-background" type="string" value=""/>
  <property name="color-selection-use-default" type="bool" value="true"/>
  <property name="color-use-theme" type="bool" value="true"/>
  <property name="font-allow-bold" type="bool" value="true"/>
  <property name="font-name" type="string" value="Liberation Mono Bold 13"/>
  <property name="font-use-system" type="bool" value="false"/>
  <property name="tab-activity-color" type="string" value="#aa0000"/>
  <property name="misc-copy-on-select" type="bool" value="true"/>
  <property name="misc-show-unsafe-paste-dialog" type="bool" value="false"/>
</channel>
EOF

cat <<EOF > /config/.config/xfce4/xfconf/xfce-perchannel-xml/xsettings.xml
<?xml version="1.1" encoding="UTF-8"?>

<channel name="xsettings" version="1.0">
  <property name="Net" type="empty">
    <property name="ThemeName" type="string" value="adw-gtk3-dark"/>
    <property name="IconThemeName" type="empty"/>
    <property name="DoubleClickTime" type="empty"/>
    <property name="DoubleClickDistance" type="empty"/>
    <property name="DndDragThreshold" type="empty"/>
    <property name="CursorBlink" type="empty"/>
    <property name="CursorBlinkTime" type="empty"/>
    <property name="SoundThemeName" type="empty"/>
    <property name="EnableEventSounds" type="empty"/>
    <property name="EnableInputFeedbackSounds" type="empty"/>
  </property>
  <property name="Xft" type="empty">
    <property name="DPI" type="int" value="120"/>
    <property name="Antialias" type="empty"/>
    <property name="Hinting" type="empty"/>
    <property name="HintStyle" type="empty"/>
    <property name="RGBA" type="empty"/>
  </property>
  <property name="Gtk" type="empty">
    <property name="CanChangeAccels" type="empty"/>
    <property name="ColorPalette" type="empty"/>
    <property name="FontName" type="empty"/>
    <property name="MonospaceFontName" type="empty"/>
    <property name="IconSizes" type="empty"/>
    <property name="KeyThemeName" type="empty"/>
    <property name="MenuImages" type="empty"/>
    <property name="ButtonImages" type="empty"/>
    <property name="MenuBarAccel" type="empty"/>
    <property name="CursorThemeName" type="empty"/>
    <property name="CursorThemeSize" type="int" value="40"/>
    <property name="DecorationLayout" type="string" value="icon,menu:minimize,maximize,close"/>
    <property name="DialogsUseHeader" type="empty"/>
    <property name="TitlebarMiddleClick" type="empty"/>
  </property>
  <property name="Gdk" type="empty">
    <property name="WindowScalingFactor" type="empty"/>
  </property>
  <property name="Xfce" type="empty">
    <property name="SyncThemes" type="bool" value="true"/>
  </property>
</channel>
EOF

chmod +x /config/Desktop/LazyPM.desktop
chown -R abc:abc /config/Desktop /config/.config

git config --global user.name "LazyPM"
git config --global user.email "lazy@pm.com"

echo "**** Desktop setup complete ****"