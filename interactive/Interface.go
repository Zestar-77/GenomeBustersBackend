/*Package interactive Used for at runtime user control of the server.

This is currently the only supported means of adding to the gene database,
by telling the server to load in and parse a genebank file
*/
package interactive

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	cnf "GenomeBustersBackend/configurationHandler"
	"GenomeBustersBackend/genedatabase"
	"GenomeBustersBackend/global"

	"github.com/c-bata/go-prompt"
)

// printLn prints a line to the logger and to the console.
func printLn(s ...string) {
	st := strings.Join(s, " ")
	fmt.Println(st)
	global.Log.Println(st)
}

// printF does a formated print to both the console and to the logger
func printF(format string, i ...interface{}) {
	fmt.Printf(format, i...)
	global.Log.Printf(format, i...)
}

// addCommand is used for adding genes to the database
// currently, only handles genbank files.
// Will provide the use with file completions
func addCommand(args []string) {
	switch args[0] {
	case "gb":
		path := ""
		if len(args) >= 2 {
			path = args[1]
		} else {
			for path == "" {
				path = prompt.Input("[File Path]> ", getFileCompletions)
				global.Log.Println("[File Path]", path)
			}
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			printLn("Specified file does not exist!")
			return
		}

		genedatabase.AddGenBank(path)

	default:
		printF("Unrecognized sub command of add: %s\n", args[0])
	}

}

// RunTui runs a command prompt for user interaction
// Takes in a channel reading in os.Signal for handling interrupts.
//
// warning, this will panic if the server is not run in a command line interface.
func RunTui(interupt chan os.Signal) error {
	for {
		t := prompt.Input("[]> ", getCompletions)
		global.Log.Println("[]> ", t)
		split := strings.Fields(t)
		if len(split) == 0 {
			continue
		}
		switch split[0] {
		case "exit":
			return nil
		case "add": // Pull Genes out of specified
			if len(split) <= 1 {
				printF("add requires a subcommand\n")
			} else {
				addCommand(split[1:])
			}
			break
		case "dump":
			genedatabase.Dump()
		case "rui": // Rebuild front end
			npmBuild := exec.Command("npm", "run-script", "build")
			npmBuild.Dir = cnf.GetConfig().GetString("serverRoot")
			npmBuild.Stderr = os.Stderr
			if err := npmBuild.Run(); err != nil {
				printF("Unable to rebuild frontend\n%e\n", err)
			} else {
				printLn("Frontend rebuilt")
			}
			break
		case "echo":
			if len(split) > 1 {
				printLn(split[1:]...)
			}
		default:
			printF("Unrecognized command: %s\n", split[0])
		}
	}
}

// getCompletions holds a list of possible commands at the base level of the prompt.
// This is used for tab completions
func getCompletions(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "End this instance of the server"},
		{Text: "add gb", Description: "Add genes from a genbank file to the local gene database for lookups"},
		{Text: "rui", Description: "Rebuild frontend"},
		{Text: "dump", Description: "Print out the entirety of the database to the console."},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
