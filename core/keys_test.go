package core_test

import (
	"testing"

	"github.com/thelotter-enterprise/usergo/core"
)

func TestBuildWithPrefix(t *testing.T) {
	want := "prefix-name-a-b-c"
	k := core.NewKeys()
	is := k.Build("pRefiX", "nAme", "A", "B", "C")

	if want != is {
		t.Fail()
	}
}

func TestBuildWithoutPrefix(t *testing.T) {
	want := "name-a-b-c"
	k := core.NewKeys()
	is := k.Build("", "nAme", "A", "B", "C")

	if want != is {
		t.Fail()
	}
}
