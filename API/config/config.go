package config

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	SeedUser SeedUserConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port   int
	Secret string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	Server     string
	DBUser     string
	DBPassword string
	DBName     string
}

// SeedUserConfigurations exported
type SeedUserConfigurations struct {
	Email           string
	InitialPassword string
}
