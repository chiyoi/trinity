package eroira

type Image struct {
	Title  string `json:"title"`
	Pid    int    `json:"pid"`
	Author string `json:"author"`
	Uid    int    `json:"uid"`
	Urls   struct {
		Original string `json:"original"`
	} `json:"urls"`
}
type LoliconResponse struct {
	Data []Image `json:"data"`
}
