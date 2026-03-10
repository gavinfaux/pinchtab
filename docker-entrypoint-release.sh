#!/bin/sh
set -eu

home_dir="${HOME:-/data}"
xdg_config_home="${XDG_CONFIG_HOME:-$home_dir/.config}"
default_config_path="$xdg_config_home/pinchtab/config.json"

mkdir -p "$home_dir" "$xdg_config_home" "$(dirname "$default_config_path")"

if [ -z "${PINCHTAB_CONFIG:-}" ] && [ ! -f "$default_config_path" ]; then
  /usr/local/bin/pinchtab config init >/dev/null
  if [ -n "${PINCHTAB_TOKEN:-}" ]; then
    /usr/local/bin/pinchtab config set server.token "$PINCHTAB_TOKEN" >/dev/null
  fi
fi

# Docker port publishing needs a non-loopback bind inside the container, but
# keep the persisted config on its secure local default unless the user
# explicitly overrides it.
if [ -z "${PINCHTAB_CONFIG:-}" ] && [ -z "${PINCHTAB_BIND:-}" ]; then
  export PINCHTAB_BIND=0.0.0.0
fi

if [ -z "${PINCHTAB_CONFIG:-}" ] && [ -f "$default_config_path" ]; then
  current_flags="$(/usr/local/bin/pinchtab config get browser.extraFlags 2>/dev/null || true)"
  if [ -z "$current_flags" ]; then
    /usr/local/bin/pinchtab config set browser.extraFlags "--no-sandbox --disable-gpu" >/dev/null
  fi
fi

exec "$@"
