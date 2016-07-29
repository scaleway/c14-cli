package pgbar

import "gopkg.in/cheggaaa/pb.v1"

func NewProgressBar(title string) (bar *pb.ProgressBar) {
	bar = pb.New(100).SetWidth(80).SetMaxWidth(80).Format("[=> ]").Prefix(title)
	bar.ShowFinalTime = false
	bar.ShowTimeLeft = false
	bar.ShowCounters = false
	return
}
