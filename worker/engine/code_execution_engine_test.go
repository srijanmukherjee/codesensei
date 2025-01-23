package engine

import (
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/srijanmukherjee/codesensei/shared/config"
	"github.com/stretchr/testify/assert"
)

var (
	PROJECT_ROOT = os.Getenv("PROJECT_ROOT")
	CONFIG_FILE  = path.Join(PROJECT_ROOT, "config/codesensei.dev.yaml")

	Config = config.LoadConfigFile(CONFIG_FILE)
)

func TestCodeExecutionInitAndCleanup(t *testing.T) {
	context := CodeExecutionContext{
		Config:         Config,
		Id:             rand.Uint32(),
		SourceFilename: "main.c",
	}

	assert.Equal(t, nil, context.Init(), "Init should return nil")

	assertFileExists(t, context, "./")
	assertFileExists(t, context, "box/main.c")
	assertFileExists(t, context, "metadata.txt")
	assertFileExists(t, context, "stdin.txt")
	assertFileExists(t, context, "box/run.sh")

	assert.Equal(t, nil, context.Cleanup(), "Cleanup should return nil")

	assertFileNotExists(t, context, "./")
}

func TestSuccessfulCompilation(t *testing.T) {
	context := CodeExecutionContext{
		Config:          Config,
		Id:              rand.Uint32(),
		CompileCommand:  "gcc %s -o main main.c",
		CompilerOptions: "-Wall -O2",
		SourceFilename:  "main.c",
		SourceContent:   "#include <stdio.h>\nint main() { printf(\"Hello, World!\\n\"); return 0; }",
	}

	assert.NoError(t, context.Init(), "Init should return nil")

	compileResult, err := context.Compile()
	assert.NoError(t, err, "Compile should not return an error")
	assert.Equal(t, true, compileResult.Success, "Compile should be successful")
	assert.Equal(t, 0, compileResult.ExitCode, "Exit code should be 0")

	assert.NoError(t, context.Cleanup(), "Cleanup should return nil")
}

// TODO: compile failure test
// TODO: inavlid compile command test

func TestSuccessfulRun(t *testing.T) {
	context := CodeExecutionContext{
		Config:                               Config,
		Id:                                   rand.Uint32(),
		SourceFilename:                       "main.c",
		RunCommand:                           "./main",
		SourceContent:                        "#include <stdio.h>\n#include <unistd.h>\nint main() { int n; int*p=NULL; scanf(\"%d\", &n); printf(\"Hello, World! %d\\n\", *p); return 0; }",
		CompileCommand:                       "gcc %s -o main main.c",
		CompilerOptions:                      "-Wall -O2",
		EnablePerProcessAndThreadMemoryLimit: true,
		EnablePerProcessAndThreadTimeLimit:   false,
		MemoryLimit:                          512000,
		CpuTimeLimit:                         1,
		CpuExtraTime:                         0,
		WallTimeLimit:                        2,
		StackLimit:                           512000,
		EnableNetwork:                        false,
	}

	assert.NoError(t, context.Init(), "Init should return nil")

	compileResult, err := context.Compile()
	assert.NoError(t, err, "Compile should not return an error")
	assert.Equal(t, true, compileResult.Success, "Compile should be successful")
	assert.Equal(t, 0, compileResult.ExitCode, "Exit code should be 0")

	_, err = context.Run("1")
	assert.NoError(t, err, "Run should not return an error")
	// assert.Equal(t, true, runResult.Success, "Run should be successful")
	// assert.Equal(t, 0, runResult.ExitCode, "Exit code should be 0")
	// assert.Equal(t, "Hello, World!\n", runResult.Output, "Output should be 'Hello, World!\n'")

	assert.NoError(t, context.Cleanup(), "Cleanup should return nil")
}

func assertFileExists(t *testing.T, context CodeExecutionContext, filename string) {
	filepath := path.Join(context.workdir, filename)
	_, err := os.Stat(filepath)
	assert.NoError(t, err, "File %s should exist", filepath)
}

func assertFileNotExists(t *testing.T, context CodeExecutionContext, filename string) {
	filepath := path.Join(context.workdir, filename)
	_, err := os.Stat(filepath)
	assert.Error(t, err, "File %s should not exist", filepath)
}
