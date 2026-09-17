#!/data/data/com.termux/files/usr/bin/bash

set -e

echo "[*] Installing dependencies..."
pkg update -y
pkg install -y zsh git curl wget figlet toilet ncurses-utils nano

echo "[*] Installing Oh My Zsh..."
if [ ! -d "$HOME/.oh-my-zsh" ]; then
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
fi

echo "[*] Installing zsh plugins..."
ZSH_CUSTOM="$HOME/.oh-my-zsh/custom"
[ -d "$ZSH_CUSTOM/plugins/zsh-autosuggestions" ] || git clone https://github.com/zsh-users/zsh-autosuggestions "$ZSH_CUSTOM/plugins/zsh-autosuggestions"
[ -d "$ZSH_CUSTOM/plugins/zsh-syntax-highlighting" ] || git clone https://github.com/zsh-users/zsh-syntax-highlighting "$ZSH_CUSTOM/plugins/zsh-syntax-highlighting"

echo "[*] Setting zsh as default shell..."
chsh -s zsh

echo "[*] Installing Powerlevel10k..."
[ -d "$ZSH_CUSTOM/themes/powerlevel10k" ] || git clone --depth=1 https://github.com/romkatv/powerlevel10k.git "$ZSH_CUSTOM/themes/powerlevel10k"

mkdir -p "$HOME/.termux"

echo "[+] Installation complete."
echo "    Run: bash ~/.termux-theme/apply.sh <theme>"
ls "$HOME/.termux-theme/themes"
