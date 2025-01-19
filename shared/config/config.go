package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ENVIRONMENT_DEVELOPMENT = "DEVELOPMENT"
	ENVIRONMENT_PRODUCTION  = "PRODUCTION"
)

type Config struct {
	Environment                               string  `yaml:"environment"`
	CpuTimeLimit                              float64 `yaml:"cpu_time_limit"`
	MaxCpuTimeLimit                           float64 `yaml:"max_cpu_time_limit"`
	CpuExtraTime                              float64 `yaml:"cpu_extra_time"`
	MaxCpuExtraTime                           float64 `yaml:"max_cpu_extra_time"`
	WallTimeLimit                             float64 `yaml:"wall_time_limit"`
	MaxWallTimeLimit                          float64 `yaml:"max_wall_time_limit"`
	MemoryLimit                               int     `yaml:"memory_limit"`
	MaxMemoryLimit                            int     `yaml:"max_memory_limit"`
	StackLimit                                int     `yaml:"stack_limit"`
	MaxStackLimit                             int     `yaml:"max_stack_limit"`
	MaxProcessesAndOrThreads                  int     `yaml:"max_processes_and_or_threads"`
	MaxMaxProcessesAndOrThreads               int     `yaml:"max_max_processes_and_or_threads"`
	EnablePerProcessAndThreadTimeLimit        bool    `yaml:"enable_per_process_and_thread_time_limit"`
	AllowEnablePerProcessAndThreadTimeLimit   bool    `yaml:"allow_enable_per_process_and_thread_time_limit"`
	EnablePerProcessAndThreadMemoryLimit      bool    `yaml:"enable_per_process_and_thread_memory_limit"`
	AllowEnablePerProcessAndThreadMemoryLimit bool    `yaml:"allow_enable_per_process_and_thread_memory_limit"`
	MaxFileSize                               int     `yaml:"max_file_size"`
	MaxMaxFileSize                            int     `yaml:"max_max_file_size"`
	RedirectStderrToStdout                    bool    `yaml:"redirect_stderr_to_stdout"`
	AllowEnableNetwork                        bool    `yaml:"allow_enable_network"`
	EnableNetwork                             bool    `yaml:"enable_network"`
}

type ConfigContext struct {
	Config Config
}

func LoadConfigFile(configFile string) Config {
	config := Config{
		Environment:                               ENVIRONMENT_DEVELOPMENT,
		CpuTimeLimit:                              5,
		MaxCpuTimeLimit:                           15,
		CpuExtraTime:                              1,
		MaxCpuExtraTime:                           5,
		WallTimeLimit:                             10,
		MaxWallTimeLimit:                          20,
		MemoryLimit:                               128000,
		MaxMemoryLimit:                            512000,
		StackLimit:                                64000,
		MaxStackLimit:                             128000,
		MaxProcessesAndOrThreads:                  60,
		MaxMaxProcessesAndOrThreads:               120,
		EnablePerProcessAndThreadTimeLimit:        false,
		AllowEnablePerProcessAndThreadTimeLimit:   true,
		EnablePerProcessAndThreadMemoryLimit:      false,
		AllowEnablePerProcessAndThreadMemoryLimit: true,
		MaxFileSize:                               1024,
		MaxMaxFileSize:                            4096,
		RedirectStderrToStdout:                    false,
		AllowEnableNetwork:                        true,
		EnableNetwork:                             false,
	}

	configFileContents, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	err = yaml.Unmarshal(configFileContents, &config)
	if err != nil {
		log.Fatalf("Failed to parse yaml configuration: %v", err)
	}

	if config.Environment != ENVIRONMENT_DEVELOPMENT && config.Environment != ENVIRONMENT_PRODUCTION {
		log.Fatalf("Unknown environment '%s'", config.Environment)
	}

	return config
}
