package engine

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/srijanmukherjee/codesensei/shared/config"
)

// TODO:
// - [ ] Create custom errors and wrap the errors returned by the standard library functions

const (
	ISOLATE_CMD             = "/usr/local/bin/isolate"
	BOXDIR                  = "/box"
	TMPDIR                  = "/tmp"
	RUN_SCRIPT_FILENAME     = "run.sh"
	COMPILE_SCRIPT_FILENAME = "compile.sh"
	COMPILE_OUTPUT_FILENAME = "compile_output.txt"
	METADATA_FILENAME       = "metadata.txt"
	STDIN_FILENAME          = "stdin.txt"
	STDOUT_FILENAME         = "stdout.txt"
	STDERR_FILENAME         = "stderr.txt"
)

type CodeExecutionContext struct {
	Config config.Config

	Id              uint32
	SourceFilename  string
	SourceContent   string
	CompilerOptions string
	CompileCommand  string
	RunCommand      string

	CpuTimeLimit                         int
	CpuExtraTime                         int
	WallTimeLimit                        int
	StackLimit                           int
	MaxProcessesAndOrThreads             int
	MemoryLimit                          int
	MaxFileSize                          int
	EnableNetwork                        bool
	RedirectStderrToStdout               bool
	EnablePerProcessAndThreadTimeLimit   bool
	EnablePerProcessAndThreadMemoryLimit bool

	boxId             string
	workdir           string
	boxdir            string
	compileScriptFile string
	compileOutputFile string
	runScriptFile     string
	stdinFile         string
	stdoutFile        string
	stderrFile        string
	sourceFile        string
	metadataFile      string
	cgroupsEnabled    bool
}

type RunResult struct {
	Stdout     string
	Stderr     string
	Time       float64
	WallTime   int
	Memory     int
	ExitCode   int
	ExitSignal int
	Message    string
	Status     string
}

type CompileResult struct {
	Output   string
	Success  bool
	ExitCode int
}

// Initializes the code execution context by creating a new isolate box
func (c *CodeExecutionContext) Init() error {
	c.boxId = fmt.Sprint(c.Id % 2147483647)
	c.cgroupsEnabled = !c.EnablePerProcessAndThreadMemoryLimit || !c.EnablePerProcessAndThreadTimeLimit

	args := []string{
		"-b", c.boxId, "--init",
	}

	if c.cgroupsEnabled {
		args = append(args, "--cg")
	}

	workdir, err := exec.Command(ISOLATE_CMD, args...).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Println(string(exitErr.Stderr))
		}
		log.Printf("%v", err)
		return err
	}

	c.workdir = strings.TrimSpace(string(workdir))
	c.boxdir = path.Join(c.workdir, BOXDIR)
	c.runScriptFile = path.Join(c.workdir, BOXDIR, RUN_SCRIPT_FILENAME)
	c.compileScriptFile = path.Join(c.workdir, BOXDIR, COMPILE_SCRIPT_FILENAME)
	c.compileOutputFile = path.Join(c.workdir, COMPILE_OUTPUT_FILENAME)
	c.sourceFile = path.Join(c.workdir, BOXDIR, c.SourceFilename)
	c.metadataFile = path.Join(c.workdir, METADATA_FILENAME)
	c.stdinFile = path.Join(c.workdir, STDIN_FILENAME)
	c.stdoutFile = path.Join(c.workdir, STDOUT_FILENAME)
	c.stderrFile = path.Join(c.workdir, STDERR_FILENAME)

	if err := c.initializeFile(c.sourceFile); err != nil {
		return err
	}

	if err := os.WriteFile(c.sourceFile, []byte(c.SourceContent), 0666); err != nil {
		return err
	}

	if err := c.initializeFile(c.metadataFile); err != nil {
		return err
	}

	if err := c.initializeFile(c.stdinFile); err != nil {
		return err
	}

	if err := os.WriteFile(c.runScriptFile, []byte(c.RunCommand), 0755); err != nil {
		return err
	}

	return nil
}

