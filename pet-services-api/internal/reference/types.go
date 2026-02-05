package reference

type Country struct {
	Name     string `json:"name"`
	DialCode string `json:"dial_code"`
	Code     string `json:"code"`
	Flag     string `json:"flag"`
}

type State struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type City struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	StateID int    `json:"state_id"`
}