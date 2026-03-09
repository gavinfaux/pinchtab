#!/bin/bash
# Common utilities for E2E tests

set -uo pipefail
# Note: not using -e because arithmetic operations can return 1

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Defaults from environment
PINCHTAB_URL="${PINCHTAB_URL:-http://localhost:9999}"
FIXTURES_URL="${FIXTURES_URL:-http://localhost:8080}"
RESULTS_DIR="${RESULTS_DIR:-/results}"

# Test tracking
CURRENT_TEST=""
TESTS_PASSED=0
TESTS_FAILED=0
ASSERTIONS_PASSED=0
ASSERTIONS_FAILED=0

# Start a test
start_test() {
  CURRENT_TEST="$1"
  echo -e "${BLUE}▶ ${CURRENT_TEST}${NC}"
}

# End a test
end_test() {
  if [ "$ASSERTIONS_FAILED" -eq 0 ]; then
    echo -e "${GREEN}✓ ${CURRENT_TEST} passed${NC}\n"
    ((TESTS_PASSED++)) || true
  else
    echo -e "${RED}✗ ${CURRENT_TEST} failed${NC}\n"
    ((TESTS_FAILED++)) || true
  fi
  ASSERTIONS_PASSED=0
  ASSERTIONS_FAILED=0
}

# Assert HTTP status
assert_status() {
  local expected="$1"
  local url="$2"
  local method="${3:-GET}"
  local body="${4:-}"
  
  local actual
  if [ -n "$body" ]; then
    actual=$(curl -s -o /dev/null -w '%{http_code}' -X "$method" -H "Content-Type: application/json" -d "$body" "$url")
  else
    actual=$(curl -s -o /dev/null -w '%{http_code}' -X "$method" "$url")
  fi
  
  if [ "$actual" = "$expected" ]; then
    echo -e "  ${GREEN}✓${NC} $method $url → $actual"
    ((ASSERTIONS_PASSED++)) || true
  else
    echo -e "  ${RED}✗${NC} $method $url → $actual (expected $expected)"
    ((ASSERTIONS_FAILED++)) || true
  fi
}

# Assert command succeeds (exit 0)
assert_ok() {
  local desc="$1"
  shift
  
  if "$@" >/dev/null 2>&1; then
    echo -e "  ${GREEN}✓${NC} $desc"
    ((ASSERTIONS_PASSED++)) || true
  else
    echo -e "  ${RED}✗${NC} $desc (exit $?)"
    ((ASSERTIONS_FAILED++)) || true
  fi
}

# Assert JSON field equals value
assert_json_eq() {
  local json="$1"
  local path="$2"
  local expected="$3"
  local desc="${4:-$path = $expected}"
  
  local actual
  actual=$(echo "$json" | jq -r "$path")
  
  if [ "$actual" = "$expected" ]; then
    echo -e "  ${GREEN}✓${NC} $desc"
    ((ASSERTIONS_PASSED++)) || true
  else
    echo -e "  ${RED}✗${NC} $desc (got: $actual)"
    ((ASSERTIONS_FAILED++)) || true
  fi
}

# Assert JSON field contains value
assert_json_contains() {
  local json="$1"
  local path="$2"
  local needle="$3"
  local desc="${4:-$path contains '$needle'}"
  
  local actual
  actual=$(echo "$json" | jq -r "$path")
  
  if [[ "$actual" == *"$needle"* ]]; then
    echo -e "  ${GREEN}✓${NC} $desc"
    ((ASSERTIONS_PASSED++)) || true
  else
    echo -e "  ${RED}✗${NC} $desc (got: $actual)"
    ((ASSERTIONS_FAILED++)) || true
  fi
}

# Assert JSON array length
assert_json_length() {
  local json="$1"
  local path="$2"
  local expected="$3"
  local desc="${4:-$path length = $expected}"
  
  local actual
  actual=$(echo "$json" | jq "$path | length")
  
  if [ "$actual" -eq "$expected" ]; then
    echo -e "  ${GREEN}✓${NC} $desc"
    ((ASSERTIONS_PASSED++)) || true
  else
    echo -e "  ${RED}✗${NC} $desc (got: $actual)"
    ((ASSERTIONS_FAILED++)) || true
  fi
}

# Assert JSON array length >= value
assert_json_length_gte() {
  local json="$1"
  local path="$2"
  local expected="$3"
  local desc="${4:-$path length >= $expected}"
  
  local actual
  actual=$(echo "$json" | jq "$path | length")
  
  if [ "$actual" -ge "$expected" ]; then
    echo -e "  ${GREEN}✓${NC} $desc"
    ((ASSERTIONS_PASSED++)) || true
  else
    echo -e "  ${RED}✗${NC} $desc (got: $actual)"
    ((ASSERTIONS_FAILED++)) || true
  fi
}

# HTTP helpers
pt_get() {
  curl -s "${PINCHTAB_URL}$1"
}

pt_post() {
  curl -s -X POST -H "Content-Type: application/json" -d "$2" "${PINCHTAB_URL}$1"
}

# Print summary
print_summary() {
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo -e "${BLUE}E2E Test Summary${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo -e "Passed: ${GREEN}${TESTS_PASSED}${NC}"
  echo -e "Failed: ${RED}${TESTS_FAILED}${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  if [ "$TESTS_FAILED" -gt 0 ]; then
    exit 1
  fi
}
