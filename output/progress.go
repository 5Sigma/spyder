package output

import (
	"fmt"
	"github.com/gosuri/uilive"
	"strings"
	"time"
)

type ProgressBar struct {
	writer    *uilive.Writer
	Max       int
	Current   int
	startTime time.Time
}

func NewProgress(max int) *ProgressBar {
	writer := uilive.New()
	writer.Start()
	return &ProgressBar{
		writer:    writer,
		Max:       max,
		startTime: time.Now(),
	}
}

func (bar *ProgressBar) Inc() {
	bar.Current++
	bar.write()
	if bar.Current == bar.Max {
		bar.Stop()
	}
}

func (bar *ProgressBar) Stop() {
	bar.writer.Stop()
}

func (bar *ProgressBar) write() {
	elapsed := time.Since(bar.startTime)
	perc := (float64(bar.Current) / float64(bar.Max)) * float64(10)
	barStr := strings.Repeat("⚪", 10)
	barStr = strings.Replace(barStr, "⚪", "⚫", int(perc))
	line := fmt.Sprintf("[%d/%d] %s %s", bar.Current, bar.Max, barStr, elapsed)
	fmt.Fprintln(bar.writer, line)
}
