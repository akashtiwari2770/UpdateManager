#!/bin/bash
# Script to fix Docker permissions by adding user to docker group

echo "Fixing Docker permissions..."
echo ""

# Check if user is already in docker group
if groups | grep -q docker; then
    echo "✓ User is already in docker group"
    echo "  If you still get permission errors, try: newgrp docker"
    exit 0
fi

# Check if docker group exists
if ! getent group docker > /dev/null 2>&1; then
    echo "Error: docker group does not exist"
    echo "Creating docker group..."
    sudo groupadd docker
fi

# Add user to docker group
echo "Adding user to docker group..."
sudo usermod -aG docker $USER

echo ""
echo "✓ User added to docker group"
echo ""
echo "IMPORTANT: You need to either:"
echo "  1. Log out and log back in, OR"
echo "  2. Run: newgrp docker"
echo ""
echo "After that, test with: docker ps"
echo ""
echo "Would you like to apply the changes now? (y/n)"
read -r response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    echo "Applying changes..."
    newgrp docker <<EOF
echo "Docker group activated in this shell"
docker ps
EOF
    echo "If docker ps worked above, you're all set!"
else
    echo "Please log out and log back in for changes to take effect."
fi

