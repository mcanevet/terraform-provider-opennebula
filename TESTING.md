# Testing Guide

## Quick Start

### Unit Tests (No OpenNebula Required) ✅
```bash
go test ./opennebula -v -short
```

### Acceptance Tests (Requires OpenNebula)

#### Option 1: Automated VM Setup (Recommended)
```bash
./test-acceptance.sh
```

This will:
1. Start Ubuntu VM with Vagrant
2. Install OpenNebula 6.10 via miniONE
3. Run migration tests
4. ~10 minutes first run, ~2 minutes after

#### Option 2: Manual VM Setup
```bash
# Start VM
vagrant up

# Run tests from host
export OPENNEBULA_ENDPOINT="http://localhost:2633/RPC2"
export OPENNEBULA_USERNAME="oneadmin"
export OPENNEBULA_PASSWORD="opennebula"
export OPENNEBULA_FLOW_ENDPOINT="http://localhost:2474"
export TF_ACC=1

go test ./opennebula -v -tags=acceptance -run TestMigrate -timeout 30m
```

## Cleanup

```bash
vagrant halt      # Stop VM (keep for later)
vagrant destroy   # Delete VM completely
```

## Requirements

- **Vagrant** 2.0+ ([install](https://www.vagrantup.com/downloads))
- **VirtualBox** 6.1+ ([install](https://www.virtualbox.org/wiki/Downloads))
- 4GB RAM available for VM
- 20GB disk space

## What Gets Tested

### Unit Tests
- Schema definitions
- Type conversions
- Null handling
- Model correctness

### Migration Tests
- SDKv2 → Framework transition (no plan changes)
- Framework → SDKv2 transition (no plan changes)
- Filter by name/ID
- `plancheck.ExpectEmptyPlan()` assertions

## Troubleshooting

### VM won't start
```bash
# Check VirtualBox is running
vboxmanage --version

# Check VM status
vagrant status

# Destroy and recreate
vagrant destroy -f && vagrant up
```

### Tests fail to connect
```bash
# Verify ports are forwarded
vagrant port

# Should show:
# 2633 (guest) => 2633 (host)
# 2474 (guest) => 2474 (host)

# Test connection manually
curl http://localhost:2633/RPC2
```

### Out of disk space
```bash
# Check VM size
vagrant ssh -c "df -h"

# Increase disk in Vagrantfile if needed
```

## CI/CD

GitHub Actions runs full acceptance tests automatically:
- Installs OpenNebula 6.10 and 7.0
- Runs all tests
- No local setup needed for contributors

See `.github/workflows/ci.yaml` for details.
