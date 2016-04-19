package syntax

import (
	"fmt"
	"testing"

	"github.com/gobs/pretty"
)

func Test(t *testing.T) {
	syn, err := Load("testdata/Go.sublime-syntax")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", pretty.PrettyFormat(syn))
}
