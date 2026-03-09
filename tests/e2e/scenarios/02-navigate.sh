#!/bin/bash
# 02-navigate.sh — Navigation and tab management

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab nav <url>"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/\"}"
assert_json_contains "$RESULT" '.title' 'E2E Test'
assert_json_contains "$RESULT" '.url' 'fixtures'

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab nav (multiple pages)"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/form.html\"}"
assert_json_contains "$RESULT" '.title' 'Form'

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/table.html\"}"
assert_json_contains "$RESULT" '.title' 'Table'

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tabs"

pt_get /tabs
assert_json_length_gte "$RESULT" '.tabs' 2

# Verify tab structure
FIRST_TAB=$(echo "$RESULT" | jq '.tabs[0]')
assert_json_contains "$FIRST_TAB" '.id' ''
assert_json_contains "$FIRST_TAB" '.url' 'http'

end_test
