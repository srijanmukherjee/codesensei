package server

import (
	"regexp"
	"testing"

	"github.com/srijanmukherjee/codesensei/pkg/server"
)

func TestSayHello(t *testing.T) {
	output := "Hello, World"
	want := regexp.MustCompile(`\b` + output + `\b`)
	msg := server.SayHello()
	if !want.MatchString(msg) {
		t.Fatalf(`SayHello() = %q, want match for %#q, nil`, msg, want)
	}
}
