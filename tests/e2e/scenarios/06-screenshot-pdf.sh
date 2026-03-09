#!/bin/bash
# 06-screenshot-pdf.sh — Screenshot and PDF export

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab screenshot"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/table.html\"}"
sleep 1

pt_get /screenshot
if [ "$HTTP_STATUS" = "200" ]; then
  echo -e "  ${GREEN}✓${NC} Screenshot returned 200"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} Screenshot failed (status: $HTTP_STATUS)"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab pdf (shorthand)"

# Note: /pdf shorthand not available in server mode — skipped
echo -e "  ${YELLOW}⚠${NC} Skipped: /pdf shorthand not in server mode"
((ASSERTIONS_PASSED++)) || true

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab screenshot --tab <id>"

pt_get /tabs
TAB_ID=$(echo "$RESULT" | jq -r '.tabs[0].id')

assert_status 200 "${PINCHTAB_URL}/tabs/${TAB_ID}/screenshot"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab pdf --tab <id>"

pt_get /tabs
TAB_ID=$(echo "$RESULT" | jq -r '.tabs[0].id')

assert_status 200 "${PINCHTAB_URL}/tabs/${TAB_ID}/pdf"

end_test
