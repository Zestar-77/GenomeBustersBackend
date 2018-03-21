package interactive

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	cnf "GenomeBustersBackend/configurationHandler"

	"github.com/c-bata/go-prompt"
)

func printLn(s ...string) {
	st := strings.Join(s, " ")
	fmt.Println(st)
	log.Println(st)
}

func printF(format string, i ...interface{}) {
	fmt.Printf(format, i...)
	log.Printf(format, i...)
}

func addCommand(args []string) {
	switch args[0] {
	case "gb":
		path := ""
		if len(args) >= 2 {
			path = args[1]
		} else {
			for path == "" {
				path = prompt.Input("[File Path]> ", getFileCompletions)
			}
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("Specified file does not exist!")
			return
		}

		// TODO Load genbank file

	default:
		printF("Unrecognized sub command of add: %s\n", args[0])
	}

}

// RunTui runs a command prompt
func RunTui(interupt chan os.Signal) error {
	for {
		t := prompt.Input("[]> ", getCompletions)
		log.Println("[]> ", t)
		split := strings.Fields(t)
		if len(split) == 0 {
			continue
		}
		switch split[0] {
		case "exit":
			return nil
		case "add": // Pull Genes out of specified
			if len(split) <= 1 {
				printF("add requires a subcommand")
			} else {
				addCommand(split[1:])
			}
			break
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
		default:
			printF("Unrecognized command: %s\n", split[0])
		}
	}
}

func getCompletions(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "End this instance of the server"},
		{Text: "add gb", Description: "Add genes from a genbank file to the local gene database for lookups"},
		{Text: "rui", Description: "Rebuild frontend"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}