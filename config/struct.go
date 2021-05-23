package config

type (
	MNConfig struct {
		Database MNCDatabase `validate:"required"`
		Server   MNCServer   `validate:"required"`
		Redis    MNCRedis    `validate:"required"`
		Mail     MNCMail     `validate:"required"`
	}

	MNCDatabase struct {
		Host       string `validate:"required,hostname_rfc1123"`
		Port       uint16 `validate:"gte=0,lte=65535"`
		User       string `validate:"required"`
		Password   string
		Name       string `validate:"required"`
		Sslmode    string `validate:"required"`
		MaxConnect uint   `toml:"max_connect" validate:"required"`
	}

	MNCServer struct {
		ListenAddress                   string `toml:"listen_address" validate:"required,hostname_rfc1123"`
		ListenPort                      uint16 `toml:"listen_port" validate:"gte=0,lte=65535"`
		DisableAccountRegistrationLimit bool   `toml:"disable_account_registration_limit"`
		AccountRegistrationLimitCount   uint   `toml:"account_registration_limit_count"`
		CookieSecure                    bool   `toml:"cookie_secure"`
		Endpoint                        string `validate:"required"`
		OAuthSecretKey                  string `toml:"oauth_secret_key" validate:"required,len=32"`
	}

	MNCRedis struct {
		Address  string `validate:"required,hostname_rfc1123"`
		Port     uint16 `validate:"gte=0,lte=65535"`
		Username string
		Password string
		Database int
	}

	MNCMail struct {
		From     string `validate:"required"`
		Username string
		Password string
		SmtpHost string `validate:"required" toml:"smtp_host"`
		SmtpPort uint16 `validate:"gte=0,lte=65535" toml:"smtp_port"`
	}
)
