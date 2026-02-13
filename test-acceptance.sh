#!/bin/bash
set -e

echo "🚀 Starting OpenNebula test environment (VM)..."
echo ""

# Check if VM exists
if vagrant status | grep -q "running"; then
    echo "✅ VM already running"
else
    echo "📦 Starting VM (first run takes ~10 minutes)..."
    vagrant up
fi

echo ""
echo "🧪 Running acceptance tests..."
echo ""

# Set environment variables
export OPENNEBULA_ENDPOINT="http://localhost:2633/RPC2"
export OPENNEBULA_USERNAME="oneadmin"
export OPENNEBULA_PASSWORD="opennebula"
export OPENNEBULA_FLOW_ENDPOINT="http://localhost:2474"
export TF_ACC=1

# Run migration tests
go test ./opennebula -v -tags=acceptance -run TestMigrate -timeout 30m

echo ""
echo "✅ Tests complete!"
echo ""
echo "Cleanup:"
echo "  vagrant halt     # Stop VM (keep it for later)"
echo "  vagrant destroy  # Delete VM completely"
