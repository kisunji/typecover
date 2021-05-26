package interfaces

// typecover:MyInterface ~MyFunc1
func sequentialMethodCallsWithExcludes() {
	m := &myStruct{}
	m.MyFunc2()
}

// typecover:MyInterface ~MyFunc1
func sequentialMethodCallsWithExcludesFails() { // want `Type interfaces.MyInterface is missing MyFunc2`
	_ = &myStruct{}
}
