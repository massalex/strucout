package strucout

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	ColorBlack    = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorPurpure
	ColorLightBlue
	ColorWhite
	ShowName      = 1
	ShowPackage   = 1 << 1
	ShowType      = 1 << 2
	ShowTags      = 1 << 3
	ShowOffset    = 1 << 4
	ShowIndex     = 1 << 5
	ShowAnonymous = 1 << 6
	ShowValue     = 1 << 7
)

type Format struct {
	Width int
	Color int
	Left bool
	flag int
}

type outFormat map[string]Format

type instance struct {
	str     interface{}
	value   reflect.Value
	Flags   int
	tag     string
	valid   bool
	maxFlag int
	Format  outFormat
	columns []string
}

func New(str interface{}) *instance {
	strucot := &instance{
		str: str,
		value: indirect(reflect.ValueOf(str)),
		Flags: ShowName | ShowType | ShowValue,
	}
	strucot.columns = []string{"name","type","value","tags","package","offset","anonymous","index"}
	strucot.Format = make(outFormat)
	strucot.Format["name"] 		= Format{20,ColorYellow,true,ShowName }
	strucot.Format["type"] 		= Format{20,ColorLightBlue,true,ShowType }
	strucot.Format["value"] 	= Format{16,ColorGreen,false,ShowValue }
	strucot.Format["package"] 	= Format{10,ColorPurpure,true,ShowPackage }
	strucot.Format["tags"] 		= Format{10,ColorLightBlue,true,ShowTags }
	strucot.Format["offset"] 	= Format{10,ColorRed,false,ShowOffset }
	strucot.Format["index"] 	= Format{10,ColorWhite,false,ShowIndex }
	strucot.Format["anonymous"] = Format{11,ColorWhite,false,ShowAnonymous }
	return strucot
}

func (strucout *instance) IsValid() bool {
	strucout.valid = strucout.value.IsValid()
	return strucout.valid
}

func (strucout *instance) SetTag(tag string) *instance {
	if tag == "" {
		strucout.Flags &^= ShowTags
	} else {
		strucout.Flags |= ShowTags
	}
	strucout.tag = tag
	return strucout
}

func (strucout *instance) AllColumns() *instance {
	strucout.Flags = ShowValue ^ (ShowValue - 1)
	return strucout
}

func (strucout *instance) DropColors() *instance {
	for _, column := range strucout.columns {
		if elem, is := strucout.Format[column]; is {
			elem.Color = ColorWhite
			strucout.Format[column] = elem
		}
	}
	return strucout
}

func (strucout *instance) ChangeColumn(field string, width int, color int, left bool) *instance {
	if elem, is := strucout.Format[field]; is {
		elem.Width = width
		elem.Color = color
		elem.Left = left
		strucout.Format[field] = elem
	}
	return strucout
}

func (strucout *instance) Out() {
	if !strucout.IsValid() {
		panic("Structure is invalid")
	}
	border := "+"
	title := "|"
	field := ""
	outColumns := [0]interface{}{}
	oc := outColumns[:]
	outFields := [0]interface{}{}
	of := outFields[:]
	for _, f := range strucout.columns {
		format := strucout.Format[f]
		if format.flag & strucout.Flags > 0 {
			border += strings.Repeat("-", format.Width + 2) + "+"
			title +=  " \x1b[7;" + strconv.Itoa(format.Color) + "m%"
			if format.Left {
				title += "-"
			}
			title += strconv.Itoa(format.Width) + "s\x1b[0m |"
			oc = append(oc, " " + strings.Title(f) + " ")
			of = append(of, "-")
		}
	}

	field = regexp.MustCompile(`\[7;`).ReplaceAllLiteralString(title, "[")
	fmt.Println("  " + strucout.value.String())
	fmt.Println(border)
	fmt.Printf(title + "\n", oc...)
	fmt.Println(border)
	valueType := reflect.TypeOf(strucout.str)
	structFields := getFields(valueType)
	j := 0
	for i, structField := range structFields {
		j = 0
		if strucout.Format["name"].flag & strucout.Flags > 0 {
			of[j] = structField.Name
			j++
		}
		if strucout.Format["type"].flag & strucout.Flags > 0 {
			of[j] = structField.Type
			j++
		}
		if strucout.Format["value"].flag & strucout.Flags > 0 {
			val := ""
			if strucout.value.Kind() != reflect.Struct {
				val = "<not struct>"
			} else {
				f := strucout.value.Field(i)
				switch f.Kind() {
				case reflect.String:
					val = fmt.Sprintf("%v", f.String())
				case reflect.Interface:
					if f.CanInterface() {
						val = fmt.Sprintf("%v", f.Interface())
					} else {
						val = "<interface value>"
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val = fmt.Sprintf("%d", f.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val = fmt.Sprintf("%d", f.Uint())
				case reflect.Float32, reflect.Float64:
					val = fmt.Sprintf("%f", f.Float())
				case reflect.Complex64, reflect.Complex128:
					val = fmt.Sprintf("%v", f.Complex())
				case reflect.Bool:
					val = fmt.Sprintf("%v", f.Bool())
				case reflect.Struct, reflect.Array, reflect.Slice:
					val = f.String()
					stringMethod := f.MethodByName("String") // default or set String method (alt strucout.value.FieldByName(structField.Name))
					if stringMethod.IsValid() {
						v := stringMethod.Call([]reflect.Value{})
						if len(v) > 0 {
							val = v[0].String()
						}
					}
				case reflect.Map:
					val = fmt.Sprintf("%v", f.MapKeys())
				case reflect.Ptr:
					val = fmt.Sprintf("%#v", f.Pointer())
				default:
					val = "<unknown>"
				}
			}
			of[j] = val
			j++
		}
		if strucout.Format["tags"].flag & strucout.Flags > 0 {
			if strucout.tag == "" {
				of[j] = structField.Tag
			} else {
				of[j] = structField.Tag.Get(strucout.tag)
			}
			j++
		}
		if strucout.Format["package"].flag & strucout.Flags > 0 {
			of[j] = structField.PkgPath
			j++
		}
		if strucout.Format["offset"].flag & strucout.Flags > 0 {
			of[j] = fmt.Sprintf("%v", structField.Offset)
			j++
		}
		if strucout.Format["index"].flag & strucout.Flags > 0 {
			of[j] = fmt.Sprintf("%v", structField.Anonymous)
			j++
		}
		if strucout.Format["index"].flag & strucout.Flags > 0 {
			of[j] = fmt.Sprintf("%d", structField.Index)
			j++
		}
		fmt.Printf(field + "\n", of...) //structField.Name, structField.Type, structField.Anonymous, val)
	}
	fmt.Println(border)
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func getFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, getFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}
	return fields
}