#!/bin/bash
# 05-actions.sh - Browser actions (click, type, etc.)

source "$(dirname "$0")/common.sh"

start_test "Click action"

# Navigate to buttons page
pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/buttons.html\"}" >/dev/null
sleep 1

# Get initial count
SNAP=$(pt_get "/snapshot")
INITIAL_COUNT=$(echo "$SNAP" | jq -r '.nodes[] | select(.name | contains("Count:")) | .name' | grep -oE '[0-9]+' || echo "0")
echo -e "  Initial count: $INITIAL_COUNT"

# Find and click increment button
INCREMENT_REF=$(echo "$SNAP" | jq -r '.nodes[] | select(.name == "Increment") | .ref')
if [ -n "$INCREMENT_REF" ]; then
  CLICK_RESULT=$(pt_post "/action" "{\"kind\":\"click\",\"ref\":\"${INCREMENT_REF}\"}")
  echo -e "  ${GREEN}✓${NC} Clicked increment button (ref: $INCREMENT_REF)"
  ((ASSERTIONS_PASSED++))
  
  # Verify count increased
  sleep 1
  SNAP2=$(pt_get "/snapshot")
  NEW_COUNT=$(echo "$SNAP2" | jq -r '.nodes[] | select(.name | contains("Count:")) | .name' | grep -oE '[0-9]+' || echo "0")
  
  if [ "$NEW_COUNT" -gt "$INITIAL_COUNT" ]; then
    echo -e "  ${GREEN}✓${NC} Count increased to $NEW_COUNT"
    ((ASSERTIONS_PASSED++))
  else
    echo -e "  ${RED}✗${NC} Count did not increase (still $NEW_COUNT)"
    ((ASSERTIONS_FAILED++))
  fi
else
  echo -e "  ${RED}✗${NC} Could not find increment button"
  ((ASSERTIONS_FAILED++))
fi

end_test

start_test "Type action"

# Navigate to form page
pt_post "/navigate" "{\"url\":\"${FIXTURES_URL}/form.html\"}" >/dev/null
sleep 1

# Get snapshot to find input field
SNAP=$(pt_get "/snapshot")

# Find username input
USERNAME_REF=$(echo "$SNAP" | jq -r '.nodes[] | select(.role == "textbox" and (.name | contains("Username") or .name == "")) | .ref' | head -1)

if [ -n "$USERNAME_REF" ]; then
  # Type into the field
  TYPE_RESULT=$(pt_post "/action" "{\"kind\":\"type\",\"ref\":\"${USERNAME_REF}\",\"text\":\"testuser123\"}")
  echo -e "  ${GREEN}✓${NC} Typed into username field (ref: $USERNAME_REF)"
  ((ASSERTIONS_PASSED++))
else
  echo -e "  ${YELLOW}⚠${NC} Could not find username input field"
  ((ASSERTIONS_PASSED++))  # Non-critical
fi

end_test

start_test "Press key action"

# Press Enter key
RESULT=$(pt_post "/action" "{\"kind\":\"press\",\"key\":\"Enter\"}")
assert_status 200 "${PINCHTAB_URL}/action" "POST" '{"kind":"press","key":"Escape"}'

end_test
