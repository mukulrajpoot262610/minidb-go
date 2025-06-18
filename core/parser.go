package core

import (
	"fmt"
	"minidb/commands"
	"minidb/utils"
	"strings"
)

func removeSemicolon(command string) string {
	return strings.TrimSuffix(strings.TrimSpace(command), ";")
}

func ParseAndExecute(command string) {

	command = removeSemicolon(command)

	switch {
	case command == "":
		return

	case command == "help":
		utils.PrintHelp()
		return

	case command == "exit":
		fmt.Println("Exiting MiniDB.")
		utils.ExitProgram()

	case strings.HasPrefix(command, "insert into"):
		commands.Insert(command)

	case strings.HasPrefix(command, "select"):
		commands.Select(command)

	case strings.HasPrefix(command, "delete from"):
		commands.Delete(command)

	case strings.HasPrefix(command, "create table"):
		commands.CreateTable(command)

	case strings.HasPrefix(command, "show tables"):
		commands.ShowTables()

	case strings.HasPrefix(command, "drop table"):
		commands.DropTable(command)

	default:
		fmt.Println("Unrecognized command:", command)
	}
}

