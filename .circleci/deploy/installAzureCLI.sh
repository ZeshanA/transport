#!/usr/bin/env bash

echo "Installing azure-cli..."

# Get install dependencies
sudo apt-get install apt-transport-https lsb-release

# Modify APT sources list to include azure-cli repo
export AZ_REPO=$(lsb_release -cs)
echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" | \
sudo tee /etc/apt/sources.list.d/azure-cli.list

# Get Microsoft signing key
sudo apt-key --keyring /etc/apt/trusted.gpg.d/Microsoft.gpg adv \
    --keyserver packages.microsoft.com \
    --recv-keys BC528686B50D79E339D3721CEB3E94ADBE1229CF

# Update sources and get azure-cli packages
sudo apt-get update && sudo apt-get install azure-cli && echo "azure-cli installed successfully"