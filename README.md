# pikvm-automator

API processor for simple and short macros language to automate PiKVM keyboard and mouse emulation.

[![Hub](https://badgen.net/docker/pulls/sealbro/pikvm-automator?icon=docker&label=pikvm-automator)](https://hub.docker.com/r/sealbro/pikvm-automator/)

```mermaid
flowchart LR
  subgraph Requests
    direction RL
    grpc("`gRPC`")
    openapi("`OpenAPI v3`")
  end
  subgraph PiKVM
    direction TB
    subgraph Docker
        pikvm-automator("`pikvm-automator`")
    end
    subgraph pikvm-api
        direction RL
        api1("`wss (websocket)`")
        api2("`unix socket (kvmd.sock)s`")
    end
  end
  subgraph Device
    direction RL
    sub1("`Keyboard emulation`")
    sub2("`Mouse emulation`")
  end
  Requests --> PiKVM  --> Device
```

## How to use

Run the docker container on PiKVM device with default credentials:

```shell
docker run -it --rm \
 -v $(pwd):/commands.yml \
 -v /run/kvmd:/run/kvmd \
 -e PIKVM_ADDRESS=unix:/run/kvmd/kvmd.sock/ws?stream=0 \
 -e PIKVM_SKIP_TLS_VERIFY=true \
 sealbro/pikvm-automator
```

### API integration

- [gRPC](./proto/pikvm_automator.proto)
- [OpenAPI v3](./generated/openapiv3/openapi.yaml)
- [API examples](./endpoints.http)

### Environment variables

| Name                    | Description               | Default         |
|-------------------------|---------------------------|-----------------|
| `PIKVM_ADDRESS`         | PiKVM api address         | required        |
| `PIKVM_SKIP_TLS_VERIFY` | Skip TLS verify           | `false`         |
| `PIKVM_USERNAME`        | PiKVM username            | `admin`         |
| `PIKVM_PASSWORD`        | PiKVM password            | `admin`         |
| `COMMANDS_PATH`         | Path with commands file   | `/commands.yml` |
| `TEMPLATE_MAX_DEEP`     | Max recursive replacement | `10`            |
| `CALL_DEBOUNCE_SECONDS` | Next command cooldown     | `2`             |
| `GRPC_ADDRESS`          | gRPC endpoint             | `0.0.0.0:32023` |
| `GRPC_GATEWAY_ADDRESS`  | openapi endpoint          | `0.0.0.0:8032`  |

- `PIKVM_ADDRESS` is required
  - `unix:/run/kvmd/kvmd.sock/ws?stream=0` for unix socket on device
  - `wss://<IP or HOST>/api/ws?stream=0` for remote connection

## Macros language

Example of `commands.yml`:

```yaml
bios_enable_virtualization:
  id: bios_enable_virtualization
  description: Bios enable virtualization
  expression: '@850''190|@left|700ms|@350''650|@left|700ms|@850''420|@left|300ms|End|50ms|Enter'
bios_save_exit:
  id: bios_save_exit
  description: Bios save configuration and exit
  expression: F10|200ms|@850'770|@left
proxmox_mode:
  id: proxmox_mode
  description: Proxmox bios mode
  expression: 2s|F2|5s|%bios_enable_virtualization%|1s|%bios_save_exit%
```

- `id` - unique identifier
- `description` - description of command
- `expression` - sequence of key presses
  - `F1` - press F1 key ([more keys](./pkg/pikvm/keyboard/key.go))
  - `@850'190` - move mouse to x=850, y=190
  - `@left` - press left mouse button
  - `700ms` - delay 700ms
  - `%command_id%` - execute another command, allowed use recursion but not more than `TEMPLATE_MAX_DEEP` times