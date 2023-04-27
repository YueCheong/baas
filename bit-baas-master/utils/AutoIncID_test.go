package utils

import (
	"fmt"
	"testing"
)

func TestAutoIncIDGen(t *testing.T) {
	dg := NewAutoIncID()
	if dg.GetStart() != 1 ||
		dg.GetStep() != 1 ||
		dg.GetCurrentID() != NoIdCreated {
		t.FailNow()
	}
	if dg.GenID() != 1 ||
		dg.GenID() != 2 ||
		dg.GenID() != 3 ||
		dg.currentID != 3 {
		t.FailNow()
	}
	fmt.Println("Default gen passed")

	g := NewAutoIncID(WithStartID(5), WithStep(3))
	if g.GetStart() != 5 ||
		g.GetStep() != 3 ||
		g.currentID != NoIdCreated {
		t.FailNow()
	}

	if g.GenID() != 5 ||
		g.GenID() != 8 ||
		g.GenID() != 11 ||
		g.currentID != 11 {
		t.FailNow()
	}
	fmt.Println("Gen start 5 step 3 passed")
}
