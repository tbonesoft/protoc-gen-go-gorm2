@echo off

protoc ^
--go_out=proto/ ^
--go_opt=paths=source_relative ^
--proto_path=proto/ ^
proto/gorm/v1/gorm.proto


set CGO_ENABLED=0
go build -o protoc-gen-go-gorm2.exe cmd\protoc-gen-go-gorm2\main.go

if not exist "gen" (
    mkdir "gen"
)

protoc ^
--go_out=gen/ ^
--go_opt=paths=source_relative ^
--go-gorm2_out=gen/ ^
--go-gorm2_opt=paths=source_relative ^
--go-gorm2_opt=engine=postgres ^
--proto_path=examples/bookstore/proto/ ^
--proto_path=proto/ ^
examples/bookstore/proto/bookstore/v1/bookstore.proto
