package message

type CountMessage struct {
	Like    int64 `json:"like"`
	Dislike int64 `json:"dislike"`
}
