package http

type Config struct {
	Targets struct {
		Users string `koanf:"users"`
		Books string `koanf:"books"`
	} `koanf:"targets"`
}
