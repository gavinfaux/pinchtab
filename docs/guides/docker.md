# Docker Deployment

PinchTab can run in Docker with a mounted data volume for config, profiles, and state.
The bundled image now manages its default config under `/data/.config/pinchtab/config.json`.
If you want full control over the config file path, you can still mount your own file and point `PINCHTAB_CONFIG` at it.

## Quick Start

Build the image from this repository:

```bash
docker build -t pinchtab .
```

Run the container with a persistent data volume:

```bash
docker run -d \
  --name pinchtab \
  -p 127.0.0.1:9867:9867 \
  -v pinchtab-data:/data \
  --shm-size=2g \
  pinchtab
```

On first boot, the image creates `/data/.config/pinchtab/config.json` and generates a token if needed. When you use the managed-config path, the container binds to `0.0.0.0` at runtime via `PINCHTAB_BIND`, but the persisted config remains on the normal loopback default unless you override it yourself.

If you inspect the startup security summary from inside Docker, the loopback bind check will still report the effective runtime bind as non-loopback. That is expected: the process is listening on `0.0.0.0` inside the container so Docker port publishing can forward traffic to it.

This does not automatically mean the service is exposed beyond your machine. Host exposure still depends on how you publish the container port. For example:

- `-p 127.0.0.1:9867:9867` keeps the service reachable only from the host machine
- `-p 9867:9867` exposes it on the host's network interfaces

Treat the Docker runtime bind and the host-published address as separate layers. If you expose PinchTab beyond localhost, keep an auth token set and put it behind TLS or a trusted reverse proxy.

## Supplying Your Own `config.json`

If you want to manage the config file yourself, mount it and point `PINCHTAB_CONFIG` at it:

```text
docker-data/
└── config.json
```

Example `docker-data/config.json`:

```json
{
  "server": {
    "bind": "0.0.0.0",
    "port": "9867",
    "stateDir": "/data/state"
  },
  "profiles": {
    "baseDir": "/data/profiles",
    "defaultProfile": "default"
  },
  "instanceDefaults": {
    "mode": "headless",
    "noRestore": true
  }
}
```

Run with an explicit config file:

```bash
docker run -d \
  --name pinchtab \
  -p 127.0.0.1:9867:9867 \
  -e PINCHTAB_CONFIG=/config/config.json \
  -v "$PWD/docker-data:/data" \
  -v "$PWD/docker-data/config.json:/config/config.json:ro" \
  --shm-size=2g \
  pinchtab
```

Check it:

```bash
curl http://localhost:9867/health
curl http://localhost:9867/instances
```

## What To Persist

If you want data to survive container restarts, persist:

- the managed config directory or your mounted config file
- the profile directory
- the state directory

Without a mounted volume, profiles and saved session state are ephemeral.

## Runtime Configuration

For current runtime overrides, rely on:

- `PINCHTAB_CONFIG`
- `PINCHTAB_BIND`
- `PINCHTAB_PORT`
- `PINCHTAB_TOKEN`
- `CHROME_BIN`

Everything else should go in `config.json`.

In the bundled image, you usually do not need to set `PINCHTAB_BIND`, `PINCHTAB_PORT`, or `CHROME_BIN` manually unless you are overriding the generated config or replacing the browser binary. The managed-config entrypoint supplies `PINCHTAB_BIND=0.0.0.0` at runtime so Docker port publishing works without broadening the persisted config.

## Compose

The repository includes a `docker-compose.yml` that follows the managed-config pattern:

1. mount a persistent `/data` volume
2. let the entrypoint create and maintain `/data/.config/pinchtab/config.json`
3. optionally pass `PINCHTAB_TOKEN`

If you prefer a fully user-managed config file, mount it separately and set `PINCHTAB_CONFIG`.

If you expose PinchTab beyond localhost, set an auth token and put it behind TLS or a trusted reverse proxy.

## Resource Notes

Chrome in containers usually needs:

- larger shared memory, such as `--shm-size=2g`
- enough RAM for your tab count and workload

For heavier scraping or testing workloads, also consider:

- lowering `instanceDefaults.maxTabs`
- setting block options like `blockImages` in config
- running multiple smaller containers instead of one oversized browser

## Multi-Instance In Containers

You can run orchestrator mode inside one container and start managed instances from the API, but many teams prefer one browser service per container because:

- lifecycle is simpler
- container-level resource limits are clearer
- restart behavior is easier to reason about

Choose based on whether you want container-level isolation or PinchTab-managed multi-instance orchestration.
