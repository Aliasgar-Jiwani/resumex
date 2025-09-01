// File: README.md

# resumex - Universal Resumable Command Wrapper

A CLI tool that allows you to run any command with state management and resume capability if interrupted.

## Features

- **Universal command wrapping**: Run any command with resumable session tracking
- **Smart resume logic**: Built-in plugins for common tools (wget, rsync, tar, curl)
- **Session management**: List, resume, and delete saved sessions
- **Live output forwarding**: See command output in real-time while logging to file
- **Signal handling**: Gracefully handle interruptions (Ctrl+C, SIGTERM)

## Installation

```bash
# Clone and build
git clone <your-repo-url>
cd resumex
go build -o resumex
sudo mv resumex /usr/local/bin/
```

## Usage

### Run a command with session tracking

```bash
resumex run wget https://example.com/bigfile.iso
resumex run rsync -avz source/ destination/
resumex run tar -czf backup.tar.gz /important/data/
```

### List all sessions

```bash
resumex list
```

### Resume an interrupted session

```bash
resumex resume <session-id>
```

### Delete a session

```bash
resumex delete <session-id>
```

## How it Works

1. **Session Creation**: When you run `resumex run <command>`, it creates a unique session with metadata
2. **Command Execution**: The original command is executed with stdout/stderr captured to a log file
3. **Signal Handling**: If interrupted (Ctrl+C), the session is marked as incomplete
4. **Smart Resume**: When resuming, resumex uses command-specific plugins to add appropriate resume flags
5. **Logging**: All output is saved to `~/.resumex/logs/<session-id>.log`

## Supported Commands (Plugins)

- **wget**: Automatically adds `-c` (continue) flag
- **rsync**: Automatically adds `--partial` flag
- **tar**: Adds appropriate flags based on operation type
- **curl**: Automatically adds `-C -` (continue from offset)

For unsupported commands, resumex will simply re-run the original command.

## File Structure

```
~/.resumex/
├── sessions/           # Session metadata (JSON files)
└── logs/              # Command output logs
```

## Examples

```bash
# Start a large download
$ resumex run wget https://releases.ubuntu.com/22.04/ubuntu-22.04.3-desktop-amd64.iso

# If interrupted, resume it
$ resumex list
SESSION ID   STATUS      COMMAND                    START TIME       WORKING DIR
a1b2c3d4     interrupted wget https://releases.u...  2023-12-01 14:30 /home/user/downloads

$ resumex resume a1b2c3d4
# This will run: wget https://releases.ubuntu.com/22.04/ubuntu-22.04.3-desktop-amd64.iso -c
```

## Development

To extend resumex with custom plugins:

```go
import "github.com/Aliasgar-Jiwani/resumex/pkg/plugins"

// Register a custom plugin
plugins.RegisterPlugin("mycommand", func(sess *session.Session) string {
    return sess.Command + " --resume-flag"
})
```

## License

MIT License
