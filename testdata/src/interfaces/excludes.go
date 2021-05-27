package interfaces

// typecover:MyInterface -exclude MyFunc1,MyFunc2
func sequentialMethodCallsWithExcludes() {
	_ = &myStruct{}
}

// typecover:MyInterface -exclude=MyFunc1
func sequentialMethodCallsWithExcludesFails() { // want `Type interfaces.MyInterface is missing MyFunc2`
	_ = &myStruct{}
}
