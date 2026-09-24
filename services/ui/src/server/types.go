package server

type manifest struct {
	Platform struct {
		Arch string `json:"architecture"`
		OS   string `json:"os"`
	} `json:"platform"`
}

type annotations struct {
	Arch string `json:"com.docker.official-images.bashbrew.arch"`
}

type tag struct {
	Name string
	Size string
	Arch []string
}

type layer struct {
	Size int `json:"size"`
}

// Response types

type repositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

type tagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type manifestResponse struct {
	Annotations annotations `json:"annotations"`
	Layers      []layer     `json:"layers"`
	Manifests   []manifest  `json:"manifests"`
}
