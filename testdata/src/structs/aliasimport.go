package structs

import f "flag"

func testAliasedImport() {
	//typecover:f.Flag
	_ = f.Flag{
		Name:     "",
		Usage:    "",
		Value:    nil,
		DefValue: "",
	}

	//typecover:f.Flag
	_ = f.Flag{ // want `Type flag.Flag is missing DefValue`
		Name:  "",
		Usage: "",
		Value: nil,
	}
}
