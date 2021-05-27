package structs

// typecover:MyStruct -exclude MyField1
func testBlockWithExcludesPass() {
	m := MyStruct{}
	if true {
		m.MyField2 = "world"
	}
}

// typecover:MyStruct -exclude=MyField1
func testBlockWithExcludesFail() { // want `Type structs.MyStruct is missing MyField2`
	m := MyStruct{}
	m.myField3 = "bye"
}
