package gocmp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type Struct struct {
	ID int64
}

type PrivateFieldContainer struct {
	ID   int64
	Name string

	metadata Metadata
}

type Metadata struct {
	ID int
}

func TestGoCmp(t *testing.T) {
	sWant, sGot := Struct{1}, Struct{1}

	// same
	if diff := cmp.Diff(sWant, sGot); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	// diff(ignore field)
	sWant, sGot = Struct{1}, Struct{2}
	if diff := cmp.Diff(sWant, sGot, cmpopts.IgnoreFields(Struct{}, "ID")); diff != "" {
		t.Logf("mismatch (-want +got):\n%s", diff)
	}

	want, got := PrivateFieldContainer{}, PrivateFieldContainer{}

	// cannot handle unexported field(ignore field)
	opt := cmpopts.IgnoreFields(PrivateFieldContainer{}, "metadata")
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	// cannot handle unexported field(ignore unexported)
	opt = cmpopts.IgnoreUnexported(PrivateFieldContainer{})
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	want, got = PrivateFieldContainer{
		ID:       1,
		Name:     "name",
		metadata: Metadata{1},
	}, PrivateFieldContainer{
		ID:       1,
		Name:     "name",
		metadata: Metadata{1},
	}

	// allow unexported(same)
	opt = cmp.AllowUnexported(PrivateFieldContainer{})
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	// allow unexported(diff)
	got.metadata.ID = 2
	opt = cmp.AllowUnexported(PrivateFieldContainer{})
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Logf("mismatch (-want +got):\n%s", diff)
	}

}
