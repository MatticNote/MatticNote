package config

type (
	MNConfig struct {
		Server   MNConfigServer   `validate:"required"`
		Database MNConfigDatabase `validate:"required"`
		Redis    MNConfigRedis    `validate:"required"`
		System   MNConfigSystem   `validate:"required"`
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

	MNConfigRedis struct {
		Host     string
		Port     uint16
		User     string
		Password string
		Database int
	}

	MNConfigSystem struct {
		// RegistrationMode: 0 -> Registration disabled, 1 -> Invite only, 2 -> Open
		RegistrationMode uint `toml:"registration_mode" validate:"gte=0,lte=2"`
		// InvitePermission: 0 -> Administrator only, 1 -> Administrator and Moderator only, 2 -> Everyone
		InvitePermission uint `toml:"invite_permission" validate:"gte=0,lte=2"`
	}
)
