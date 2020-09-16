package pbft

type Block struct {
	Height int    `json:"height"`
	Data   string `json:"data"`
}
