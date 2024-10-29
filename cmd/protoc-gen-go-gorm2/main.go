// Copyright 2024 TBONESOFT LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The protoc-gen-go-gorm2 binary is a protoc plugin to generate Go GORM code for
// proto3 version of the protocol buffer language.
//
// For more information about the usage of this plugin, see:
// https://github.com/tbonesoft/protoc-gen-go-gorm2.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"gorm.io/gorm/schema"

	"protoc-gen-go-gorm2/cmd/protoc-gen-go-gorm2/internal_gengo"
	"protoc-gen-go-gorm2/internal/editionssupport"
	"protoc-gen-go-gorm2/internal/version"
	gorm "protoc-gen-go-gorm2/proto/gorm/v1"
)

const genGoDocURL = "https://github.com/tbonesoft/protoc-gen-go-gorm2"

type DB_ENGINE int

const (
	DB_ENGINE_UNSET DB_ENGINE = iota
	DB_ENGINE_POSTGRES
)

// SupportedFeatures reports the set of supported protobuf language features.
var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL | pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), version.String())
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Fprintf(os.Stdout, "See "+genGoDocURL+" for usage information.\n")
		os.Exit(0)
	}

	var (
		flags  flag.FlagSet
		engine = flags.String("engine", "", `accept empty string "" or "postgres"`)
	)

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {

		builder := NewBuilder(gen, *engine)
		builder.Generate()

		return nil
	})
}

type builder struct {
	gen      *protogen.Plugin
	dbEngine DB_ENGINE
}

func NewBuilder(gen *protogen.Plugin, dbEngine string) *builder {
	builder := &builder{
		gen:      gen,
		dbEngine: DB_ENGINE_UNSET,
	}

	if dbEngine == "postgres" {
		builder.dbEngine = DB_ENGINE_POSTGRES
	}

	// Enable support for `optional` keyword and editions features in proto3.
	// ref: https://protobuf.dev/reference/go/go-generated/#singular-message
	// ref: https://protobuf.dev/editions/overview/
	gen.SupportedFeatures = SupportedFeatures

	gen.SupportedEditionsMinimum = editionssupport.Minimum
	gen.SupportedEditionsMaximum = editionssupport.Maximum

	return builder
}

func (b *builder) Generate() {
	for _, f := range b.gen.Files {
		if f.Generate {
			b.GenerateFile(f)
		}
	}
}

func (b *builder) GenerateFile(file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_gorm.pb.go"
	g := b.gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-gorm2. DO NOT EDIT.")
	g.P("// source: ", file.Desc.Path())
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	for _, msg := range file.Messages {
		opts := parseMessageOptions(msg)
		if !opts.Orm {
			continue
		}

		g.P(`type `, msg.GoIdent.GoName, `ORM struct {`)
		for _, field := range msg.Fields {
			b.genMessageField(g, file, msg, field)
		}
		g.P("}")

		b.genTableNameFunction(g, msg, opts)
	}

	return g
}

func parseMessageOptions(msg *protogen.Message) (rs *gorm.GormMessageOptions) {
	opts := msg.Desc.Options()
	msgOpts, ok := proto.GetExtension(opts, gorm.E_Opts).(*gorm.GormMessageOptions)
	if ok && msgOpts != nil {
		return msgOpts
	}

	return &gorm.GormMessageOptions{}
}

func parseFieldOptions(field *protogen.Field) (rs *gorm.GormFieldOptions) {
	opts := field.Desc.Options()
	fieldOpts, ok := proto.GetExtension(opts, gorm.E_Field).(*gorm.GormFieldOptions)
	if ok && fieldOpts != nil {
		return fieldOpts
	}

	return &gorm.GormFieldOptions{}
}

func (b *builder) genTableNameFunction(g *protogen.GeneratedFile, message *protogen.Message, msgOpts *gorm.GormMessageOptions) {
	g.P(`// TableName overrides the default table name generated by GORM.`)
	g.P(`func (x *`, message.GoIdent, `ORM) TableName() string {`)

	tableName := schema.NamingStrategy{}.TableName(message.GoIdent.GoName)
	if value := msgOpts.GetTable(); value != "" {
		tableName = value
	}

	g.P(`return "`, tableName, `"`)
	g.P(`}`)
}