// Compiles the source code in the code execution context
func (c *CodeExecutionContext) Compile() (CompileResult, error) {
	compilerOptions := sanitizeCompilerOptions(c.CompilerOptions)
	compileCommand := fmt.Sprintf(c.CompileCommand, compilerOptions)

	if err := c.initializeFile(c.compileScriptFile); err != nil {
		return CompileResult{}, err
	}

	if err := os.WriteFile(c.compileScriptFile, []byte(compileCommand), 0755); err != nil {
		return CompileResult{}, err
	}

	cmd := exec.Command(
		ISOLATE_CMD,
		"-s",
		"-b", c.boxId,
		"-M", c.metadataFile,
		"--stderr-to-stdout",
		"-i", "/dev/null",
		"-x", "0",
		"-t", fmt.Sprint(c.Config.MaxCpuTimeLimit),
		"-w", fmt.Sprint(c.Config.MaxWallTimeLimit),
		"--stack", fmt.Sprint(c.Config.MaxStackLimit),
		"-f", fmt.Sprint(c.Config.MaxMaxFileSize), // TODO: MAX_MAX_FILE_SIZE
		fmt.Sprintf("-p%v", c.Config.MaxMaxProcessesAndOrThreads),
		fmt.Sprintf("-m%v", c.Config.MaxMemoryLimit),
		"-E", "HOME=/tmp",
		"-E", "PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/usr/libexec:/bin\"",
		"-E", "LANG",
		"-E", "LANGUAGE",
		"-E", "LC_ALL",
		"-d", "/etc:noexec",
		"--run",
		"--",
		"/bin/bash",
		COMPILE_SCRIPT_FILENAME,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if !errors.Is(err, &exec.ExitError{}) {
			return CompileResult{}, err
		}
	}

	exitCode := cmd.ProcessState.ExitCode()
	metadata := c.getMetadata()

	c.resetMetadata()

	// TODO: delete the compile script

	if exitCode == 0 {
		return CompileResult{
			Output:   string(output),
			Success:  true,
			ExitCode: exitCode,
		}, nil
	}

	if metadata["status"] == "TO" {
		return CompileResult{
			Output:   "Compilation time limit exceeded",
			Success:  false,
			ExitCode: exitCode,
		}, nil
	}

	return CompileResult{
		Output:   string(output),
		Success:  false,
		ExitCode: exitCode,
	}, nil
}

// Runs the compiled code in the code execution context with the given input
func (c *CodeExecutionContext) Run(stdin string) (RunResult, error) {

	if err := os.WriteFile(c.stdinFile, []byte(stdin), 0666); err != nil {
		return RunResult{}, err
	}

	args := []string{
		"-v",
		"-b", c.boxId,
		"-M", c.metadataFile,
		"-t", fmt.Sprint(c.CpuTimeLimit),
		"-x", fmt.Sprint(c.CpuExtraTime),
		"-w", fmt.Sprint(c.WallTimeLimit),
		"-k", fmt.Sprint(c.StackLimit),
		"-f", fmt.Sprint(c.MaxFileSize),
		fmt.Sprintf("-p%d", c.MaxProcessesAndOrThreads),
		"-E", "HOME=/tmp",
		"-E", "PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/usr/libexec:/bin\"",
		"-E", "LANG",
		"-E", "LANGUAGE",
		"-E", "LC_ALL",
		"-d", "/etc:noexec",
	}

	if c.cgroupsEnabled {
		args = append(args, "--cg")
	}

	if c.RedirectStderrToStdout {
		args = append(args, "--stderr-to-stdout")
	}

	if c.EnableNetwork {
		args = append(args, "--share-net")
	}

	if c.EnablePerProcessAndThreadMemoryLimit {
		args = append(args, fmt.Sprintf("-m%d", c.MemoryLimit))
	} else {
		args = append(args, fmt.Sprintf("--cg-mem=%d", c.MemoryLimit))
	}

	if c.EnablePerProcessAndThreadTimeLimit {
		if c.cgroupsEnabled {
			args = append(args, "--no-cg-timeing")
		}
	} else {
		args = append(args, "--cg-timing")
	}

	args = append(
		args,
		"--run",
		"--",
		"/bin/bash",
		RUN_SCRIPT_FILENAME)

	cmd := exec.Command(ISOLATE_CMD, args...)
	log.Println(cmd.String())

	stdinBuffer := bytes.Buffer{}
	// TODO: handle error
	stdinBuffer.WriteString(stdin)
	cmd.Stdin = &stdinBuffer

	stdoutBuffer := bytes.Buffer{}
	cmd.Stdout = &stdoutBuffer

	stderrBuffer := bytes.Buffer{}
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running isolate: %v", err)
	}

	log.Println(string(stdoutBuffer.Bytes()))

	return RunResult{}, nil
}

func (c *CodeExecutionContext) Cleanup() error {
	err := exec.Command(ISOLATE_CMD, "-b", c.boxId, "--cleanup").Run()
	return err
}

func (c *CodeExecutionContext) initializeFile(filename string) error {
	return exec.Command(
		"/bin/bash",
		"-c",
		fmt.Sprintf("sudo touch %s && sudo chown $(whoami): %s", filename, filename),
	).Run()
}

func (c *CodeExecutionContext) getMetadata() map[string]string {
	content, err := os.ReadFile(c.metadataFile)
	if err != nil {
		log.Printf("Error reading metadata file: %v", err)
		return nil
	}

	lines := strings.Split(string(content), "\n")
	metadata := make(map[string]string)

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			metadata[parts[0]] = strings.TrimSpace(parts[1])
		}
	}

	return metadata
}

func (c *CodeExecutionContext) resetMetadata() {
	err := os.WriteFile(c.metadataFile, nil, 0666)
	if err != nil {
		log.Printf("Error resetting metadata file: %v", err)
	}
}

func sanitizeCompilerOptions(compilerOptions string) string {
	if !utf8.ValidString(compilerOptions) {
		compilerOptions = strings.ToValidUTF8(compilerOptions, "")
	}

	compilerOptions = strings.TrimSpace(compilerOptions)

	unsafeChars := regexp.MustCompile("[$&;<>|`]")
	compilerOptions = unsafeChars.ReplaceAllString(compilerOptions, "")
	return compilerOptions
}
