package handlers

type CreateBuff struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

type CreateStream struct {
	Name    string   `json:"name"`
	BuffIDs []uint64 `json:"buff_ids"`
}

type Buff struct {
	ID       uint64   `json:"id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

type Stream struct {
	ID      uint64   `json:"id"`
	Name    string   `json:"name"`
	BuffIDs []uint64 `json:"buff_ids"`
}
