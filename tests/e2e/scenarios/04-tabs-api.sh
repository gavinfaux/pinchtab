#!/bin/bash
# 04-tabs-api.sh - Tab-scoped API tests (regression test for #207)

source "$(dirname "$0")/common.sh"

start_test "Tab-scoped snapshot (regression #207)"

# Navigate to create a tab
RESULT=$(pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/index.html\"}")
CREATED_TAB=$(echo "$RESULT" | jq -r '.tabId')
echo -e "  Created tab: ${CREATED_TAB:0:12}..."

# Get tab ID from /tabs endpoint
TABS=$(pt_get "/tabs")
LISTED_TAB=$(echo "$TABS" | jq -r '.tabs[-1].id')
echo -e "  Listed tab:  ${LISTED_TAB:0:12}..."

# Test: /tabs/{id}/snapshot should work
assert_status 200 "${PINCHTAB_URL}/tabs/${LISTED_TAB}/snapshot"

# Verify snapshot content
SNAP=$(pt_get "/tabs/${LISTED_TAB}/snapshot")
assert_json_contains "$SNAP" '.title' 'E2E Test' "snapshot has correct title"

end_test

start_test "Tab-scoped operations"

# Navigate to form page
RESULT=$(pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/form.html\"}")
TAB_ID=$(echo "$RESULT" | jq -r '.tabId')

# Test /tabs/{id}/text
assert_status 200 "${PINCHTAB_URL}/tabs/${TAB_ID}/text"

# Test /tabs/{id}/screenshot
assert_status 200 "${PINCHTAB_URL}/tabs/${TAB_ID}/screenshot"

end_test

start_test "Tab close"

# Get current tabs
BEFORE=$(pt_get "/tabs" | jq '.tabs | length')

# Navigate to create a new tab
RESULT=$(pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/buttons.html\"}")
TAB_ID=$(echo "$RESULT" | jq -r '.tabId')

# Close the tab via /tabs/{id}/close
assert_status 200 "${PINCHTAB_URL}/tabs/${TAB_ID}/close" "POST"

# Verify tab count decreased or stayed same (close removes tab)
sleep 1
AFTER=$(pt_get "/tabs" | jq '.tabs | length')
if [ "$AFTER" -le "$BEFORE" ]; then
  echo -e "  ${GREEN}✓${NC} Tab closed successfully (before: $BEFORE, after: $AFTER)"
  ((ASSERTIONS_PASSED++))
else
  echo -e "  ${RED}✗${NC} Tab not closed (before: $BEFORE, after: $AFTER)"
  ((ASSERTIONS_FAILED++))
fi

end_test
