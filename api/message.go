package api

type Message struct {
	Item    string        `json:"item"`
	Subject string        `json:"subject"`
	Mtype   int           `json:"type"` //#0发起 1反馈
	Status  bool          `json:"status"`
	Id      int64         `json:"id"`
	Data    []interface{} `json:"data"`
}
