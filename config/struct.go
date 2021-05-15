package config

type (
	MNConfig struct {
		Database MNCDatabase
		Server   MNCServer
	}

	MNCDatabase struct {
		Host       string
		Port       uint16
		User       string
		Password   string
		Name       string
		Sslmode    string
		MaxConnect uint `toml:"max_connect"`
	}

	MNCServer struct {
		DisableAccountRegistrationLimit bool `toml:"disable_account_registration_limit"`
		AccountRegistrationLimitCount   uint `toml:"account_registration_limit_count"`
	}
)
