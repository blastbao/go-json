package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type opType struct {
	Op                    string
	Code                  string
	Escaped               func() string
	HeadToPtrHead         func() string
	HeadToNPtrHead        func() string
	HeadToAnonymousHead   func() string
	HeadToOmitEmptyHead   func() string
	HeadToStringTagHead   func() string
	HeadToOnlyHead        func() string
	PtrHeadToHead         func() string
	FieldToEnd            func() string
	FieldToOmitEmptyField func() string
	FieldToStringTagField func() string
}

func (t opType) IsEscaped() bool {
	return t.Op != t.Escaped()
}

func (t opType) IsHeadToPtrHead() bool {
	return t.Op != t.HeadToPtrHead()
}

func (t opType) IsHeadToNPtrHead() bool {
	return t.Op != t.HeadToNPtrHead()
}

func (t opType) IsHeadToAnonymousHead() bool {
	return t.Op != t.HeadToAnonymousHead()
}

func (t opType) IsHeadToOmitEmptyHead() bool {
	return t.Op != t.HeadToOmitEmptyHead()
}

func (t opType) IsHeadToStringTagHead() bool {
	return t.Op != t.HeadToStringTagHead()
}

func (t opType) IsPtrHeadToHead() bool {
	return t.Op != t.PtrHeadToHead()
}

func (t opType) IsHeadToOnlyHead() bool {
	return t.Op != t.HeadToOnlyHead()
}

func (t opType) IsFieldToEnd() bool {
	return t.Op != t.FieldToEnd()
}

func (t opType) IsFieldToOmitEmptyField() bool {
	return t.Op != t.FieldToOmitEmptyField()
}

func (t opType) IsFieldToStringTagField() bool {
	return t.Op != t.FieldToStringTagField()
}

func createOpType(op, code string) opType {
	return opType{
		Op:                    op,
		Code:                  code,
		Escaped:               func() string { return op },
		HeadToPtrHead:         func() string { return op },
		HeadToNPtrHead:        func() string { return op },
		HeadToAnonymousHead:   func() string { return op },
		HeadToOmitEmptyHead:   func() string { return op },
		HeadToStringTagHead:   func() string { return op },
		HeadToOnlyHead:        func() string { return op },
		PtrHeadToHead:         func() string { return op },
		FieldToEnd:            func() string { return op },
		FieldToOmitEmptyField: func() string { return op },
		FieldToStringTagField: func() string { return op },
	}
}

