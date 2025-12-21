#!/bin/bash
# TAPI Quick Start Script

set -e

echo "🚀 TAPI - Terminal API Explorer Quick Start"
echo ""

# Check if binary exists
if [ ! -f "./bin/tapi" ]; then
    echo "📦 Building TAPI..."
    make build
    echo "✅ Build complete!"
    echo ""
fi

# Display options
echo "Choose an option:"
echo ""
echo "  1) Explore example Pet Store API (local)"
echo "  2) Explore Pet Store API (remote)"
echo "  3) Validate example spec"
echo "  4) Show help"
echo "  5) Exit"
echo ""
read -p "Enter your choice (1-5): " choice

case $choice in
    1)
        echo ""
        echo "🎯 Launching TAPI with local Pet Store spec..."
        echo "💡 Tip: Use j/k to navigate, Enter to select, ? for help, q to quit"
        echo ""
        sleep 2
        ./bin/tapi explore -f example-petstore.yaml
        ;;
    2)
        echo ""
        echo "🌐 Fetching remote Pet Store spec..."
        echo "💡 Tip: Use j/k to navigate, Enter to select, ? for help, q to quit"
        echo ""
        sleep 2
        ./bin/tapi explore -u https://petstore3.swagger.io/api/v3/openapi.json
        ;;
    3)
        echo ""
        echo "✅ Validating example-petstore.yaml..."
        echo ""
        ./bin/tapi validate -f example-petstore.yaml
        ;;
    4)
        echo ""
        ./bin/tapi --help
        ;;
    5)
        echo ""
        echo "👋 Goodbye!"
        exit 0
        ;;
    *)
        echo ""
        echo "❌ Invalid choice. Please run again and select 1-5."
        exit 1
        ;;
esac
