package main

import (
	"GNBS/compiler"
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gnbs",
		Short: "Compiler",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			vm := compiler.NewVM()
			if len(args) == 0 {
				repl(vm)
			} else if len(args) == 1 {
				runFile(args[0], vm)
			} else {
				os.Exit(64)
				return
			}

		},
	}
	return cmd
}

func repl(vm *compiler.VM) {
	reader := bufio.NewReader(os.Stdin)
	var buffer bytes.Buffer
	for {
		fmt.Print("> ")
		read, _ := reader.ReadString('\n')
		buffer.WriteString(read)
		vm.Interpret(buffer.Bytes())
	}
}

func runFile(path string, vm *compiler.VM) {
	fileBytes := readFile(path)
	result := vm.Interpret(fileBytes)

	if result == compiler.InterpretCompileError {
		os.Exit(65)
	}
	if result == compiler.InterpretRuntimeError {
		os.Exit(70)
	}
}

func readFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		handleError(err)
	}
	info, _ := file.Stat()
	fileSize := info.Size()

	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		handleError(err)
	}
	return buffer
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