func _main() error {
	tmpl, err := template.New("").Parse(`// Code generated by cmd/generator. DO NOT EDIT!
package json

import (
  "strings"
)

type codeType int

const (
{{- range $index, $type := .CodeTypes }}
  code{{ $type }} codeType = {{ $index }}
{{- end }}
)

type opType int

const (
{{- range $index, $type := .OpTypes }}
  op{{ $type.Op }} opType = {{ $index }}
{{- end }}
)

func (t opType) String() string {
  switch t {
{{- range $type := .OpTypes }}
  case op{{ $type.Op }}:
    return "{{ $type.Op }}"
{{- end }}
  }
  return ""
}

func (t opType) codeType() codeType {
  if strings.Contains(t.String(), "Struct") {
    if strings.Contains(t.String(), "End") {
      return codeStructEnd
    }
    return codeStructField
  }
  if strings.Contains(t.String(), "ArrayHead") {
    return codeArrayHead
  }
  if strings.Contains(t.String(), "ArrayElem") {
    return codeArrayElem
  }
  if strings.Contains(t.String(), "SliceHead") {
    return codeSliceHead
  }
  if strings.Contains(t.String(), "SliceElem") {
    return codeSliceElem
  }
  if strings.Contains(t.String(), "MapHead") {
    return codeMapHead
  }
  if strings.Contains(t.String(), "MapKey") {
    return codeMapKey
  }
  if strings.Contains(t.String(), "MapValue") {
    return codeMapValue
  }
  if strings.Contains(t.String(), "MapEnd") {
    return codeMapEnd
  }

  return codeOp
}

func (t opType) toEscaped() opType {
  if strings.Index(t.String(), "Escaped") > 0 {
    return t
  }
  if t.String() == "String" {
    return opType(int(t) + 1)
  }

  fieldHeadIdx := strings.Index(t.String(), "Head")
  if fieldHeadIdx > 0 && strings.Contains(t.String(), "Struct") {
    const toEscapedHeadOffset = 36
    return opType(int(t) + toEscapedHeadOffset)
  }
  fieldIdx := strings.Index(t.String(), "Field")
  if fieldIdx > 0 && strings.Contains(t.String(), "Struct") {
    const toEscapedFieldOffset = 3
    return opType(int(t) + toEscapedFieldOffset)
  }
  if strings.Contains(t.String(), "StructEnd") {
    const toEscapedEndOffset = 3
    return opType(int(t) + toEscapedEndOffset)
  }
  return t
}

func (t opType) headToPtrHead() opType {
  if strings.Index(t.String(), "PtrHead") > 0 {
    return t
  }
  if strings.Index(t.String(), "PtrAnonymousHead") > 0 {
    return t
  }

  idx := strings.Index(t.String(), "Field")
  if idx == -1 {
    return t
  }
  suffix := "Ptr"+t.String()[idx+len("Field"):]

  const toPtrOffset = 12
  if strings.Contains(opType(int(t) + toPtrOffset).String(), suffix) {
    return opType(int(t) + toPtrOffset)
  }
  return t
}

func (t opType) headToNPtrHead() opType {
  if strings.Index(t.String(), "PtrHead") > 0 {
    return t
  }
  if strings.Index(t.String(), "PtrAnonymousHead") > 0 {
    return t
  }

  idx := strings.Index(t.String(), "Field")
  if idx == -1 {
    return t
  }
  suffix := "NPtr"+t.String()[idx+len("Field"):]

  const toPtrOffset = 24
  if strings.Contains(opType(int(t) + toPtrOffset).String(), suffix) {
    return opType(int(t) + toPtrOffset)
  }
  return t
}

func (t opType) headToAnonymousHead() opType {
  const toAnonymousOffset = 6
  if strings.Contains(opType(int(t) + toAnonymousOffset).String(), "Anonymous") {
    return opType(int(t) + toAnonymousOffset)
  }
  return t
}

func (t opType) headToOmitEmptyHead() opType {
  const toOmitEmptyOffset = 2
  if strings.Contains(opType(int(t) + toOmitEmptyOffset).String(), "OmitEmpty") {
    return opType(int(t) + toOmitEmptyOffset)
  }

  return t
}

func (t opType) headToStringTagHead() opType {
  const toStringTagOffset = 4
  if strings.Contains(opType(int(t) + toStringTagOffset).String(), "StringTag") {
    return opType(int(t) + toStringTagOffset)
  }
  return t
}

func (t opType) headToOnlyHead() opType {
  if strings.HasSuffix(t.String(), "Head") || strings.HasSuffix(t.String(), "HeadOmitEmpty") || strings.HasSuffix(t.String(), "HeadStringTag") {
    return t
  }

  const toOnlyOffset = 1
  if opType(int(t) + toOnlyOffset).String() == t.String() + "Only" {
    return opType(int(t) + toOnlyOffset)
  }
  return t
}

func (t opType) ptrHeadToHead() opType {
  idx := strings.Index(t.String(), "Ptr")
  if idx == -1 {
    return t
  }
  suffix := t.String()[idx+len("Ptr"):]

  const toPtrOffset = 12
  if strings.Contains(opType(int(t) - toPtrOffset).String(), suffix) {
    return opType(int(t) - toPtrOffset)
  }
  return t
}

func (t opType) fieldToEnd() opType {
  switch t {
{{- range $type := .OpTypes }}
{{- if $type.IsFieldToEnd }}
  case op{{ $type.Op }}:
    return op{{ call $type.FieldToEnd }}
{{- end }}
{{- end }}
  }
  return t
}

func (t opType) fieldToOmitEmptyField() opType {
  const toOmitEmptyOffset = 1
  if strings.Contains(opType(int(t) + toOmitEmptyOffset).String(), "OmitEmpty") {
    return opType(int(t) + toOmitEmptyOffset)
  }
  return t
}

func (t opType) fieldToStringTagField() opType {
  const toStringTagOffset = 2
  if strings.Contains(opType(int(t) + toStringTagOffset).String(), "StringTag") {
    return opType(int(t) + toStringTagOffset)
  }
  return t
}

`)
	if err != nil {
		return err
	}
	codeTypes := []string{
		"Op",
		"ArrayHead",
		"ArrayElem",
		"SliceHead",
		"SliceElem",
		"MapHead",
		"MapKey",
		"MapValue",
		"MapEnd",
		"StructFieldRecursive",
		"StructField",
		"StructEnd",
	}
	primitiveTypes := []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "bool", "string", "escapedString", "bytes",
		"array", "map", "mapLoad", "slice", "struct", "MarshalJSON", "MarshalText", "recursive",
		"intString", "int8String", "int16String", "int32String", "int64String",
		"uintString", "uint8String", "uint16String", "uint32String", "uint64String",
		"intPtr", "int8Ptr", "int16Ptr", "int32Ptr", "int64Ptr",
		"uintPtr", "uint8Ptr", "uint16Ptr", "uint32Ptr", "uint64Ptr",
		"float32Ptr", "float64Ptr", "boolPtr", "stringPtr", "escapedStringPtr", "bytesPtr",
		"intNPtr", "int8NPtr", "int16NPtr", "int32NPtr", "int64NPtr",
		"uintNPtr", "uint8NPtr", "uint16NPtr", "uint32NPtr", "uint64NPtr",
		"float32NPtr", "float64NPtr", "boolNPtr", "stringNPtr", "escapedStringNPtr", "bytesNPtr",
	}
	primitiveTypesUpper := []string{}
	for _, typ := range primitiveTypes {
		primitiveTypesUpper = append(primitiveTypesUpper, strings.ToUpper(string(typ[0]))+typ[1:])
	}
	opTypes := []opType{
		createOpType("End", "Op"),
		createOpType("Interface", "Op"),
		createOpType("InterfaceEnd", "Op"),
		createOpType("Ptr", "Op"),
		createOpType("NPtr", "Op"),
		createOpType("SliceHead", "SliceHead"),
		createOpType("RootSliceHead", "SliceHead"),
		createOpType("SliceElem", "SliceElem"),
		createOpType("RootSliceElem", "SliceElem"),
		createOpType("SliceEnd", "Op"),
		createOpType("ArrayHead", "ArrayHead"),
		createOpType("ArrayElem", "ArrayElem"),
		createOpType("ArrayEnd", "Op"),
		createOpType("MapHead", "MapHead"),
		createOpType("MapHeadLoad", "MapHead"),
		createOpType("MapKey", "MapKey"),
		createOpType("MapValue", "MapValue"),
		createOpType("MapEnd", "Op"),
		createOpType("StructFieldRecursiveEnd", "Op"),
		createOpType("StructAnonymousEnd", "StructEnd"),
	}
	for _, typ := range primitiveTypesUpper {
		typ := typ
		optype := createOpType(typ, "Op")
		switch typ {
		case "String", "StringPtr", "StringNPtr":
			optype.Escaped = func() string {
				return fmt.Sprintf("Escaped%s", typ)
			}
		}
		opTypes = append(opTypes, optype)
	}
	for _, typ := range append(primitiveTypesUpper, "") {
		if typ == "EscapedString" || typ == "EscapedStringPtr" || typ == "EscapedStringNPtr" {
			continue
		}
		for _, escapedOrNot := range []string{"", "Escaped"} {
			for _, ptrOrNot := range []string{"", "Ptr", "NPtr"} {
				for _, headType := range []string{"", "Anonymous"} {
					for _, opt := range []string{"", "OmitEmpty", "StringTag"} {
						for _, onlyOrNot := range []string{"", "Only"} {
							escapedOrNot := escapedOrNot
							ptrOrNot := ptrOrNot
							headType := headType
							opt := opt
							typ := typ
							onlyOrNot := onlyOrNot

							isEscaped := escapedOrNot != ""
							isString := typ == "String" || typ == "StringPtr" || typ == "StringNPtr"

							if isEscaped && isString {
								typ = "Escaped" + typ
							}

							op := fmt.Sprintf(
								"Struct%sField%s%sHead%s%s%s",
								escapedOrNot,
								ptrOrNot,
								headType,
								opt,
								typ,
								onlyOrNot,
							)
							opTypes = append(opTypes, opType{
								Op:   op,
								Code: "StructField",
								Escaped: func() string {
									switch typ {
									case "String", "StringPtr", "StringNPtr":
										return fmt.Sprintf(
											"StructEscapedField%s%sHead%sEscaped%s%s",
											ptrOrNot,
											headType,
											opt,
											typ,
											onlyOrNot,
										)
									}
									return fmt.Sprintf(
										"StructEscapedField%s%sHead%s%s%s",
										ptrOrNot,
										headType,
										opt,
										typ,
										onlyOrNot,
									)
								},
								HeadToPtrHead: func() string {
									return fmt.Sprintf(
										"Struct%sFieldPtr%sHead%s%s%s",
										escapedOrNot,
										headType,
										opt,
										typ,
										onlyOrNot,
									)
								},
								HeadToNPtrHead: func() string {
									return fmt.Sprintf(
										"Struct%sFieldNPtr%sHead%s%s%s",
										escapedOrNot,
										headType,
										opt,
										typ,
										onlyOrNot,
									)
								},
								HeadToAnonymousHead: func() string {
									return fmt.Sprintf(
										"Struct%sField%sAnonymousHead%s%s%s",
										escapedOrNot,
										ptrOrNot,
										opt,
										typ,
										onlyOrNot,
									)
								},
								HeadToOmitEmptyHead: func() string {
									return fmt.Sprintf(
										"Struct%sField%s%sHeadOmitEmpty%s%s",
										escapedOrNot,
										ptrOrNot,
										headType,
										typ,
										onlyOrNot,
									)
								},
								HeadToStringTagHead: func() string {
									return fmt.Sprintf(
										"Struct%sField%s%sHeadStringTag%s%s",
										escapedOrNot,
										ptrOrNot,
										headType,
										typ,
										onlyOrNot,
									)
								},
								HeadToOnlyHead: func() string {
									switch typ {
									case "", "Array", "Map", "MapLoad", "Slice",
										"Struct", "Recursive", "MarshalJSON", "MarshalText",
										"IntString", "Int8String", "Int16String", "Int32String", "Int64String",
										"UintString", "Uint8String", "Uint16String", "Uint32String", "Uint64String":
										return op
									}
									return fmt.Sprintf(
										"Struct%sField%s%sHead%s%sOnly",
										escapedOrNot,
										ptrOrNot,
										headType,
										opt,
										typ,
									)
								},
								PtrHeadToHead: func() string {
									return fmt.Sprintf(
										"Struct%sField%sHead%s%s%s",
										escapedOrNot,
										headType,
										opt,
										typ,
										onlyOrNot,
									)
								},
								FieldToEnd:            func() string { return op },
								FieldToOmitEmptyField: func() string { return op },
								FieldToStringTagField: func() string { return op },
							})
						}
					}
				}
			}
		}
	}
	for _, typ := range append(primitiveTypesUpper, "") {
		if typ == "EscapedString" || typ == "EscapedStringPtr" || typ == "EscapedStringNPtr" {
			continue
		}
		for _, escapedOrNot := range []string{"", "Escaped"} {
			for _, opt := range []string{"", "OmitEmpty", "StringTag"} {
				escapedOrNot := escapedOrNot
				opt := opt
				typ := typ

				isEscaped := escapedOrNot != ""
				isString := typ == "String" || typ == "StringPtr" || typ == "StringNPtr"

				if isEscaped && isString {
					typ = "Escaped" + typ
				}

				op := fmt.Sprintf(
					"Struct%sField%s%s",
					escapedOrNot,
					opt,
					typ,
				)
				opTypes = append(opTypes, opType{
					Op:   op,
					Code: "StructField",
					Escaped: func() string {
						switch typ {
						case "String", "StringPtr", "StringNPtr":
							return fmt.Sprintf(
								"StructEscapedField%sEscaped%s",
								opt,
								typ,
							)
						}
						return fmt.Sprintf(
							"StructEscapedField%s%s",
							opt,
							typ,
						)
					},
					HeadToPtrHead:       func() string { return op },
					HeadToNPtrHead:      func() string { return op },
					HeadToAnonymousHead: func() string { return op },
					HeadToOmitEmptyHead: func() string { return op },
					HeadToStringTagHead: func() string { return op },
					HeadToOnlyHead:      func() string { return op },
					PtrHeadToHead:       func() string { return op },
					FieldToEnd: func() string {
						switch typ {
						case "", "Array", "Map", "MapLoad", "Slice", "Struct", "Recursive":
							return op
						}
						return fmt.Sprintf(
							"Struct%sEnd%s%s",
							escapedOrNot,
							opt,
							typ,
						)
					},
					FieldToOmitEmptyField: func() string {
						return fmt.Sprintf(
							"Struct%sFieldOmitEmpty%s",
							escapedOrNot,
							typ,
						)
					},
					FieldToStringTagField: func() string {
						return fmt.Sprintf(
							"Struct%sFieldStringTag%s",
							escapedOrNot,
							typ,
						)
					},
				})
			}
		}
	}
	for _, typ := range append(primitiveTypesUpper, "") {
		if typ == "EscapedString" || typ == "EscapedStringPtr" || typ == "EscapedStringNPtr" {
			continue
		}
		for _, escapedOrNot := range []string{"", "Escaped"} {
			for _, opt := range []string{"", "OmitEmpty", "StringTag"} {
				escapedOrNot := escapedOrNot
				opt := opt
				typ := typ

				isEscaped := escapedOrNot != ""
				isString := typ == "String" || typ == "StringPtr" || typ == "StringNPtr"

				if isEscaped && isString {
					typ = "Escaped" + typ
				}

				op := fmt.Sprintf(
					"Struct%sEnd%s%s",
					escapedOrNot,
					opt,
					typ,
				)
				opTypes = append(opTypes, opType{
					Op:   op,
					Code: "StructEnd",
					Escaped: func() string {
						switch typ {
						case "String", "StringPtr", "StringNPtr":
							return fmt.Sprintf(
								"StructEscapedEnd%sEscaped%s",
								opt,
								typ,
							)
						}
						return fmt.Sprintf(
							"StructEscapedEnd%s%s",
							opt,
							typ,
						)
					},
					HeadToPtrHead:         func() string { return op },
					HeadToNPtrHead:        func() string { return op },
					HeadToAnonymousHead:   func() string { return op },
					HeadToOmitEmptyHead:   func() string { return op },
					HeadToStringTagHead:   func() string { return op },
					HeadToOnlyHead:        func() string { return op },
					PtrHeadToHead:         func() string { return op },
					FieldToEnd:            func() string { return op },
					FieldToOmitEmptyField: func() string { return op },
					FieldToStringTagField: func() string { return op },
				})
			}
		}
	}
	var b bytes.Buffer
	if err := tmpl.Execute(&b, struct {
		CodeTypes []string
		OpTypes   []opType
		OpLen     int
	}{
		CodeTypes: codeTypes,
		OpTypes:   opTypes,
		OpLen:     len(opTypes),
	}); err != nil {
		return err
	}
	path := filepath.Join(repoRoot(), "encode_optype.go")
	buf, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, buf, 0644)
}

func repoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	relativePathFromRepoRoot := filepath.Join("cmd", "generator")
	return strings.TrimSuffix(filepath.Dir(file), relativePathFromRepoRoot)
}

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}
