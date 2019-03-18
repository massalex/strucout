# Go strucout package
Prints structure info into stdout

---

#### How to use
For structure example, see the code below.   

+ Simle\
prints name, type, value fields by default
```go
strucout.New(from).Out()
````
+ All columns\
prints name, type, value, tags, package, offset, anonymous, index fields
```go
strucout.New(from).AllColumns().Out()
```
+ Set columns for output\
adds package and offset to output
 ```go
 so := strucout.New(from)
 so.Flags |= strucout.ShowOffset | strucout.ShowPackage
 so.Out()
 ````
+ Change column output format\
sets name column format: 40 symbols width, red color and right align
 ```go
 strucout.New(from).ChangeColumn("name", 40, strucout.ColorRed, false).Out()
 ````
 + Tag set \
 Sets `json` filter for tags output (for `Id` ouput `person_id`, for `Role` - `part`)\
 The `ShowTags` flag is set automatically 
 ```go
strucout.New(from).SetTag("json").Out()
````
---
        
####Structure for example 	
```go
type Person struct {
    Name string
    Number int
}

type Table struct {
    Id int32 `json:"person_id"`
    Role string `json:"part"`
    Factor float64
    Age uint
    Interface interface{}
    Check bool
    Person Person
    Complex complex64
    Fnc func() string
    Mp map[string] string
    Sl []string
    Ar [10]int
    Pt *int
}

iPtr := 10

from := &Table{
    Id: 56988965,
    Role: "Admin",
    Factor: 253.35,
    Interface: 342543,
    Check: true,
    Fnc: func() string { return "Hello" },
    Sl: []string{"I","am","strucout"},
    Ar: [10]int{1,2,3},
    Pt: &iPtr,
}

from.Mp = make(map[string] string)
from.Mp["A"] = "Z"
from.Mp["Y"] = "B"
```