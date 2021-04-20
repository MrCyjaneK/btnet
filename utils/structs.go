package utils

type Files struct {
	Path   string `json:"path"`
	Sha512 string `json:"sha512"`
}

type BTnetJson struct {
	Build string
	Files []Files
}
