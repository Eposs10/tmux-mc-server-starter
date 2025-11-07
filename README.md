# Minecraft Server Tmux Manager

A lightweight Go script for managing Minecraft server sessions inside **tmux** — automatically restarts the server when it stops, logs exit codes, and allows graceful shutdowns.

---

## Features

- Automatically starts a Minecraft server in a **tmux** session  
- Prevents multiple sessions with the same name  
- Auto-restarts the server after crashes  
- Allows you to stop restarts by pressing **Enter** within a configurable wait time  
- Logs all exit codes to `exit_codes/server_exit_codes.log`  
- Configurable RAM usage, wait time, and JAR file  

---

## Requirements

- **Go 1.21+** (to build)
- **tmux** installed
- **Java** (for running the Minecraft server)

---

## Installation

1. Clone or download this repository:
   ```bash
   git clone git@github.com:Eposs10/tmux-mc-server-starter.git
   cd tmux-mc-server-starter
   ```

2. Build the executable:
   ```bash
   go build -o mc-tmux
   ```

3. Move it somewhere in your `$PATH` (optional):
   ```bash
   sudo mv mc-tmux /usr/local/bin/
   ```

---

## Usage

### Basic Command
```bash
mc-tmux <session_name> <path> [options]
```

### Example
```bash
mc-tmux survival /home/minecraft/server --max-ram 8G --min-ram 4G --wait-time 10
```

This starts a new tmux session named `survival`, running the server at `/home/minecraft/server` with 4–8 GB of RAM.  
If the server crashes or stops, it will wait **10 seconds** before restarting — unless you press **Enter**.

---

## Options

| Option | Type | Default | Description |
|--------|------|----------|-------------|
| `--jar` | `string` | `server.jar` | The Minecraft server JAR file to run |
| `--min-ram` | `string` | `2G` | Minimum RAM allocation |
| `--max-ram` | `string` | `6G` | Maximum RAM allocation |
| `--wait-time` | `int` | `5` | Seconds to wait before restarting |

---

## Examples

### Create and start a new session
```bash
mc-tmux creative /home/mc/creative
```

### Specify a custom JAR file
```bash
mc-tmux modded /home/mc/modded --jar forge-1.20.1.jar
```

### Attach to an existing session
If a session with the same name already exists, the script will automatically attach to it instead of creating a new one:
```
⚠️ Session 'survival' already exists — attaching instead.
```

---

## Logs

Exit codes are saved under:
```
exit_codes/server_exit_codes.log
```

Each entry includes a timestamp and the exit code.

---

## Stopping the Server

To stop the automatic restart loop:
1. Stop the server normally (e.g., `stop` command in Minecraft console).
2. When prompted:
   ```
   ----- Press enter to prevent the server from restarting in 5 seconds -----
   ```
   Press **Enter** before the timer runs out.
