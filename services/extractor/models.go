package extractor

type Project struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}

type User struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
