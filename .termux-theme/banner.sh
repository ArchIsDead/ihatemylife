#!/data/data/com.termux/files/usr/bin/bash

DIR="$HOME/.termux-theme"

show_banner() {
  local file="$DIR/current/banner.txt"
  [ -f "$file" ] && cat "$file"
}

set_banner() {
  local src="$1"
  if [ -z "$src" ]; then
    echo "Usage: banner.sh set <file>"
    return 1
  fi
  if [ ! -f "$src" ]; then
    echo "[!] File not found: $src"
    return 1
  fi
  cp "$src" "$DIR/current/banner.txt"
  echo "[+] Banner updated"
}

reset_banner() {
  local theme=$(basename $(readlink "$DIR/current"))
  local backup="$DIR/themes/$theme/banner.default"
  if [ -f "$backup" ]; then
    cp "$backup" "$DIR/current/banner.txt"
    echo "[+] Banner reset to default"
  else
    echo "[!] No default banner found"
  fi
}

case "$1" in
  show) show_banner ;;
  set) set_banner "$2" ;;
  reset) reset_banner ;;
  edit) nano "$DIR/current/banner.txt" ;;
  *) echo "Usage: banner.sh [show|set <file>|reset|edit]" ;;
esac
