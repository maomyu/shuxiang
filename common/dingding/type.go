package dingding

// 存放链接消息的结构体
type LinkMessageResult struct {
	Msgtype string     `json:"msgtype"`
	Link    LinkResult `json:"link"`
}
type LinkResult struct {
	MessageUrl string `json:"messageUrl"`
	PicUrl     string `json:"picUrl"`
	Title      string `json:"title"`
	Text       string `json:"text"`
}
