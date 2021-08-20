package config

type (
	MNConfig struct {
		Database MNCDatabase `validate:"required"`
		Server   MNCServer   `validate:"required"`
		Redis    MNCRedis    `validate:"required"`
		Mail     MNCMail     `validate:"required"`
		Job      MNCJob      `validate:"required"`
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
		ListenAddress                 string `toml:"listen_address"`
		ListenPort                    uint16 `toml:"listen_port" validate:"gte=0,lte=65535"`
		Prefork                       bool   `toml:"prefork"`
		DisableAccountRateLimit       bool   `toml:"disable_account_rate_limit"`
		AccountRegistrationLimitCount uint   `toml:"account_registration_limit_count"`
		CookieSecure                  bool   `toml:"cookie_secure"`
		Endpoint                      string `validate:"required"`
		OAuthSecretKey                string `toml:"oauth_secret_key" validate:"required,len=32"`
		RecaptchaSiteKey              string `toml:"recaptcha_site_key"`
		RecaptchaSecretKey            string `toml:"recaptcha_secret_key"`
	}

	MNCRedis struct {
		Address  string `validate:"required,hostname_rfc1123"`
		Port     uint16 `validate:"gte=0,lte=65535"`
		Username string
		Password string
		Database int
	}

	MNCJob struct {
		MaxActive int `validate:"required" toml:"max_active"`
		MaxIdle   int `validate:"required" toml:"max_idle"`
	}

	MNCMail struct {
		From     string `validate:"required"`
		Username string
		Password string
		SmtpHost string `validate:"required" toml:"smtp_host"`
		SmtpPort uint16 `validate:"gte=0,lte=65535" toml:"smtp_port"`
	}
)
