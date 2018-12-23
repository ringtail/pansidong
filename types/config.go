package types

type GlobalConfig struct {
	Port int
	Host string
}

type Config struct {
	GlobalConfig *GlobalConfig
	Backend      *BackendConfigUnknown
	Memory       *MemoryConfig
}
