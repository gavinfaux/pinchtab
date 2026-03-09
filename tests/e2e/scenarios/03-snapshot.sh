#!/bin/bash
# 03-snapshot.sh - Accessibility tree extraction

source "$(dirname "$0")/common.sh"

start_test "Snapshot current page"

# Navigate to a page first
pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/\"}" >/dev/null

# Get snapshot
SNAP=$(pt_get "/snapshot")
assert_json_length_gte "$SNAP" '.nodes' 1 "snapshot has nodes"
assert_json_contains "$SNAP" '.title' 'E2E Test' "snapshot has title"
assert_json_contains "$SNAP" '.url' 'fixtures' "snapshot has url"

# Check node structure
FIRST_NODE=$(echo "$SNAP" | jq '.nodes[0]')
assert_json_contains "$FIRST_NODE" '.ref' 'e' "node has ref"
assert_json_contains "$FIRST_NODE" '.role' '' "node has role"

end_test

start_test "Snapshot with elements"

# Navigate to buttons page (has interactive elements)
pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/buttons.html\"}" >/dev/null
sleep 1

SNAP=$(pt_get "/snapshot")

# Should find buttons
BUTTON_COUNT=$(echo "$SNAP" | jq '[.nodes[] | select(.role == "button")] | length')
if [ "$BUTTON_COUNT" -ge 3 ]; then
  echo -e "  ${GREEN}✓${NC} Found $BUTTON_COUNT buttons"
  ((ASSERTIONS_PASSED++))
else
  echo -e "  ${RED}✗${NC} Expected >=3 buttons, found $BUTTON_COUNT"
  ((ASSERTIONS_FAILED++))
fi

end_test

start_test "Text extraction"

# Navigate to table page
pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/table.html\"}" >/dev/null
sleep 1

# Get text content
TEXT=$(pt_get "/text")
assert_json_contains "$TEXT" '.text' 'Alice Johnson' "text contains 'Alice Johnson'"
assert_json_contains "$TEXT" '.text' 'bob@example.com' "text contains 'bob@example.com'"

end_test
