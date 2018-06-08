package app

type Msg struct {
	Action string `json:"action"`
	Key    string `json:"key, omitempty"`
	Value  string `json:"value, omitempty"`
}
