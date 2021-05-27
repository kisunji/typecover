# typecover

`typecover` is a go linter that checks if a code block is assigning to all exported fields of a struct or 
calling all exported methods of an interface.

It is useful in cases where code wants to be aware of any newly added members.

## Install
```
go get -u github.com/kisunji/typecover/cmd/typecover
```

## Usage

Using the CLI
```
typecover [package/file]
```

Comment directives
```
# Local Type
// typecover:TypeName

# Imported Type
// typecover:pkg.TypeName 

# Excluding members
// typecover:TypeName -exclude Field3,Field4,Field5
```
`typecover:YourType` will check for the existence of all exported members of `YourType` in the comment's associated code
block (for details on how comments are associated with code, see [here](https://golang.org/pkg/go/ast/#NewCommentMap)).

## Examples

### Struct assignment
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

### Covering a code block
The `typecover` annotation can be placed at a higher level (func, if-stmt, for-loop) to cover the whole block.
```go
// typecover:MyStruct
func example() {
    m := MyStruct{}
    m.MyField2 = "world"    
}
```
```
Type example.MyStruct is missing MyField1
```

### Using imported types
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

### Copying fields from exported type
```go
//typecover:flag.Flag
func cloneFlag
f := MyFlag{
	Name: f.Name,
	Usage: f.Usage,
	Value: f.Value,	
}
```

## Credits

https://github.com/mbilski/exhaustivestruct

https://github.com/reillywatson/enumcover