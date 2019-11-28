package config

type Config struct {
	Targets []Target `json:"targets"`
}

type Target struct {
	Proxy      string `json:"proxy"`
	Host       string `json:"host"`
	Https      bool   `json:"https"`
	ForceHttps bool   `json:"forceHttps"`
	Default    bool   `json:"default"`
}
