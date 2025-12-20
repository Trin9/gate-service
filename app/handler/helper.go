package handler

var vllmURL string

func SetVllmURL(url string) {
	vllmURL = url
}

func GetVllmURL() string {
	return vllmURL
}
