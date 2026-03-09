#!/bin/bash
# 02-navigate.sh - Navigation and tab creation

source "$(dirname "$0")/common.sh"

start_test "Navigate to page"

# Navigate to fixtures index
RESULT=$(pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/\"}")
assert_json_contains "$RESULT" '.title' 'E2E Test' "title contains 'E2E Test'"
assert_json_contains "$RESULT" '.url' 'fixtures' "url contains 'fixtures'"

TAB_ID=$(echo "$RESULT" | jq -r '.tabId')
[ -n "$TAB_ID" ] && echo -e "  ${GREEN}✓${NC} Got tabId: ${TAB_ID:0:8}..."

end_test

start_test "Navigate to multiple pages"

# Navigate to form page
RESULT1=$(pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/form.html\"}")
assert_json_contains "$RESULT1" '.title' 'Form' "form page loaded"

# Navigate to table page
RESULT2=$(pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/table.html\"}")
assert_json_contains "$RESULT2" '.title' 'Table' "table page loaded"

end_test

start_test "List tabs"

# Get all tabs
TABS=$(pt_get "/tabs")
assert_json_length_gte "$TABS" '.tabs' 2 "at least 2 tabs exist"

# Verify tabs have required fields
FIRST_TAB=$(echo "$TABS" | jq '.tabs[0]')
assert_json_contains "$FIRST_TAB" '.id' '' "tab has id"
assert_json_contains "$FIRST_TAB" '.url' 'http' "tab has url"

end_test
