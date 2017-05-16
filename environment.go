package strongo

type Environment int8

const (
	EnvUnknown    Environment = iota
	EnvProduction
	EnvStaging
	EnvDevTest
	EnvLocal
)

var EnvironmentNames = map[Environment]string {
	EnvUnknown: "unknown",
	EnvProduction: "production",
	EnvStaging: "staging",
	EnvDevTest: "dev",
	EnvLocal: "local",
}