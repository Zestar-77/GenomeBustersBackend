// +build !windows

package interactive

import (
	"io/ioutil"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

// getFileCompletions gets a list of files/directories given the path currently entered by the user
func getFileCompletions(d prompt.Document) []prompt.Suggest {
	path := d.Text
	if len(path) < 1 {
		return []prompt.Suggest{}
	}

	if path[0] != '/' {
		if !(len(path) > 2 && path[0] == '.' && path[1] == '/') {
			path = "./" + path
		}
	}

	lastDir := path[:strings.LastIndex(path, "/")+1]
	current := path[strings.LastIndex(path, "/")+1:]
	files, err := ioutil.ReadDir(lastDir)
	if err != nil {
		return []prompt.Suggest{}
	}

	suggestions := make([]prompt.Suggest, 12)
	for _, value := range files {
		if strings.Contains(value.Name(), current) {
			suggestions = append(suggestions, prompt.Suggest{Text: value.Name(), Description: ""})
		}
	}

	return prompt.FilterContains(suggestions, d.GetWordBeforeCursor(), true)
}
