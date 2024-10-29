@echo off

@REM
@REM build the CLI protoc-gen-go-gorm2
@REM

protoc ^
--go_out=proto/ ^
--go_opt=paths=source_relative ^
--proto_path=proto/ ^
proto/gorm/v1/gorm.proto

go install -v github.com/tbonesoft/protoc-gen-go-gorm2/cmd/protoc-gen-go-gorm2

@REM
@REM build the example bookstore
@REM

if not exist "gen" (
    mkdir "gen"
)

protoc ^
--go_out=gen/ ^
--go_opt=paths=source_relative ^
--go-gorm2_out=gen/ ^
--go-gorm2_opt=paths=source_relative ^
--go-gorm2_opt=engine=postgres ^
--proto_path=proto/ ^
proto/bookstore/v1/bookstore.proto
