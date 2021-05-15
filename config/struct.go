package config

type (
	MNConfig struct {
		Database MNCDatabase
	}

	MNCDatabase struct {
		Host       string
		Port       uint16
		User       string
		Password   string
		Name       string
		Sslmode    string
		MaxConnect uint
	}
)
