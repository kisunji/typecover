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

type MyStructWithMethods struct {
	MyField string
}

func (m MyStructWithMethods) ValueMethod() int {
	return 0
}
func (m *MyStructWithMethods) PointerMethod() int {
	return 0
}

type MyEmbeddedStructWithMethods struct {
	MyStructWithMethods
}

func (m MyEmbeddedStructWithMethods) ValueMethod2() int {
	return 0
}

func (m *MyEmbeddedStructWithMethods) PointerMethod2() int {
	return 0
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
func testBlockFail() { // want `Type structs.MyStruct is missing MyField1`
	m := MyStruct{}
	if true {
		m.MyField2 = "world"
	}
}

// typecover:MyStruct
func testBlockPointerPass() {
	m := &MyStruct{}
	m.MyField1 = "hello"
	m.MyField2 = "world"
	return
}

// typecover:MyStruct
func testBlockPointerFail() { // want `Type structs.MyStruct is missing MyField2`
	m := &MyStruct{}
	m.MyField1 = "hello"
	return
}

// typecover:MyStruct
func testBlockAndCompositeLiteral() {
	m := MyStruct{
		MyField1: "hello",
	}
	m.MyField2 = "world"
}

// typecover:MyStruct
func testBlockTwoStructsSameFieldNames() { // want `Type structs.MyStruct is missing MyField1`
	m := MyStruct{}
	if true {
		m.MyField2 = "world"
	}
	n := AnotherStruct{}
	n.MyField1 = "error"
}

func testAssigningFromAllFields() {
	m := MyStruct{
		MyField1: "hello",
		MyField2: "world",
	}
	type local struct {
		MyField1 string
		MyField2 string
	}
	// typecover:MyStruct
	_ = local{
		MyField1: m.MyField1,
		MyField2: m.MyField2,
	}

	// typecover:MyStruct
	_ = local{ // want `Type structs.MyStruct is missing MyField2`
		MyField1: m.MyField1,
		MyField2: "world",
	}

	type localWPointers struct {
		MyField1 *string
		MyField2 *string
	}
	// typecover:MyStruct
	_ = localWPointers{
		MyField1: &m.MyField1,
		MyField2: &m.MyField2,
	}

	// typecover:MyStruct
	_ = localWPointers{ // want `Type structs.MyStruct is missing MyField2`
		MyField1: &m.MyField1,
		MyField2: nil,
	}
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
	_ = MyStruct{ // want `Type structs.MyStruct2 not found in project scope`
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

	// typecover:MyStruct
	_ = &MyStruct{
		MyField1: "hello",
		MyField2: "world",
	}

	// typecover:MyStruct
	_ = &MyStruct{ // want `Type structs.MyStruct is missing MyField1`
		MyField2: "world",
	}

	// typecover:MyStructWithMethods
	{

		m := &MyStructWithMethods{
			MyField: "",
		}
		_ = m.ValueMethod()
		_ = m.PointerMethod()
	}

	// typecover:MyStructWithMethods
	_ = &MyStructWithMethods{ // want `Type structs.MyStructWithMethods is missing ValueMethod, PointerMethod`
		MyField: "",
	}

	//typecover:MyEmbeddedStructWithMethods
	{
		m := &MyEmbeddedStructWithMethods{
			MyStructWithMethods: MyStructWithMethods{
				MyField: "",
			},
		}
		_ = m.PointerMethod()
		_ = m.PointerMethod2()
		_ = m.ValueMethod()
		_ = m.ValueMethod2()
	}

	//typecover:MyEmbeddedStructWithMethods
	{ // want `Type structs.MyEmbeddedStructWithMethods is missing PointerMethod`
		m := &MyEmbeddedStructWithMethods{
			MyStructWithMethods: MyStructWithMethods{
				MyField: "",
			},
		}
		_ = m.PointerMethod2()
		_ = m.ValueMethod()
		_ = m.ValueMethod2()
	}
}
