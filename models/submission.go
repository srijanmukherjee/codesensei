package models

type Submission struct {
	Id                                   int
	SourceCode                           string
	LanguageId                           int
	Stdins                               []string
	ExpectedOutputs                      []string
	Stdouts                              []string
	CreatedAt                            string
	FinishedAt                           string
	Time                                 *float32
	Memory                               *int
	Stderr                               *string
	Token                                string
	Iterations                           int
	CpuTimeLimit                         float32
	WallTimeLimit                        float32
	MemoryLimit                          int
	StackLimit                           int
	MaxProcessesAndOrThreads             int
	EnablePerProcessAndThreadTimeLimit   bool
	EnablePerProcessAndThreadMemoryLimit bool
	CompileOutput                        string
	ExitCode                             *int
	WallTime                             *float32
	CompilerOptions                      string
	CommandLineArguments                 string
	RedirectStderrToStdout               bool
	CallbackUrl                          string
	EnableNetwork                        bool
	StartedAt                            string
	QueuedAt                             string
	UpdatedAt                            string
	ExecutionHost                        string
	Language                             Language
}
