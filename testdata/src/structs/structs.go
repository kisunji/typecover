package structs

import "flag"

type MyStruct struct {
	MyField1 string
	MyField2 string
	myField3 string
}

type AnotherStruct struct {
	MyField1 string
	MyField2 string
}

type MyEmbeddedStruct struct {
	MyStruct

	MyField4 string
}

// typecover:MyStruct
func testBlockFail() { // want `Type structs.MyStruct is missing MyField1`
	m := MyStruct{}
	if true {
		m.MyField2 = "world"
	}
}

// typecover:MyStruct
func testBlockPass() {
	m := MyStruct{}
	m.MyField1 = "hello"
	if true {
		m.MyField2 = "world"
	}
}

// typecover:MyStruct
func testBlockAndCompositeLiteral() {
	m := MyStruct{
		MyField1: "hello",
	}
	m.MyField2 = "world"
}

//typecover:MyStruct
func testBlockTwoStructsSameFieldNames() { // want `Type structs.MyStruct is missing MyField1`
	m := MyStruct{}
	if true {
		m.MyField2 = "world"
	}
	n := AnotherStruct{}
	n.MyField1 = "error"
}

func testStructs() {
	// typecover:MyStruct
	_ = MyStruct{
		MyField1: "hello",
		MyField2: "world",
	}

	// typecover:MyStruct
	_ = MyStruct{ // want `Type structs.MyStruct is missing MyField2`
		MyField1: "hello",
	}

	// TEST: nonvalid struct

	// typecover:MyStruct2
	_ = MyStruct{ // want `Type structs.MyStruct2 not found in associated code block`
		MyField1: "hello",
		MyField2: "world",
	}

	// TEST: imported structs

	// typecover:flag.Flag
	_ = flag.Flag{
		Name:     "",
		Usage:    "",
		Value:    nil,
		DefValue: "",
	}

	// typecover:flag.Flag
	_ = flag.Flag{ // want `Type flag.Flag is missing DefValue`
		Name:  "",
		Usage: "",
		Value: nil,
	}

	// TEST: embedded fields
	// typecover:MyEmbeddedStruct
	_ = MyEmbeddedStruct{
		MyStruct: MyStruct{
			MyField1: "",
			MyField2: "",
			myField3: "",
		},
		MyField4: "hello",
	}

	// typecover:MyEmbeddedStruct
	_ = MyEmbeddedStruct{
		// typecover:MyStruct
		MyStruct: MyStruct{ // want `Type structs.MyStruct is missing MyField2`
			MyField1: "",
		},
		MyField4: "hello",
	}

	//typecover:MyStruct
	_ = &MyStruct{
		MyField1: "hello",
		MyField2: "world",
	}

	//typecover:MyStruct
	_ = &MyStruct{ // want `Type structs.MyStruct is missing MyField1`
		MyField2: "world",
	}
}
