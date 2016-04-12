package cli

import (
	"io"
	"text/tabwriter"
)

// newTabWriter returns an initialized tab Writer writes tabbed text as
// column aligned text to the given io.Writer.
func newTabWriter(writer io.Writer) *tabwriter.Writer {
	tw := new(tabwriter.Writer)
	tw.Init(writer, 0, 8, 1, '\t', 0)
	return tw
}
