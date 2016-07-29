package pgbar

import "gopkg.in/cheggaaa/pb.v1"

// NewProgressBar returns a pb.ProgressBar with the good settings
// TODO make the progressbar more generic
func NewProgressBar(title string) (bar *pb.ProgressBar) {
	bar = pb.New(100).SetWidth(80).SetMaxWidth(80).Format("[=> ]").Prefix(title)
	bar.ShowFinalTime = false
	bar.ShowTimeLeft = false
	bar.ShowCounters = false
	return
}
