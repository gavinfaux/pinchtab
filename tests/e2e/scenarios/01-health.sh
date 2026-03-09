#!/bin/bash
# 01-health.sh — Basic connectivity tests

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab health"

pt_get /health
assert_json_eq "$RESULT" '.status' 'ok'

end_test

# ─────────────────────────────────────────────────────────────────
start_test "fixtures server"

# Verify test fixtures are accessible
curl -sf "${FIXTURES_URL}/" > /dev/null && echo -e "  ${GREEN}✓${NC} GET ${FIXTURES_URL}/" && ((ASSERTIONS_PASSED++)) || ((ASSERTIONS_FAILED++))
curl -sf "${FIXTURES_URL}/form.html" > /dev/null && echo -e "  ${GREEN}✓${NC} GET ${FIXTURES_URL}/form.html" && ((ASSERTIONS_PASSED++)) || ((ASSERTIONS_FAILED++))
curl -sf "${FIXTURES_URL}/table.html" > /dev/null && echo -e "  ${GREEN}✓${NC} GET ${FIXTURES_URL}/table.html" && ((ASSERTIONS_PASSED++)) || ((ASSERTIONS_FAILED++))

end_test
