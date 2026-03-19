package config

type Manifest struct {
	FormatVersion int      `json:"formatVersion"`
	Tool          string   `json:"tool"`
	EnvName       string   `json:"envName"`
	EnvType       string   `json:"envType"`
	Includes      []string `json:"includes"`
}
