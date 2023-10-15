package taskrunner

import (
	"fmt"
	"math/rand"
	"os/exec"
	"path"
	"testing"

	"github.com/srijanmukherjee/codesensei/models"
)

// Test compilation
func TestCXXSuccessfulCompilation(t *testing.T) {
	ctx := TaskContext{
		Submission: models.Submission{
			Id: 1,
			SourceCode: "#include <iostream> \n" +
				"int main() { std::cout << \"Hello, World\"; return 0; }",
			Language: models.Language{
				CompileCommand: "/usr/bin/g++ %s main.cpp",
				SourceFile:     "main.cpp",
			},
			CompilerOptions: "",
		},
		Cgroups: "--cg",
	}

	err := ctx.initialize()
	if err != nil {
		t.Logf("initialization failed: %v\n", err)
		_, err = exec.Command("isolate", "--cg", "-b", ctx.BoxId, "--cleanup").Output()
		if err != nil {
			t.Log("failed to cleanup isolate")
		}
		t.FailNow()
	}

	err = ctx.compile()
	if err != nil {
		t.Log("compilation failed")
		t.Log(err.Error())
		t.Log(ctx.Submission.CompileOutput)
		t.Fail()
	} else {
		t.Log("Compilation successful")
	}

	_, err = exec.Command("isolate", "--cg", "-b", ctx.BoxId, "--cleanup").Output()
	if err != nil {
		t.Log("failed to cleanup isolate")
	}
}

// Test Helper Functions
func TestCreateWorkDir(t *testing.T) {
	ctx := TaskContext{
		Submission: models.Submission{},
		BoxId:      "1",
		Cgroups:    "--cg",
	}

	err := ctx.createWorkDirectory()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	} else {
		t.Logf("Created working directory %v (box id: %v)", ctx.Workdir, ctx.BoxId)
	}

	_, err = exec.Command("isolate", "--cg", "-b", ctx.BoxId, "--cleanup").Output()
	if err != nil {
		t.Fail()
	}
}

func TestInitializeAndRemoveFile(t *testing.T) {

	filePath := path.Join(t.TempDir(), fmt.Sprintf("temp-%v", rand.Int()))
	err := initializeFile(filePath)
	if err != nil {
		t.FailNow()
	}

	err = removeFile(filePath)
	if err != nil {
		t.FailNow()
	}
}
