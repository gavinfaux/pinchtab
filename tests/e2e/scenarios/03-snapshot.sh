#!/bin/bash
# 03-snapshot.sh — Accessibility tree and text extraction

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab snap"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/\"}"

pt_get /snapshot
assert_json_length_gte "$RESULT" '.nodes' 1
assert_json_contains "$RESULT" '.title' 'E2E Test'
assert_json_contains "$RESULT" '.url' 'fixtures'

# Check node structure
FIRST_NODE=$(echo "$RESULT" | jq '.nodes[0]')
assert_json_contains "$FIRST_NODE" '.ref' 'e'
assert_json_contains "$FIRST_NODE" '.role' ''

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab snap (interactive elements)"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/buttons.html\"}"
sleep 1

pt_get /snapshot

# Should find buttons
BUTTON_COUNT=$(echo "$RESULT" | jq '[.nodes[] | select(.role == "button")] | length')
if [ "$BUTTON_COUNT" -ge 3 ]; then
  echo -e "  ${GREEN}✓${NC} Found $BUTTON_COUNT buttons"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} Expected >=3 buttons, found $BUTTON_COUNT"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab text"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/table.html\"}"
sleep 1

pt_get /text
assert_json_contains "$RESULT" '.text' 'Alice Johnson'
assert_json_contains "$RESULT" '.text' 'bob@example.com'

end_test
