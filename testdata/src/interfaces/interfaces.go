package interfaces

import "io/fs"

type MyInterface interface {
	MyFunc1()
	MyFunc2()
	myFunc3()
}

type MyBuilder interface {
	Option1() MyBuilder
	Option2() MyBuilder
	Build() *myStruct
}

func NewMyBuilder() MyBuilder {
	return &myStruct{}
}

type myStruct struct{}

func (m *myStruct) MyFunc1() {}

func (m *myStruct) MyFunc2() {}

func (m *myStruct) myFunc3() {}

func (m *myStruct) Option1() MyBuilder { return m }

func (m *myStruct) Option2() MyBuilder { return m }

func (m *myStruct) Build() *myStruct { return m }

type myFile struct{}

func (m *myFile) Stat() (fs.FileInfo, error) { panic("implement me") }

func (m *myFile) Read(bytes []byte) (int, error) { panic("implement me") }

func (m *myFile) Close() error { panic("implement me") }

// typecover:MyInterface
func sequentialMethodCalls() {
	m := &myStruct{}
	m.MyFunc1()
	m.MyFunc2()
}

// typecover:MyInterface
func sequentialMethodCallsMissingCall() { // want `Type interfaces.MyInterface is missing MyFunc2`
	m := &myStruct{}
	m.MyFunc1()
}

// typecover:MyBuilder
func builderMethodChaining() {
	m := &myStruct{}
	_ = m.Option1().Option2().Build()
}

// typecover:MyBuilder
func builderMethodChainingMissingCall() { // want `Type interfaces.MyBuilder is missing Option1`
	m := &myStruct{}
	_ = m.Option2().Build()
}

// typecover:MyBuilder
func builderWorksInsideComposite() {
	type local struct{
		b MyBuilder
	}
	_= local{
		b: NewMyBuilder().Option1().Option2().Build(),
	}
}

// typecover:MyBuilder
func builderWorksInsideCompositeMissingCall() { // want `Type interfaces.MyBuilder is missing Option2`
	type local struct{
		b MyBuilder
	}
	_= local{
		b: NewMyBuilder().Option1().Build(),
	}
}

// typecover:MyInterface
// typecover:MyBuilder
func multipleInterfacesCovered() {
	m := &myStruct{}
	m.Option1().Option2()
	m.MyFunc1()
	m.MyFunc2()
	_ = m.Build()
}

// typecover:MyInterface
// typecover:MyBuilder
func multipleInterfacesCoveredButMissingCalls() { // want `Type interfaces.MyInterface is missing MyFunc1` `Type interfaces.MyBuilder is missing Option1`
	m := &myStruct{}
	m.Option2()
	m.MyFunc2()
	_ = m.Build()
}

// typecover:MyInterface
func nilInterface() {
	var m MyInterface
	m.MyFunc1()
	m.MyFunc2()
}

//typecover:fs.File
func importedInterface() {
	m := &myFile{}
	m.Stat()
	m.Read(nil)
	m.Close()
}

// typecover:MyInterface
func complexCodeBlock() {
	var m MyInterface
	// typecover:MyBuilder
	if true {
		b := NewMyBuilder()
		m = b.Option1().
			Option2().Build()
	}
	m.MyFunc1()
	f := func() {
		m.MyFunc2()
	}
	f()
}

// typecover:MyInterface
func complexCodeBlockMissingMyBuilderMethod() {
	var m MyInterface
	// typecover:MyBuilder
	if true { // want `Type interfaces.MyBuilder is missing Option2`
		b := NewMyBuilder()
		m = b.Option1().
			Build()
	}
	m.MyFunc1()
	f := func() {
		m.MyFunc2()
	}
	f()
}

// typecover:MyInterface
func complexCodeBlockMissingMyInterfaceMethod() { // want `Type interfaces.MyInterface is missing MyFunc1`
	var m MyInterface
	// typecover:MyBuilder
	if true {
		b := NewMyBuilder()
		m = b.Option1().Option2().
			Build()
	}
	f := func() {
		m.MyFunc2()
	}
	f()
}
