#!/data/data/com.termux/files/usr/bin/bash

DIR="$HOME/.termux-theme"
THEME="$1"

if [ -z "$THEME" ]; then
  echo "Usage: apply.sh <theme>"
  ls "$DIR/themes"
  exit 1
fi

SRC="$DIR/themes/$THEME"
if [ ! -d "$SRC" ]; then
  echo "[!] Theme not found: $THEME"
  exit 1
fi

mkdir -p "$HOME/.termux"

[ -f "$SRC/colors.properties" ] && cp "$SRC/colors.properties" "$HOME/.termux/colors.properties"
[ -f "$SRC/termux.properties" ] && cp "$SRC/termux.properties" "$HOME/.termux/termux.properties"
[ -f "$SRC/zshrc" ] && cp "$SRC/zshrc" "$HOME/.zshrc"

ln -sfn "$SRC" "$DIR/current"

termux-reload-settings

echo "[+] Theme applied: $THEME"
