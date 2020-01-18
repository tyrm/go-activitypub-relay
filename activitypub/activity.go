package activitypub

type Activity struct {
	Context interface{} `json:"@context,omitempty"`
	ID      string      `json:"id,omitempty"`
	Actor   string      `json:"actor,omitempty"`
	Type    string      `json:"type,omitempty"`
	Object  interface{} `json:"object,omitempty"`
	To      []string    `json:"to,omitempty"`
	Cc      []string    `json:"cc,omitempty"`
}