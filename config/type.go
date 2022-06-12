package config

type (
	MNConfig struct {
		Server   MNConfigServer   `validate:"required"`
		Database MNConfigDatabase `validate:"required"`
	}

	MNConfigServer struct {
		Host    string
		Port    uint16
		Prefork bool
	}

	MNConfigDatabase struct {
		Host     string
		Port     uint16
		User     string
		Password string
		Name     string
		SSLMode  string
	}
)
