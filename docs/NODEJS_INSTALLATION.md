# Node.js Installation Guide

## Overview

This guide provides multiple methods to install the latest Node.js on Ubuntu Linux. For the Update Manager project, Node.js 14+ is required (Node.js 18+ recommended).

## Method 1: Using NVM (Node Version Manager) - Recommended

NVM allows you to easily switch between Node.js versions and is ideal for development.

### Installation Steps

1. **Install NVM:**
```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
```

2. **Reload your shell configuration:**
```bash
source ~/.bashrc
# Or if using zsh:
# source ~/.zshrc
```

3. **Verify NVM installation:**
```bash
nvm --version
```

4. **Install the latest LTS version of Node.js:**
```bash
nvm install --lts
```

5. **Or install the latest version:**
```bash
nvm install node
```

6. **Set as default:**
```bash
nvm use --default node
# Or for LTS:
# nvm use --default --lts
```

7. **Verify Node.js installation:**
```bash
node --version
npm --version
```

### Using NVM

```bash
# List installed versions
nvm list

# Install a specific version
nvm install 18.17.0

# Switch to a version
nvm use 18.17.0

# Set default version
nvm alias default 18.17.0
```

---

## Method 2: Using NodeSource Repository (Official)

This method installs Node.js system-wide using the official NodeSource repository.

### Installation Steps

1. **Update package index:**
```bash
sudo apt update
```

2. **Install required packages:**
```bash
sudo apt install -y curl gnupg2 ca-certificates
```

3. **Add NodeSource repository (for Node.js 20.x - latest LTS):**
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
```

4. **Install Node.js:**
```bash
sudo apt install -y nodejs
```

5. **Verify installation:**
```bash
node --version
npm --version
```

### For Different Versions

Replace `20.x` with your desired version:
- `setup_18.x` - Node.js 18 LTS
- `setup_20.x` - Node.js 20 LTS (current)
- `setup_lts.x` - Latest LTS
- `setup_current.x` - Latest current version

---

## Method 3: Using Snap (Simple but Limited)

Snap provides an easy way to install Node.js, but version selection is limited.

### Installation Steps

1. **Install Node.js via snap:**
```bash
sudo snap install node --classic
```

2. **Verify installation:**
```bash
node --version
npm --version
```

### Note
Snap may not always have the absolute latest version. Use `snap info node` to check available versions.

---

## Method 4: Using Package Manager (Ubuntu Default)

Ubuntu's default repositories may have older versions, but it's the simplest method.

### Installation Steps

1. **Update package index:**
```bash
sudo apt update
```

2. **Install Node.js and npm:**
```bash
sudo apt install -y nodejs npm
```

3. **Verify installation:**
```bash
node --version
npm --version
```

**Note:** This method typically installs an older version. Use this only if other methods don't work.

---

## Quick Installation Script (NVM Method)

Run this single command to install NVM and the latest Node.js LTS:

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash && \
source ~/.bashrc && \
nvm install --lts && \
nvm use --default --lts && \
node --version
```

---

## Verification

After installation, verify everything works:

```bash
# Check Node.js version (should be 14+)
node --version

# Check npm version
npm --version

# Check nvm version (if using nvm)
nvm --version

# Test Artillery (if already installed)
artillery --version
```

---

## Troubleshooting

### NVM not found after installation

If `nvm` command is not found after installation:

1. **Check if NVM is installed:**
```bash
ls -la ~/.nvm
```

2. **Add to your shell profile:**
```bash
# For bash
echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.bashrc
echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> ~/.bashrc
echo '[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"' >> ~/.bashrc
source ~/.bashrc

# For zsh
echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.zshrc
echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> ~/.zshrc
echo '[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"' >> ~/.zshrc
source ~/.zshrc
```

### Permission Denied Errors

If you get permission errors:

```bash
# Fix npm permissions (if using system-wide installation)
sudo chown -R $(whoami) ~/.npm
```

### Multiple Node.js Versions

If you have multiple Node.js installations:

1. **Check which Node.js is being used:**
```bash
which node
```

2. **If using NVM, ensure it's in your PATH:**
```bash
nvm use --default node
```

3. **Remove old installations:**
```bash
# If installed via apt
sudo apt remove nodejs npm

# If installed via snap
sudo snap remove node
```

---

## Recommended Setup for Update Manager

For the Update Manager project, we recommend:

1. **Install NVM** (Method 1) - Best for development
2. **Install Node.js 18 LTS or 20 LTS**
3. **Verify Artillery works:**
```bash
npm install -g artillery
artillery --version
```

---

## Updating Node.js

### Using NVM

```bash
# Install latest LTS
nvm install --lts

# Switch to it
nvm use --lts

# Set as default
nvm alias default lts/*
```

### Using NodeSource (apt)

```bash
# Update repository
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -

# Upgrade Node.js
sudo apt upgrade nodejs
```

### Using Snap

```bash
sudo snap refresh node
```

---

## Uninstalling Node.js

### If installed via NVM

```bash
# Remove NVM directory
rm -rf ~/.nvm
# Remove from shell profile (~/.bashrc or ~/.zshrc)
```

### If installed via apt

```bash
sudo apt remove nodejs npm
sudo apt autoremove
```

### If installed via snap

```bash
sudo snap remove node
```

### If installed via NodeSource

```bash
sudo apt remove nodejs npm
sudo rm -rf /etc/apt/sources.list.d/nodesource.list
sudo apt update
```

---

## Next Steps

After installing Node.js:

1. **Install Artillery globally:**
```bash
npm install -g artillery
```

2. **Verify load tests work:**
```bash
cd /home/accops/UpdateManager
make load-test
```

3. **Check Node.js version:**
```bash
node --version  # Should show v14+ or v18+ or v20+
```

