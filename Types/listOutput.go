package LocalTypes

type ListInfo struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
	Writes int `json:"writes"`
	Reads int `json:"reads"`
	Age int64`json:"age"`
}
