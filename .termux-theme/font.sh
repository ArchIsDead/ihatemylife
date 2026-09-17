#!/data/data/com.termux/files/usr/bin/bash

FONTDIR="$HOME/.termux"

apply_font() {
  local src="$1"
  if [ -z "$src" ]; then
    echo "Usage: font.sh apply <url|path>"
    return 1
  fi

  mkdir -p "$FONTDIR"
  local tmp="/tmp/font_download"

  if echo "$src" | grep -qE '^https?://'; then
    echo "[*] Downloading font..."
    curl -L -o "$tmp" "$src"
    if [ ! -f "$tmp" ]; then
      echo "[!] Download failed"
      return 1
    fi
    src="$tmp"
  fi

  if [ ! -f "$src" ]; then
    echo "[!] File not found: $src"
    return 1
  fi

  case "$src" in
    *.zip)
      echo "[*] Extracting zip..."
      mkdir -p /tmp/font_extract
      rm -rf /tmp/font_extract/*
      unzip -o "$src" -d /tmp/font_extract >/dev/null
      local ttf=$(find /tmp/font_extract -iname "*.ttf" | head -n 1)
      if [ -z "$ttf" ]; then
        echo "[!] No .ttf found in zip"
        return 1
      fi
      src="$ttf"
      ;;
  esac

  cp "$src" "$FONTDIR/font.ttf"
  termux-reload-settings
  echo "[+] Font applied"
}

remove_font() {
  rm -f "$FONTDIR/font.ttf"
  termux-reload-settings
  echo "[+] Font removed"
}

case "$1" in
  apply) apply_font "$2" ;;
  remove) remove_font ;;
  *) echo "Usage: font.sh [apply <url|path>|remove]" ;;
esac
