// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal_gengo is internal to the protobuf module.
package internal_gengo

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/tbonesoft/protoc-gen-go-gorm2/internal/encoding/tag"
)

// FieldGoType returns the Go type used for a field.
//
// If it returns pointer=true, the struct field is a pointer to the type.
func FieldGoType(g *protogen.GeneratedFile, file *protogen.File, field *protogen.Field) (goType string, pointer bool) {
	if field.Desc.IsWeak() {
		return "struct{}", false
	}

	pointer = field.Desc.HasPresence()
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = g.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false // rely on nullability of slices for presence
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + g.QualifiedGoIdent(field.Message.GoIdent)
		pointer = false // pointer captured as part of the type
	}

	switch {
	case field.Desc.IsList():
		return "[]" + goType, false
	case field.Desc.IsMap():
		keyType, _ := FieldGoType(g, file, field.Message.Fields[0])
		valType, _ := FieldGoType(g, file, field.Message.Fields[1])
		return fmt.Sprintf("map[%v]%v", keyType, valType), false
	}
	return goType, pointer
}

// AppendDeprecationSuffix optionally appends a deprecation notice as a suffix.
func AppendDeprecationSuffix(prefix protogen.Comments, parentFile protoreflect.FileDescriptor, deprecated bool) protogen.Comments {
	fileDeprecated := parentFile.Options().(*descriptorpb.FileOptions).GetDeprecated()
	if !deprecated && !fileDeprecated {
		return prefix
	}
	if prefix != "" {
		prefix += "\n"
	}
	if fileDeprecated {
		return prefix + " Deprecated: The entire proto file " + protogen.Comments(parentFile.Path()) + " is marked as deprecated.\n"
	}
	return prefix + " Deprecated: Marked as deprecated in " + protogen.Comments(parentFile.Path()) + ".\n"
}

func FieldJSONTagValue(field *protogen.Field) string {
	return string(field.Desc.Name()) + ",omitempty"
}

func FieldProtobufTagValue(field *protogen.Field) string {
	var enumName string
	if field.Desc.Kind() == protoreflect.EnumKind {
		enumName = protoimpl.X.LegacyEnumName(field.Enum.Desc)
	}
	return tag.Marshal(field.Desc, enumName)
}

// TrailingComment is like protogen.Comments, but lacks a trailing newline.
type TrailingComment protogen.Comments

func (c TrailingComment) String() string {
	s := strings.TrimSuffix(protogen.Comments(c).String(), "\n")
	if strings.Contains(s, "\n") {
		// We don't support multi-lined trailing comments as it is unclear
		// how to best render them in the generated code.
		return ""
	}
	return s
}

// StructTags is a data structure for build idiomatic Go struct tags.
// Each [2]string is a key-value pair, where value is the unescaped string.
//
// Example: StructTags{{"key", "value"}}.String() -> `key:"value"`
type StructTags [][2]string

func (tags StructTags) String() string {
	if len(tags) == 0 {
		return ""
	}
	var ss []string
	for _, tag := range tags {
		// NOTE: When quoting the value, we need to make sure the backtick
		// character does not appear. Convert all cases to the escaped hex form.
		key := tag[0]
		val := strings.Replace(strconv.Quote(tag[1]), "`", `\x60`, -1)
		ss = append(ss, fmt.Sprintf("%s:%s", key, val))
	}
	return "`" + strings.Join(ss, " ") + "`"
}
