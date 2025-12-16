#!/bin/bash
# Quick test script for pending updates feature

BASE_URL="http://localhost:8080/api/v1"

echo "=== Testing Pending Updates Feature ==="
echo ""

# Test 1: Get deployment pending updates
echo "Test 1: Get deployment pending updates"
echo "GET $BASE_URL/customers/test-customer/tenants/test-tenant/deployments/test-deployment/updates"
curl -s "$BASE_URL/customers/test-customer/tenants/test-tenant/deployments/test-deployment/updates" | jq '.' || echo "Failed or no jq installed"
echo ""

# Test 2: Get tenant pending updates
echo "Test 2: Get tenant pending updates"
echo "GET $BASE_URL/customers/test-customer/tenants/test-tenant/deployments/pending-updates"
curl -s "$BASE_URL/customers/test-customer/tenants/test-tenant/deployments/pending-updates" | jq '.' || echo "Failed or no jq installed"
echo ""

# Test 3: Get customer pending updates
echo "Test 3: Get customer pending updates"
echo "GET $BASE_URL/customers/test-customer/deployments/pending-updates"
curl -s "$BASE_URL/customers/test-customer/deployments/pending-updates" | jq '.' || echo "Failed or no jq installed"
echo ""

# Test 4: Get all pending updates
echo "Test 4: Get all pending updates (admin view)"
echo "GET $BASE_URL/updates/pending"
curl -s "$BASE_URL/updates/pending" | jq '.' || echo "Failed or no jq installed"
echo ""

# Test 5: Filter by priority
echo "Test 5: Filter by priority (critical)"
echo "GET $BASE_URL/updates/pending?priority=critical"
curl -s "$BASE_URL/updates/pending?priority=critical" | jq '.' || echo "Failed or no jq installed"
echo ""

echo "=== Tests Complete ==="
