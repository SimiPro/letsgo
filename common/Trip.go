package common

type Trip struct {
	Id        string `json:"id"`
	User	  User   `json:"user"`
	Stops	  []Stop `json:"stops"`
}

type Stop struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Long float64 `json:"long"`
	Images []Image `json:"images"`
}

type Image struct {
	Id string `json:"id"`
	Path string `json:"path"`
	Story string `json:"story"`
}