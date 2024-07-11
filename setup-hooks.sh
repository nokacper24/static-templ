#!/bin/bash

# Set the hooks directory to the `.githooks` folder in the repository
git config core.hooksPath .githooks

# Make the hooks executable
chmod +x .githooks/commit-msg
chmod +x .githooks/post-commit

echo "Git hooks have been set up."
