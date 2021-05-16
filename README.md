# typecover

`typecover` is a go linter that checks if a struct is assigned all the exported fields of its type.

## Install
```
go get -u github.com/kisunji/typecover/cmd/typecover
```

## Usage
```
typecover [package/file]
```

## Examples
`typecover` will check that all exported fields are assigned
```go
type MyStruct struct {
	MyField1 string
	MyField2 string
	myField3 string
}

// typecover:MyStruct
m := MyStruct{ 
    MyField1: "hello",
}
```

```
Type example.MyStruct is missing MyField2
```


`typecover` works with imported structs as well
```go
// typecover:flag.Flag
f := flag.Flag{ 
    Name:  "test",
    Usage: "usage instructions",
    Value: nil,
}
```

```
Type flag.Flag is missing DefValue
```

## Credits

https://github.com/mbilski/exhaustivestruct

https://github.com/reillywatson/enumcover