# Docker Permission Setup

## Quick Fix (Temporary - Uses sudo)

The Makefile now automatically detects if you need sudo and will use it. You can run:

```bash
make db-start
```

It will automatically use `sudo` if needed.

## Permanent Fix (Recommended)

To avoid using sudo every time, add your user to the docker group:

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Apply the changes (you need to log out and log back in)
# OR use newgrp to apply immediately
newgrp docker

# Verify it works
docker ps
```

After this, you won't need sudo for docker commands.

## Verify Docker Access

```bash
# Test without sudo
docker ps

# If it works, you're all set!
# If you get permission denied, use the permanent fix above
```

## Alternative: Use sudo explicitly

If you prefer to always use sudo, you can modify the Makefile or run commands directly:

```bash
cd src/database
sudo docker-compose -f docker-compose.mongodb.yml up -d
```

