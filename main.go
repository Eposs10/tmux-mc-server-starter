package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"time"
)

type Config struct {
	name     string
	path     string
	jar      string
	maxRAM   string
	minRAM   string
	waitTime int
}

func main() {
	if slices.Contains(os.Args, "-h") || slices.Contains(os.Args, "--help") {
		printHelp(0)
	}

	if len(os.Args) < 3 {
		printHelp(1)
	}

	config := Config{
		name:     os.Args[1],
		path:     os.Args[2],
		jar:      "server.jar",
		maxRAM:   "6G",
		minRAM:   "2G",
		waitTime: 5,
	}

	// Parse optional arguments
	for i := 3; i < len(os.Args); i += 2 {
		switch os.Args[i] {
		case "--jar":
			if i+1 < len(os.Args) {
				config.jar = os.Args[i+1]
			}
		case "--min-ram":
			if i+1 < len(os.Args) {
				config.minRAM = os.Args[i+1]
			}
		case "--max-ram":
			if i+1 < len(os.Args) {
				config.maxRAM = os.Args[i+1]
			}
		case "--wait-time":
			if i+1 < len(os.Args) {
				if t, err := strconv.Atoi(os.Args[i+1]); err == nil {
					config.waitTime = t
				}
			}
		}
	}

	check := exec.Command("tmux", "has-session", "-t", config.name)
	if err := check.Run(); err == nil {
		fmt.Printf("⚠️ Session '%s' already exists — attaching instead.\n", config.name)
		time.Sleep(500 * time.Millisecond)
		attachToSession(config.name)
		return
	}

	// --- Build the script to run inside tmux ---
	script := fmt.Sprintf(`
cd "%s"

JAR="%s"
MAXRAM="%s"
MINRAM="%s"
TIME=%d

while true; do
    java -Xmx$MAXRAM -Xms$MINRAM -jar $JAR nogui
    mkdir -p "exit_codes"
    touch "exit_codes/server_exit_codes.log"
    echo "[$(date +"%%d.%%m.%%Y %%T")] ExitCode: $?" >> exit_codes/server_exit_codes.log
    echo "----- Press enter to prevent the server from restarting in $TIME seconds -----"
    read -t $TIME input
    if [ $? == 0 ]; then
        break
    else
        echo "------------------- SERVER RESTARTS -------------------"
    fi
done
`, config.path, config.jar, config.maxRAM, config.minRAM, config.waitTime)

	// --- Run tmux new-session ---
	cmd := exec.Command("tmux", "new", "-d", "-s", config.name, "bash", "-s")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("❌ Failed to get stdin pipe:", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ Failed to start tmux session:", err)
		os.Exit(1)
	}

	// Write the script to tmux's stdin
	if _, err := stdin.Write([]byte(script)); err != nil {
		fmt.Println("❌ Failed to write script to tmux stdin:", err)
		os.Exit(1)
	}

	if err := stdin.Close(); err != nil {
		fmt.Println("❌ Error closing tmux stdin:", err)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("❌ Error waiting for tmux:", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Tmux session '%s' started in '%s'.\n", config.name, config.path)
	fmt.Println("Attaching now... (detach with Ctrl+B, then D)")
	time.Sleep(500 * time.Millisecond)

	attachToSession(config.name)
}

func attachToSession(name string) {
	attachCmd := exec.Command("tmux", "attach", "-t", name)
	attachCmd.Stdin = os.Stdin
	attachCmd.Stdout = os.Stdout
	attachCmd.Stderr = os.Stderr

	if err := attachCmd.Run(); err != nil {
		fmt.Println("❌ Failed to attach to tmux session:", err)
		os.Exit(1)
	}
}

func printHelp(exitCode int) {
	fmt.Println("Usage:")
	fmt.Println("  <session_name> <path> [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --jar string       JAR file name (default: server.jar)")
	fmt.Println("  --min-ram string   Minimum RAM (default: 2G)")
	fmt.Println("  --max-ram string   Maximum RAM (default: 6G)")
	fmt.Println("  --wait-time int    Seconds to wait before restart (default: 5)")

	os.Exit(exitCode)
}
