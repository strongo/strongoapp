package strongo

// Environment defines environment
type Environment int8

const (
	// EnvUnknown is unknown environment
	EnvUnknown Environment = iota
	// EnvProduction is production environment
	EnvProduction
	// EnvStaging is staging environment
	EnvStaging
	// EnvDevTest is developers test environment
	EnvDevTest
	// EnvLocal is developer's local environment
	EnvLocal
)

// EnvironmentNames is mapping of environment codes to environment names
var EnvironmentNames = map[Environment]string{
	EnvUnknown:    "unknown",
	EnvProduction: "production",
	EnvStaging:    "staging",
	EnvDevTest:    "dev",
	EnvLocal:      "local",
}
