package structs

import (
	"myvendor.org/somepkg"
)

// typecover:somepkg.Foo
func testImportPathWithDotSuccess() {
	f := &somepkg.Foo{
		Bar: "bar",
		Baz: "baz",
	}
	_ = f
}

// typecover:somepkg.Foo
func testImportPathWithDotFailure() { // want `Type myvendor.org/somepkg.Foo is missing Baz`
	f := &somepkg.Foo{
		Bar: "bar",
	}
	_ = f
}
