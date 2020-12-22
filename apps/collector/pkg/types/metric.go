package types

type Metric struct {
	Source   string `json:"source"`
	Group    string `json:"group"`
	Instance string `json:"instance"`
	Clients  int    `json:"clients"`
	Messages int64  `json:"messages"`
	Volume   int64  `json:"volume"`
	Errors   int64  `json:"errors"`
}
