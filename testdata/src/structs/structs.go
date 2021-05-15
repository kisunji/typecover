package structs

import "flag"

type MyStruct struct {
	MyField1 string
	MyField2 string
	myField3 string
}

type MyEmbeddedStruct struct {
	MyStruct

	MyField4 string
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

	// TEST: anonymous structs

	// typecover:MyStruct
	_ = MyStruct{"hello", "world", "!"} // this should pass

	// TEST: embedded fields
	// typecover:MyEmbeddedStruct
	_ = MyEmbeddedStruct{
		MyStruct{
			MyField1: "",
			MyField2: "",
			myField3: "",
		},
		"hello",
	}

	// typecover:MyEmbeddedStruct
	_ = MyEmbeddedStruct{
		// typecover:MyStruct
		MyStruct{ // want `Type structs.MyStruct is missing MyField2`
			MyField1: "",
		},
		"hello",
	}
}
