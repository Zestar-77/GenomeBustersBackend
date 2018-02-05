package configurationHandler

/*

import (
	"time"

	"github.com/gizak/termui"
)

var timeWhenStarted = time.Now()
var output *termui.Par
var writtenToOutput chan []byte

// ToOutput is a struct that writes to the output pane
type ToOutput struct{}

func (to *ToOutput) Write(b []byte) (n int, err error) {
	writtenToOutput <- b
	return len(b), nil
}

// GetOutput Retrives an writer for writing to the log pane
func GetOutput() *ToOutput {
	return &ToOutput{}
}

func visualOut() {
	defer termui.Close()
	uptime := termui.NewPar("0")
	uptime.BorderLabel = "Uptime"

	output = termui.NewPar("")
	output.BorderLabel = "Log"
	output.Text = "Hi!"

	termui.Handle("/timer/1s", func(event termui.Event) {
		uptime.Text = time.Since(timeWhenStarted).String()
	})

	termui.Handle("/sys/kbd/q", func(event termui.Event) {
		termui.StopLoop()
	})

	termui.Body.AddRows(
		termui.NewCol(6, 0, output),
		termui.NewCol(6, 0,
			termui.NewRow(
				termui.NewCol(2, 0, uptime))))

	go func() {
		for text := range <-writtenToOutput {
			output.Text += string(text)
		}
	}()

	termui.Loop()
	close(writtenToOutput)
	// TODO Handle exit condition
}
*/
