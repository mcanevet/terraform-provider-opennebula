# Plugin Framework Migration Testing

This document describes the testing strategy for migrating resources from SDKv2 to the Plugin Framework.

## Test Types

### 1. Unit Tests
**Purpose**: Verify Framework implementation correctness without requiring OpenNebula

**Location**: `*_unit_test.go` files

**Run**:
```bash
go test ./opennebula -v -run "TestZoneDataSource" -short
```

**What they test**:
- Schema definitions
- Type conversions (SDKv2 types → Framework types)
- Null value handling
- Data model correctness

### 2. Migration Tests
**Purpose**: Verify SDKv2 and Framework implementations produce identical results

**Location**: `*_migration_test.go` files (require `acceptance` build tag)

**Run**:
```bash
# Requires running OpenNebula instance
go test ./opennebula -v -tags=acceptance -run "TestMigrate"
```

**What they test**:
- Apply config with SDKv2 provider
- Switch to Framework provider with same config
- Verify no plan differences (`plancheck.ExpectEmptyPlan()`)
- Switch back to SDKv2 and verify still no changes

### 3. Acceptance Tests
**Purpose**: Test Framework resources against real OpenNebula

**Location**: `TestAcc*` functions

**Run**:
```bash
# Requires OpenNebula with credentials
TF_ACC=1 go test ./opennebula -v -tags=acceptance -run "TestAccZone"
```

## Testing Requirements

### For Unit Tests
- ✅ No external dependencies
- ✅ Fast execution
- ✅ Run in CI/CD

### For Migration Tests
- ⚠️ Requires running OpenNebula instance
- ⚠️ Needs valid credentials
- ⚠️ Run manually or in integration environment

**Setup**:
```bash
# Start OpenNebula (example with Docker)
docker run -d --name opennebula \
  -p 2633:2633 \
  opennebula/opennebula:latest

# Configure credentials
export OPENNEBULA_ENDPOINT="http://localhost:2633/RPC2"
export OPENNEBULA_USERNAME="oneadmin"
export OPENNEBULA_PASSWORD="opennebula"
```

## Migration Test Pattern

Based on [HashiCorp's migration testing guide](https://developer.hashicorp.com/terraform/plugin/framework/migrating/testing):

```go
func TestMigrateResource_IdenticalBehavior(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            {
                // Step 1: SDKv2 provider
                ProtoV5ProviderFactories: sdkv2ProviderFactories,
                Config: testConfig,
                Check: resource.ComposeTestCheckFunc(
                    // Verify resource created correctly
                ),
            },
            {
                // Step 2: Framework provider - expect no changes
                ProtoV5ProviderFactories: frameworkProviderFactories,
                Config: testConfig,
                ConfigPlanChecks: resource.ConfigPlanChecks{
                    PreApply: []plancheck.PlanCheck{
                        plancheck.ExpectEmptyPlan(), // ← Key assertion
                    },
                },
            },
        },
    })
}
```

## Current Test Coverage

### opennebula_zone (Data Source)

**Unit Tests** ✅:
- Schema definition
- Type conversion (Int64, String)
- Null handling
- Model instantiation

**Migration Tests** ⚠️ (require OpenNebula):
- SDKv2 → Framework transition
- Framework → SDKv2 transition
- Filter by name
- Filter by ID

## Running Tests Locally

### Quick Check (Unit Tests Only)
```bash
# Fast, no dependencies
go test ./opennebula -v -short
```

### Full Migration Tests
```bash
# 1. Start OpenNebula
docker-compose up -d

# 2. Run migration tests
go test ./opennebula -v -tags=acceptance -run "TestMigrate"

# 3. Run all acceptance tests
TF_ACC=1 go test ./opennebula -v -tags=acceptance -timeout 30m
```

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Unit Tests
  run: go test ./opennebula -v -short

- name: Migration Tests
  run: |
    docker-compose up -d opennebula
    sleep 30  # Wait for OpenNebula to start
    go test ./opennebula -v -tags=acceptance -run "TestMigrate"
```

## Troubleshooting

### Flag Redefinition Error
```
panic: flag redefined: sweep
```

**Cause**: SDKv2 and terraform-plugin-testing define conflicting flags

**Solution**: Migration tests use `// +build acceptance` tag to separate them

### Empty Plan Check Fails
```
Expected an empty plan, but got changes
```

**Possible causes**:
1. Framework implementation differs from SDKv2
2. Type conversion issue (e.g., int vs int64)
3. Null handling mismatch
4. Computed attribute differences

**Debug**:
```bash
# Enable detailed logging
TF_LOG=DEBUG go test ./opennebula -v -tags=acceptance -run "TestMigrate"
```

## References

- [HashiCorp Migration Testing Guide](https://developer.hashicorp.com/terraform/plugin/framework/migrating/testing)
- [Plugin Testing Module](https://developer.hashicorp.com/terraform/plugin/testing)
- [Plan Checks Documentation](https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/plan-checks)

Sources:
- [Testing migration | Terraform](https://developer.hashicorp.com/terraform/plugin/framework/migrating/testing)
- [Plugin Development: Plan Checks](https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/plan-checks)
- [Terraform AWS Provider - Plugin Migrations](https://hashicorp.github.io/terraform-provider-aws/terraform-plugin-migrations/)
