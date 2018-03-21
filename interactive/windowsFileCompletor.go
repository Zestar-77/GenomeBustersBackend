// +build windows

package interactive

import (
	"io/ioutil"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

func getFileCompletions(d prompt.Document) []prompt.Suggest {
	path := d.Text
	if len(path) < 1 {
		return []prompt.Suggest{}
	}

	if !strings.HasPrefix(path, "C:\\") {
		if !strings.HasPrefix(path, ".\\") {
			path = ".\\" + path
		}
	}

	lastDir := path[:strings.LastIndex(path, "\\")+1]
	current := path[strings.LastIndex(path, "\\")+1:]
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

	return suggestions
}
