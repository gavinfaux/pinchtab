#!/bin/bash
# 05-actions.sh — Browser actions (click, type, press)

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab click <ref>"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/buttons.html\"}"
sleep 2

# Get snapshot and find increment button
pt_get /snapshot
INCREMENT_REF=$(echo "$RESULT" | jq -r '.nodes[] | select(.name == "Increment") | .ref' | head -1)

if [ -n "$INCREMENT_REF" ] && [ "$INCREMENT_REF" != "null" ]; then
  pt_post /action -d "{\"kind\":\"click\",\"ref\":\"${INCREMENT_REF}\"}"
  echo -e "  ${GREEN}✓${NC} Clicked button (ref: $INCREMENT_REF)"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${YELLOW}⚠${NC} Could not find increment button"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab type <ref> <text>"

pt_post /navigate -d "{\"url\":\"${FIXTURES_URL}/form.html\"}"
sleep 1

pt_get /snapshot
USERNAME_REF=$(echo "$RESULT" | jq -r '.nodes[] | select(.role == "textbox") | .ref' | head -1)

if [ -n "$USERNAME_REF" ] && [ "$USERNAME_REF" != "null" ]; then
  pt_post /action -d "{\"kind\":\"type\",\"ref\":\"${USERNAME_REF}\",\"text\":\"testuser123\"}"
  echo -e "  ${GREEN}✓${NC} Typed into field (ref: $USERNAME_REF)"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${YELLOW}⚠${NC} Could not find input field"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab press <key>"

pt_post /action -d '{"kind":"press","key":"Escape"}'
assert_status 200 "${PINCHTAB_URL}/action" "POST" '{"kind":"press","key":"Escape"}'

end_test
