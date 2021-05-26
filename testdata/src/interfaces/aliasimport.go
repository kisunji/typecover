package interfaces

import f "io/fs"

type MyFS interface {
	f.File
}

//typecover:f.File
func aliasedImport() {
	m := &myFile{}
	m.Stat()
	m.Read(nil)
	m.Close()
}

//typecover:f.File
func aliasedImportMissingMethod() { // want `Type io/fs.File is missing Close`
	m := &myFile{}
	m.Stat()
	m.Read(nil)
}


//typecover:MyFS
func embeddedInterface() {
	m := &myFile{}
	m.Stat()
	m.Read(nil)
	m.Close()
}

//typecover:MyFS
func embeddedInterfaceMissingMethod() { // want `Type interfaces.MyFS is missing Read`
	m := &myFile{}
	m.Stat()
	m.Close()
}