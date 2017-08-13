package output

import (
	"fmt"
	"github.com/gosuri/uilive"
	"strings"
	"time"
)

// ProgressBar - A structure that controls displaying a progress bar in the
// console.
type ProgressBar struct {
	writer    *uilive.Writer
	Max       int
	Current   int
	startTime time.Time
}

// NewProgress - Creates a new progress bar with a given maximum value.
func NewProgress(max int) *ProgressBar {
	writer := uilive.New()
	writer.Start()
	return &ProgressBar{
		writer:    writer,
		Max:       max,
		startTime: time.Now(),
	}
}

// Inc - Increments the current value by one.
func (bar *ProgressBar) Inc() {
	bar.Current++
	bar.write()
	if bar.Current == bar.Max {
		bar.Stop()
	}
}

// Stop - Stops the progress bar rendering. This is automatically called if the
// current value reaches the maximum.
func (bar *ProgressBar) Stop() {
	bar.writer.Stop()
}

// write - Writes the progress bar line to the console.
func (bar *ProgressBar) write() {
	elapsed := time.Since(bar.startTime)
	perc := (float64(bar.Current) / float64(bar.Max)) * float64(10)
	barStr := strings.Repeat("⚪", 10)
	barStr = strings.Replace(barStr, "⚪", "⚫", int(perc))
	line := fmt.Sprintf("[%d/%d] %s %s", bar.Current, bar.Max, barStr, elapsed)
	fmt.Fprintln(bar.writer, line)
}
