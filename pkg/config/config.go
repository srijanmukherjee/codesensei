package config

import (
	"github.com/spf13/viper"
)

// Defaults
const (
	DefaultMaxMemoryLimit                          = 1024000
	DefaultMemoryLimit                             = 256000
	DefaultMaxCpuTimeLimit                         = 15
	DefaultCpuTimeLimit                            = 5
	DefaultMaxCpuExtraTime                         = 5
	DefaultCpuExtraTime                            = 1
	DefaultMaxWallTimeLimit                        = 20
	DefaultWallTimeLimit                           = 10
	DefaultMaxStackLimit                           = 128000
	DefaultStackLimit                              = 64000
	DefaultMaxMaxProcessesAndOrThreads             = 120
	DefaultMaxProcessesAndOrThreads                = 60
	DefaultEnablePerProcessAndThreadTimeLimit      = false
	DefaultAllowEnablePerProcessAndThreadTimeLimit = false
	DefaultMaxIterations                           = 20
	DefaultIterations                              = 1
	DefaultRedirectStderrToStdout                  = false
	DefaultAllowNetwork                            = false
	DefaultServerPort                              = 3000
)

type AppConfig struct {
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     int    `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDb       string `mapstructure:"POSTGRES_DB"`

	MaxMemoryLimit                          int  `mapstructure:"MAX_MEMORY_LIMIT"`
	MemoryLimit                             int  `mapstructure:"MEMORY_LIMIT"`
	MaxCpuTimeLimit                         int  `mapstructure:"MAX_CPU_TIME_LIMIT"`
	CpuTimeLimit                            int  `mapstructure:"CPU_TIME_LIMIT"`
	MaxCpuExtraTime                         int  `mapstructure:"MAX_CPU_EXTRA_TIME"`
	CpuExtraTime                            int  `mapstructure:"CPU_EXTRA_TIME"`
	MaxWallTimeLimit                        int  `mapstructure:"MAX_WALL_TIME_LIMIT"`
	WallTimeLimit                           int  `mapstructure:"WALL_TIME_LIMIT"`
	MaxStackLimit                           int  `mapstructure:"MAX_STACK_LIMIT"`
	StackLimit                              int  `mapstructure:"STACK_LIMIT"`
	MaxMaxProcessesAndOrThreads             int  `mapstructure:"MAX_MAX_PROCESSES_AND_OR_THREADS"`
	MaxProcessesAndOrThreads                int  `mapstructure:"MAX_PROCESSES_AND_OR_THREADS"`
	EnablePerProcessAndThreadTimeLimit      bool `mapstructure:"ENABLE_PER_PROCESS_AND_THREAD_MEMORY_LIMIT"`
	AllowEnablePerProcessAndThreadTimeLimit bool `mapstructure:"ALLOW_ENABLE_PER_PROCESS_AND_THREAD_MEMORY_LIMIT"`
	MaxIterations                           int  `mapstructure:"MAX_ITERATIONS"`
	Iterations                              int  `mapstructure:"ITERATIONS"`
	RedirectStderrToStdout                  bool `mapstructure:"REDIRECT_STDERR_TO_STDOUT"`
	AllowNetwork                            bool `mapstructure:"ALLOW_NETWORK"`
	ServerPort                              int  `mapstructure:"SERVER_PORT"`
}

func LoadConfig(path string) (config AppConfig) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(path + "/codesensei.conf")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	config = AppConfig{
		MaxMemoryLimit:                          DefaultMaxMemoryLimit,
		MemoryLimit:                             DefaultMemoryLimit,
		MaxCpuTimeLimit:                         DefaultMaxCpuTimeLimit,
		CpuTimeLimit:                            DefaultCpuTimeLimit,
		MaxCpuExtraTime:                         DefaultMaxCpuExtraTime,
		CpuExtraTime:                            DefaultCpuExtraTime,
		MaxWallTimeLimit:                        DefaultMaxWallTimeLimit,
		WallTimeLimit:                           DefaultWallTimeLimit,
		MaxStackLimit:                           DefaultMaxStackLimit,
		StackLimit:                              DefaultStackLimit,
		MaxMaxProcessesAndOrThreads:             DefaultMaxMaxProcessesAndOrThreads,
		MaxProcessesAndOrThreads:                DefaultMaxProcessesAndOrThreads,
		EnablePerProcessAndThreadTimeLimit:      DefaultEnablePerProcessAndThreadTimeLimit,
		AllowEnablePerProcessAndThreadTimeLimit: DefaultAllowEnablePerProcessAndThreadTimeLimit,
		MaxIterations:                           DefaultMaxIterations,
		Iterations:                              DefaultIterations,
		RedirectStderrToStdout:                  DefaultRedirectStderrToStdout,
		AllowNetwork:                            DefaultAllowNetwork,
		ServerPort:                              DefaultServerPort,
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return
}

var Config = LoadConfig("/home/srijan/code/codesensei")
