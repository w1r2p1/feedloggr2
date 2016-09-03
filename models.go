package feedloggr2

import "time"

type Feed struct {
	Title string
	URL   string
	Items []*FeedItem
	Error error
}

type FeedSlice []*Feed

func (fs FeedSlice) Len() int {
	return len(fs)
}

func (fs FeedSlice) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs FeedSlice) Less(i, j int) bool {
	return fs[i].Title < fs[j].Title
}

type FeedItem struct {
	ID    int
	Title string `sql:"type:varchar(100)"`
	URL   string `sql:"type:varchar(255);unique"`
	Date  time.Time
	Feed  string `sql:"type:varchar(64)"`
}