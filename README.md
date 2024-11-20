
# go-runner

This script allows you to manage and execute commands specified in a `.runner` file within the current directory. This is a **go** implementation of my previous [runner](https://github.com/mhs003/runner/) script which I wrote in python.

## features
- **Run predefined commands** from a `.runner` file.
- **Interactive mode** for selecting commands.
- **Verbose mode** for detailed output.

## requirements
- [Go-lang](https://go.dev/) `(latest)`

## installation

1. Clone the script.
    ```bash
    git clone https://github.com/mhs003/go-runner
    cd go-runner/
    ```
2. Build the binary using:
   ```bash
   go build -o run
   ```

## usage

### Running the Script
To run the script, execute:
```bash
./run [options] [command] [arguments]
```

### options
- `--verbose`: Enables verbose mode.
- `--list`: Lists all available commands.

### commands
- Specify the name of the command as defined in the `.runner` file.

### interactive mode
If no command or option is specified, the script enters **Interactive Mode**, prompting you to select a command to execute.

## usage example

Create a `.runner` file in your working directory:
```plaintext
main: go run main.go
hello: python3 hello.py
goodbye: echo "Goodbye!"
main: node main.js
```

Example usage
```bash
$ ./run
# will execute `go run main.go` and will skip `node main.js`

$ ./run main
# will execute `go run main.go` and will skip `node main.js`

$ ./run goodbye
# will execute `echo "Goodbye!"`

$ ./run --list
Commands: Available commands
- main
- hello
- goodbye

```

## conclusion
I just started learning Go, so the code may look messy. Feel free to use and modify this script if you found it interesting.