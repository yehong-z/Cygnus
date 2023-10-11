package message

type CountMessage struct {
	ObjectId int64 `json:"objectId"`
	Like     int64 `json:"like"`
	Dislike  int64 `json:"dislike"`
}
