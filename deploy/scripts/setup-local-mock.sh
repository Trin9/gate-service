#!/bin/bash
echo "Creating mock FSx directory at /tmp/mock-fsx..."
mkdir -p /tmp/mock-fsx/qwen/Qwen1.5-4B-Chat

echo "Touching config.json..."
touch /tmp/mock-fsx/qwen/Qwen1.5-4B-Chat/config.json

echo "Setting permissions..."
chmod -R 777 /tmp/mock-fsx

echo "Done. HostPath ready."