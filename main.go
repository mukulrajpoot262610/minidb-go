package main

import (
	"bufio"
	"fmt"
	"minidb/core"
	"minidb/utils"
	"os"
	"strings"
)

func printWelcomeMessage() {
	fmt.Println("MiniDB v0.1")
	fmt.Println("Type 'help' for available commands.")
	fmt.Println("Type 'exit' to exit the program.")

}

func main() {
	utils.LoadSchemasFromFile()
	printWelcomeMessage()

	for {
		fmt.Print("db > ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		core.ParseAndExecute(command)
	}
}
