package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	TextBold     = "\033[1m"
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorCyan    = "\033[36m"
	ColorMagenta = "\033[35m"
	ColorGray    = "\033[37m"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		fmt.Println("\n" + ColorGray + "\\> cancel" + ColorReset)
		os.Exit(0)
	}()

	runnerFile := filepath.Join(getCwd(), ".runner")

	if _, err := os.Stat(runnerFile); os.IsNotExist(err) {
		fmt.Println(TextBold + ColorRed + "Error: " + ColorReset + "No .runner file found in the directory")
		return
	}

	lines, err := readLines(runnerFile)
	if err != nil {
		fmt.Printf(TextBold+ColorRed+"Error: "+ColorReset+"reading .runner file: %v\n", err)
		return
	}

	args := os.Args
	verbose := containsArg("--verbose", args)

	if verbose {
		fmt.Println(TextBold + ColorCyan + "Info: " + ColorReset + "Verbose mode enabled")
		fmt.Printf(TextBold+ColorBlue+"Info: "+ColorReset+"Loaded %d lines from .runner\n", len(lines))
	}

	if len(args) > 1 && !strings.HasPrefix(args[1], "--") {
		handleCommand(args[1:], lines, verbose)
	} else {
		if containsArg("--list", args) {
			displayCommands(lines)
		} else {
			mainCommand := findMainCommand(lines)
			if mainCommand != "" {
				if verbose {
					fmt.Println(TextBold + ColorGreen + "Info: " + ColorReset + "Running `main` command...")
				}
				handleCommand([]string{"main"}, lines, verbose)
			} else {
				interactiveMode(lines)
			}
		}
	}
}

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf(TextBold+ColorRed+"Error: "+ColorReset+"getting current directory: %v\n", err)
		os.Exit(1)
	}
	return dir
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func handleCommand(args []string, lines []string, verbose bool) {
	command := args[0]
	if command == "" {
		fmt.Println(TextBold + ColorRed + "Error: " + ColorReset + "No command specified")
		displayCommands(lines)
		return
	}

	var commandStr string
	for _, line := range lines {
		if strings.HasPrefix(line, command+":") {
			commandStr = strings.TrimSpace(line[len(command)+1:])
			if len(args) > 1 {
				commandStr += " " + strings.Join(args[1:], " ")
			}
			break
		}
	}

	if commandStr == "" {
		fmt.Printf(TextBold+ColorRed+"Error: "+ColorReset+"No `%s` command found\n", command)
		displayCommands(lines)
		return
	}

	if verbose {
		fmt.Printf(TextBold+ColorGreen+"Info: "+ColorReset+"Executing command: %s\n", commandStr)
	}

	if err := runCommand(commandStr); err != nil {
		fmt.Printf(TextBold+ColorRed+"Error: "+ColorReset+"executing command: %v\n", err)
	}
}

func findMainCommand(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "main:") {
			return strings.TrimSpace(line[5:])
		}
	}
	return ""
}

func displayCommands(lines []string) {
	fmt.Println(TextBold + ColorMagenta + "Commands: " + ColorReset + "Available commands")
	printed := []string{}
	for _, line := range lines {
		if idx := strings.Index(line, ":"); idx != -1 {
			command := strings.TrimSpace(line[:idx])
			if containsArg(command, printed) || strings.HasPrefix(command, "#") {
				continue
			}
			printed = append(printed, command)
			fmt.Printf(ColorBlue+"- %s"+ColorReset+"\n", command)
		}
	}
}

func interactiveMode(lines []string) {
	fmt.Println(TextBold + ColorCyan + "Info: " + ColorReset + "Interactive mode: Select a command to run")
	displayCommands(lines)

	fmt.Print(TextBold + ColorGreen + "Prompt: " + ColorReset + "Enter your choice: ")
	var choice string
	fmt.Scanln(&choice)

	handleCommand([]string{choice}, lines, false)
}

func runCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func containsArg(arg string, args []string) bool {
	for _, a := range args {
		if a == arg {
			return true
		}
	}
	return false
}
