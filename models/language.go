package models

type Language struct {
	Id             int
	Name           string
	CompileCommand string
	RunCommand     string
	SourceFile     string
}
