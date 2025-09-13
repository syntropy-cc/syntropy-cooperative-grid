#!/bin/bash

# Test Runner for Syntropy Cooperative Grid

set -e

echo "Running Syntropy Cooperative Grid test suite..."

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

run_tests() {
    local test_type=$1
    local test_path=$2
    
    echo -e "${YELLOW}Running $test_type tests...${NC}"
    
    if [ -d "$test_path" ] && [ "$(find "$test_path" -name "*.go" | wc -l)" -gt 0 ]; then
        go test -v "$test_path/..."
        echo -e "${GREEN}✓ $test_type tests passed${NC}"
    else
        echo -e "${YELLOW}⚠ No $test_type tests found in $test_path${NC}"
    fi
}

# Run different test suites
run_tests "Unit" "./tests/unit"
run_tests "Integration" "./tests/integration"
run_tests "Security" "./tests/security"

# Placeholder for future test types
echo -e "${YELLOW}Note: End-to-end and performance tests will be added in future phases${NC}"

echo -e "${GREEN}All available tests completed successfully!${NC}"
