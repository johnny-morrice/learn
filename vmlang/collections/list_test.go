package collections

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	list := List[string]{}
	listFoo := list.Append("foo")
	fooSl1 := listFoo.Slice()
	if !reflect.DeepEqual(fooSl1, []string{"foo"}) {
		t.Fatalf("unexpected list foo 1: %v\n%v", fooSl1, listFoo)
	}

	listBar := listFoo.Append("bar")
	barSl1 := listBar.Slice()
	if !reflect.DeepEqual(barSl1, []string{"foo", "bar"}) {
		t.Fatalf("unexpected list bar 1: %v\n%v", barSl1, listBar)
	}
	fooSl2 := listFoo.Slice()
	if !reflect.DeepEqual(fooSl2, []string{"foo"}) {
		t.Fatalf("unexpected list foo 2: %v", fooSl2)
	}

	listBaz := listBar.Append("baz")

	expectedBazNode := &listNode[string]{
		Value: "baz",
	}

	expectedBarNode := &listNode[string]{
		Value: "bar",
		Next:  expectedBazNode,
	}

	expectedFooNode := &listNode[string]{
		Value: "foo",
		Next:  expectedBarNode,
	}

	if !reflect.DeepEqual(listBaz.zeroNode, expectedFooNode) {
		t.Errorf("unexpected node: %v\nactual: %v", expectedFooNode, listBaz.zeroNode)
	}

	bazSl1 := listBaz.Slice()
	if !reflect.DeepEqual(bazSl1, []string{"foo", "bar", "baz"}) {
		t.Errorf("unexpected list baz 1: %v\n%v", bazSl1, listBaz)
	}

	barSl2 := listBar.Slice()
	if !reflect.DeepEqual(barSl2, []string{"foo", "bar"}) {
		t.Errorf("unexpected list bar 2: %v", barSl2)
	}
	fooSl3 := listFoo.Slice()
	if !reflect.DeepEqual(fooSl3, []string{"foo"}) {
		t.Errorf("unexpected list foo 3: %v", fooSl3)
	}
}
