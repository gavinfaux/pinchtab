#!/bin/bash
# Common utilities for E2E tests

set -uo pipefail
# Note: not using -e because arithmetic operations can return 1

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
MUTED='\033[0;90m'
NC='\033[0m'

# Defaults from environment
PINCHTAB_URL="${PINCHTAB_URL:-http://localhost:9999}"
FIXTURES_URL="${FIXTURES_URL:-http://localhost:8080}"
RESULTS_DIR="${RESULTS_DIR:-/results}"

# Test tracking (only initialize if not already set)
CURRENT_TEST="${CURRENT_TEST:-}"
TESTS_PASSED="${TESTS_PASSED:-0}"
TESTS_FAILED="${TESTS_FAILED:-0}"
ASSERTIONS_PASSED="${ASSERTIONS_PASSED:-0}"
ASSERTIONS_FAILED="${ASSERTIONS_FAILED:-0}"

# Test timing (using seconds, Alpine doesn't support ms)
TEST_START_TIME="${TEST_START_TIME:-0}"
TEST_START_NS="${TEST_START_NS:-0}"
if [ -z "${TEST_RESULTS_INIT:-}" ]; then
  TEST_RESULTS=()
  TEST_RESULTS_INIT=1
fi

# Get time in milliseconds (cross-platform)
get_time_ms() {
  if [ -f /proc/uptime ]; then
    # Linux: use /proc/uptime (gives centiseconds)
    awk '{printf "%.0f", $1 * 1000}' /proc/uptime
  elif command -v gdate &>/dev/null; then
    # macOS with coreutils
    gdate +%s%3N
  elif command -v perl &>/dev/null; then
    # Perl fallback
    perl -MTime::HiRes=time -e 'printf "%.0f", time * 1000'
  else
    # Last resort: seconds * 1000
    echo $(($(date +%s) * 1000))
  fi
}

# Start a test
start_test() {
  CURRENT_TEST="$1"
  TEST_START_TIME=$(get_time_ms)
  echo -e "${BLUE}▶ ${CURRENT_TEST}${NC}"
}

# End a test
end_test() {
  local end_time=$(get_time_ms)
  local duration=$((end_time - TEST_START_TIME))
  
  if [ "$ASSERTIONS_FAILED" -eq 0 ]; then
    echo -e "${GREEN}✓ ${CURRENT_TEST} passed${NC} ${MUTED}(${duration}ms)${NC}\n"
    TEST_RESULTS+=("✅ ${CURRENT_TEST}|${duration}ms|passed")
    ((TESTS_PASSED++)) || true
  else
    echo -e "${RED}✗ ${CURRENT_TEST} failed${NC} ${MUTED}(${duration}ms)${NC}\n"
    TEST_RESULTS+=("❌ ${CURRENT_TEST}|${duration}ms|failed")
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

# ================================================================
# Visible curl wrapper — shows exact command when running
# ================================================================

RESULT=""
HTTP_STATUS=""

pinchtab() {
  local method="$1"
  local path="$2"
  shift 2

  # Print the curl command in cyan so you see what's executed
  echo -e "${BLUE}→ curl -X $method ${PINCHTAB_URL}$path $@${NC}" >&2

  # Execute and capture response + status
  local response
  response=$(curl -s -w "\n%{http_code}" \
    -X "$method" \
    "${PINCHTAB_URL}$path" \
    -H "Content-Type: application/json" \
    "$@")

  RESULT=$(echo "$response" | head -n -1)
  HTTP_STATUS=$(echo "$response" | tail -n 1)
}

# Aliases for cleaner test files
pt() { pinchtab "$@"; }
pt_get() { pinchtab GET "$1"; echo "$RESULT"; }
pt_post() { pinchtab POST "$1" -d "$2"; echo "$RESULT"; }

# Print summary
print_summary() {
  local total=$((TESTS_PASSED + TESTS_FAILED))
  local total_time=0
  
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo -e "${BLUE}E2E Test Summary${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""
  printf "  %-40s %10s %10s\n" "Test" "Duration" "Status"
  echo "  ────────────────────────────────────────────────────────"
  
  for result in "${TEST_RESULTS[@]}"; do
    IFS='|' read -r name duration status <<< "$result"
    local time_num=${duration%ms}
    ((total_time += time_num)) || true
    if [ "$status" = "passed" ]; then
      printf "  %-40s %10s ${GREEN}%10s${NC}\n" "$name" "$duration" "✓"
    else
      printf "  %-40s %10s ${RED}%10s${NC}\n" "$name" "$duration" "✗"
    fi
  done
  
  echo "  ────────────────────────────────────────────────────────"
  printf "  %-40s %10s\n" "Total" "${total_time}ms"
  echo ""
  echo -e "  ${GREEN}Passed:${NC} ${TESTS_PASSED}/${total}"
  echo -e "  ${RED}Failed:${NC} ${TESTS_FAILED}/${total}"
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  # Generate markdown report for CI
  if [ -d "${RESULTS_DIR:-}" ]; then
    generate_markdown_report > "${RESULTS_DIR}/report.md"
    echo "passed=$TESTS_PASSED" > "${RESULTS_DIR}/summary.txt"
    echo "failed=$TESTS_FAILED" >> "${RESULTS_DIR}/summary.txt"
    echo "total_time=${total_time}ms" >> "${RESULTS_DIR}/summary.txt"
    echo "timestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "${RESULTS_DIR}/summary.txt"
  fi
  
  if [ "$TESTS_FAILED" -gt 0 ]; then
    exit 1
  fi
}

# Generate markdown report
generate_markdown_report() {
  local total=$((TESTS_PASSED + TESTS_FAILED))
  local total_time=0
  
  echo "## 🦀 PinchTab E2E Test Report"
  echo ""
  if [ "$TESTS_FAILED" -eq 0 ]; then
    echo "**Status:** ✅ All tests passed"
  else
    echo "**Status:** ❌ ${TESTS_FAILED} test(s) failed"
  fi
  echo ""
  echo "| Test | Duration | Status |"
  echo "|------|----------|--------|"
  
  for result in "${TEST_RESULTS[@]}"; do
    IFS='|' read -r name duration status <<< "$result"
    local time_num=${duration%ms}
    ((total_time += time_num)) || true
    local icon="✅"
    [ "$status" = "failed" ] && icon="❌"
    # Remove emoji from name for cleaner table
    local clean_name="${name#✅ }"
    clean_name="${clean_name#❌ }"
    echo "| ${clean_name} | ${duration} | ${icon} |"
  done
  
  echo ""
  echo "**Summary:** ${TESTS_PASSED}/${total} passed in ${total_time}ms"
  echo ""
  echo "<sub>Generated at $(date -u +%Y-%m-%dT%H:%M:%SZ)</sub>"
}
