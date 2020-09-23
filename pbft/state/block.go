package state

type Block struct {
	Height int    `json:"height"`
	Data   string `json:"data"`
}
