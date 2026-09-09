# suic1.de

# ⚠️ DO NOT RENAME ⚠️

Created by **RFour** ([ArchIsDead](https://github.com/ArchIsDead))

Profile: [dont.suic1.de](https://dont.suic1.de)

---

## Install

```bash
pkg update -y && pkg upgrade -y && pkg install make golang git mpv -y && git clone https://github.com/ArchIsDead/ihatemylife.git && cd ihatemylife && make install
```

---

Make Commands Guide

Command What It Does
make run Check deps, then run the tool directly
make start Check deps, then run the built binary
make build Check deps, then build the binary only
make install Check deps, tidy modules, build binary
make update Git pull, rebuild, and run
make clean Remove binary and clean Go cache
make fix Reset presets, rebuild binary
make stop Kill running instance
make restart Stop then start again

---

Quick Start

```bash
make run
```

Update To Latest

```bash
make update
```

Build Only

```bash
make build
```

Reset Everything

```bash
make fix
```

Stop Running

```bash
make stop
```

Restart

```bash
make restart
```

---

License

MIT License. See LICENSE for details.
