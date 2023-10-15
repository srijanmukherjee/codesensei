package taskrunner

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/srijanmukherjee/codesensei/models"
	"github.com/srijanmukherjee/codesensei/pkg/config"
	"github.com/srijanmukherjee/codesensei/pkg/utils"
)

const (
	stdinFileName    = "stdin.txt"
	stdoutFileName   = "stdout.txt"
	stderrFileName   = "stderr.txt"
	metadataFileName = "metadata.txt"
)

type TaskContext struct {
	BoxId        string
	Cgroups      string
	Workdir      string
	Boxdir       string
	Tempdir      string
	SourceFile   string
	StdinFile    string
	StdoutFile   string
	StderrFile   string
	MetadataFile string
	Submission   models.Submission
}

func Perform(submission models.Submission) error {
	context := TaskContext{
		Submission: submission,
	}

	err := context.initialize()
	if err != nil {
		return err
	}

	err = context.compile()
	if err != nil {
		return err
	}

	return nil
}

func (ctx *TaskContext) initialize() (e error) {
	ctx.BoxId = strconv.Itoa(ctx.Submission.Id % math.MaxInt)

	err := ctx.createWorkDirectory()
	if err != nil {
		return err
	}

	ctx.Boxdir = ctx.Workdir + "/box"
	ctx.Tempdir = ctx.Workdir + "/tmp"
	ctx.SourceFile = ctx.Boxdir + "/" + ctx.Submission.Language.SourceFile
	ctx.StdinFile = ctx.Workdir + "/" + stdinFileName
	ctx.StdoutFile = ctx.Workdir + "/" + stdoutFileName
	ctx.StderrFile = ctx.Workdir + "/" + stderrFileName
	ctx.MetadataFile = ctx.Workdir + "/" + metadataFileName

	for _, file := range []string{ctx.StdinFile, ctx.StdoutFile, ctx.StderrFile, ctx.MetadataFile} {
		err = initializeFile(file)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(ctx.SourceFile, []byte(ctx.Submission.SourceCode), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *TaskContext) compile() (err error) {
	if ctx.Submission.Language.CompileCommand == "" {
		return
	}

	compileScript := ctx.Boxdir + "/compile.sh"
	compilerOptions := strings.TrimSpace(ctx.Submission.CompilerOptions)
	compilerOptions = regexp.MustCompile("[$&;<>|`]").ReplaceAllString(compilerOptions, "")
	compileScriptContent := fmt.Sprintf(ctx.Submission.Language.CompileCommand, compilerOptions)
	log.Println(compileScriptContent)
	err = os.WriteFile(compileScript, []byte(compileScriptContent), 0644)
	if err != nil {
		return
	}

	command := exec.Command(
		"isolate",
		"-s",
		"-b", ctx.BoxId, "--cg",
		"-M", ctx.MetadataFile,
		"--stderr-to-stdout",
		"-i", "/dev/null",
		"-t", utils.Str(config.Config.MaxCpuTimeLimit),
		"-x", "0",
		"-w", utils.Str(config.Config.MaxWallTimeLimit),
		"-k", utils.Str(config.Config.MaxStackLimit),
		"-p"+utils.Str(config.Config.MaxMaxProcessesAndOrThreads),
		// "-f", // TODO: add this
		"-E", "HOME=/tmp",
		"-E", "PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"",
		"-E", "LANG",
		"-E", "LANGUAGE",
		"-E", "LC_ALL",
		"-d", "/etc:noexec",
		"--run",
		"--",
		"/bin/bash", path.Base(compileScript),
	)

	log.Printf("Compiling submission %v (%v)\n", ctx.Submission.Token, ctx.Submission.Id)
	log.Println(command.String())

	commandOutput, err := command.Output()

	ctx.Submission.CompileOutput = string(commandOutput[:])

	// cleanup compilation files
	for _, file := range []string{compileScript} {
		if removeFile(file) != nil {
			return
		}
	}

	// handle compilation error or timeout
	metadata, _err := ctx.getMetadata()
	if _err != nil {
		return
	}

	_err = ctx.resetMetadataFile()
	if _err != nil {
		return _err
	}

	// compiled successfully
	if err == nil {
		return
	}

	if metadata["status"] == "TO" {
		ctx.Submission.CompileOutput = "Compilation time limit exceeded."
	}

	finishedAt, _ := time.Now().MarshalText()
	ctx.Submission.FinishedAt = string(finishedAt)
	ctx.Submission.Time = nil
	ctx.Submission.WallTime = nil
	ctx.Submission.Memory = nil
	ctx.Submission.Stdouts = nil
	ctx.Submission.Stderr = nil
	ctx.Submission.ExitCode = nil

	return
}

/* Helper functions */

func initializeFile(name string) error {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	file.Chown(os.Geteuid(), os.Getegid())

	return file.Close()
}

func removeFile(name string) (err error) {
	err = os.Chown(name, os.Geteuid(), os.Getegid())
	if err != nil {
		return
	}

	return os.Remove(name)
}

func (ctx *TaskContext) createWorkDirectory() error {
	dir, err := exec.Command(
		"isolate",
		ctx.Cgroups,
		"-b",
		ctx.BoxId,
		"--init").Output()

	if err != nil {
		return err
	}

	// NOTE: dir contains a newline at the end
	ctx.Workdir = strings.TrimRight(string(dir[:]), "\n")

	return nil
}

func (ctx *TaskContext) getMetadata() (metadata map[string]string, err error) {
	metadata = map[string]string{}

	data, err := os.ReadFile(ctx.MetadataFile)
	if err != nil {
		return
	}

	str := string(data)

	for _, line := range strings.Split(str, "\n") {
		sepIndex := strings.Index(line, ":")
		if sepIndex <= 0 {
			continue
		}

		key := line[:sepIndex]
		value := line[sepIndex+1:]
		metadata[key] = value
	}

	return
}

func (ctx *TaskContext) resetMetadataFile() (err error) {
	err = os.Remove(ctx.MetadataFile)
	if err != nil {
		return
	}
	return initializeFile(ctx.MetadataFile)
}
