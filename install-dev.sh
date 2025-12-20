#!/bin/bash
# Install the extension to the local azd extensions directory for development

set -e

echo "Building extension..."
go build -o azd-ext-doctor .

echo "Installing to ~/.azd/extensions/spboyer.azd.doctor/..."
cp azd-ext-doctor ~/.azd/extensions/spboyer.azd.doctor/

echo "✓ Extension installed successfully!"
echo ""
echo "Test with: azd doctor check"