func (b *builder) genMessageField(g *protogen.GeneratedFile, file *protogen.File, m *protogen.Message, field *protogen.Field) {
	fieldOpts := parseFieldOptions(field)
	gormTags := fieldOpts.GetTag()
	goType := b.ParseFieldGoType(g, file, field, gormTags)

	var tags internal_gengo.StructTags
	gormTagValue := fieldGormTagValue(field, gormTags)
	if len(gormTagValue) > 0 {
		tags = internal_gengo.StructTags{
			{"gorm", gormTagValue},
		}
	}

	leadingComments := internal_gengo.AppendDeprecationSuffix(field.Comments.Leading,
		field.Desc.ParentFile(),
		field.Desc.Options().(*descriptorpb.FieldOptions).GetDeprecated())

	g.P(leadingComments,
		field.GoName, " ", goType, tags,
		internal_gengo.TrailingComment(field.Comments.Trailing))
}

// ref: https://github.com/lib/pq/blob/master/array.go
const (
	PG_BOOL_ARRAY    = "pq.BoolArray"
	PG_FLOAT64_ARRAY = "pq.Float64Array"
	PG_FLOAT32_ARRAY = "pq.Float32Array"
	PG_INT64_ARRAY   = "pq.Int64Array"
	PG_INT32_ARRAY   = "pq.Int32Array"

	PG_BYTEA_ARRAY  = "pq.ByteaArray" // ref: https://www.postgresql.org/docs/17/datatype-binary.html
	PG_STRING_ARRAY = "pq.StringArray"
)

const (
	PQ_IMPORT = "github.com/lib/pq"
)

func (it *builder) ParseFieldGoType(g *protogen.GeneratedFile, file *protogen.File, field *protogen.Field, tags *gorm.GormTag) (rs string) {
	goType, pointer := internal_gengo.FieldGoType(g, file, field)

	if it.dbEngine == DB_ENGINE_POSTGRES {
		mapGoType2PgDriverType := map[string]string{
			"[]bool":    PG_BOOL_ARRAY,
			"[]float64": PG_FLOAT64_ARRAY,
			"[]float32": PG_FLOAT32_ARRAY,
			"[]int64":   PG_INT64_ARRAY,
			"[]int32":   PG_INT32_ARRAY,
			"[]bytea":   PG_BYTEA_ARRAY,
			"[]string":  PG_STRING_ARRAY,
		}

		if pgDriverType, ok := mapGoType2PgDriverType[goType]; ok {
			goType = pgDriverType
			generateImport(pgDriverType, PQ_IMPORT, g)
		}
	}

	if pointer {
		goType = "*" + goType
	}

	return goType
}

func generateImport(name string, importPath string, g *protogen.GeneratedFile) string {
	return g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       name,
		GoImportPath: protogen.GoImportPath(importPath),
	})
}

func fieldGormTagValue(field *protogen.Field, tags *gorm.GormTag) (rs string) {
	var value string

	// part 1
	column := tags.GetColumn()
	if column != "" {
		value += fmt.Sprintf("column:%s;", column)
	}

	_type := tags.GetType()
	if _type != "" {
		value += fmt.Sprintf("type:%s;", _type)
	}

	serializer := tags.GetSerializer()
	if serializer != "" {
		value += fmt.Sprintf("serializer:%s;", serializer)
	}

	size := tags.GetSize()
	if size > 0 {
		value += fmt.Sprintf("size:%d;", size)
	}

	primaryKey := tags.GetPrimaryKey()
	if primaryKey {
		value += "primaryKey;"
	}

	// part 2
	unique := tags.GetUnique()
	if unique {
		value += "unique;"
	}

	_default := tags.GetDefault()
	if _default != "" {
		value += fmt.Sprintf("default:%s;", _default)
	}

	precision := tags.GetPrecision()
	if precision > 0 {
		value += fmt.Sprintf("precision:%d;", precision)
	}

	scale := tags.GetScale()
	if scale > 0 {
		value += fmt.Sprintf("scale:%d;", scale)
	}

	notNull := tags.GetNotNull()
	if notNull {
		value += "not null;"
	}

	value = strings.TrimRight(value, ";")

	// TODO: part 3
	return value
}
