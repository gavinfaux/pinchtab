#!/bin/bash
# 01-health.sh - Basic connectivity tests

source "$(dirname "$0")/common.sh"

start_test "Health endpoint"

# Test /health returns 200
assert_status 200 "${PINCHTAB_URL}/health"

# Test /health returns valid JSON
HEALTH=$(pt_get "/health")
assert_json_eq "$HEALTH" '.status' 'ok' "status = ok"

end_test

start_test "Fixtures server"

# Test fixtures are accessible
assert_status 200 "${FIXTURES_URL}/"
assert_status 200 "${FIXTURES_URL}/form.html"
assert_status 200 "${FIXTURES_URL}/table.html"
assert_status 200 "${FIXTURES_URL}/buttons.html"

end_test
