package template

const (
	StartFile = "main"

	GoMod = `module {{nameModule}}

go 1.14

require github.com/fatih/color v1.9.0`

	Main = `package main

import (
    "os"
	"{{nameModule}}/pkg/{{nameModule}}"
)

func main() {
	input1 := os.Getenv("SAMPLE_TEXT")
	input2 := os.Getenv("SAMPLE_LIST")
	input3 := os.Getenv("SAMPLE_BOOL")

	{{nameModule}}.Input{
    	Text:    input1,
    	List:    input2,
    	Boolean: input3,
    }.Run()
}`

	Makefile = `# Go parameters
BIN_FOLDER=../bin
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat
BIN_NAME=main
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
CMD_PATH=./main.go
BIN_FOLDER_DARWIN=$(BIN_FOLDER)/darwin
BIN_DARWIN=$(BIN_FOLDER_DARWIN)/$(BIN_NAME)
BIN_FOLDER_LINUX=$(BIN_FOLDER)/linux
BIN_LINUX=$(BIN_FOLDER_LINUX)/$(BIN_NAME)
BIN_FOLDER_WINDOWS=$(BIN_FOLDER)/windows
BIN_WINDOWS=$(BIN_FOLDER_WINDOWS)/$(BIN_NAME).exe


build: go-build sh-unix bat-windows

go-build:
	mkdir -p $(BIN_FOLDER_DARWIN) $(BIN_FOLDER_LINUX) $(BIN_FOLDER_WINDOWS)
	export MODULE=$(GO111MODULE=on go list -m)
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o '$(BIN_LINUX)' -v $(CMD_PATH)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o '$(BIN_DARWIN)' -v $(CMD_PATH)
	#WINDOWS 64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o '$(BIN_WINDOWS)' -v $(CMD_PATH)

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'if [ $$(uname) = "Darwin" ]; then' >> $(SH)
	echo '  ./darwin/$(BIN_NAME)' >> $(SH)
	echo 'else' >> $(SH)
	echo '  ./linux/$(BIN_NAME)' >> $(SH)
	echo 'fi' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo 'cd windows' >> $(BAT)
	echo 'start $(BIN_NAME).exe' >> $(BAT)

test:
	$(GOTEST) -short ` + "`go list ./... | grep -v vendor/`"

	Dockerfile = `
FROM golang:alpine AS builder

ADD . /app
WORKDIR /app
RUN go build -o main -v main.go

FROM alpine:latest


COPY --from=builder /app/main main
COPY --from=builder /app/set_umask.sh set_umask.sh
RUN chmod +x main
RUN chmod +x set_umask.sh

WORKDIR /app
ENTRYPOINT ["/set_umask.sh"]
CMD ["/main"]`

	Pkg = `package {{nameModule}}

import (
	"fmt"
	"github.com/fatih/color"
)

type Input struct {
	Text string
	List string
	Boolean string
}

func(in Input)Run()  {
	fmt.Println("Hello world!")
	color.Green(fmt.Sprintf("You receive %s in text.", in.Text ))
	color.Red(fmt.Sprintf("You receive %s in list.", in.List ))
	color.Yellow(fmt.Sprintf("You receive %s in boolean.", in.Boolean ))
}`

	WindowsBuild = `:: Go parameters
echo off
SETLOCAL
SET BINARY_NAME={{bin-name}}
SET GOCMD=go
SET GOBUILD=%GOCMD% build
SET GOTEST=%GOCMD% test
SET CMD_PATH=main.go
SET DIST=..\dist
SET DIST_MAC_DIR=%DIST%\darwin\bin
SET BIN_MAC=%BINARY_NAME%-darwin
SET DIST_LINUX_DIR=%DIST%\linux\bin
SET BIN_LINUX=%BINARY_NAME%-linux
SET DIST_WIN_DIR=%DIST%\windows\bin
SET BIN_WIN=%BINARY_NAME%-windows.exe
:build
	mkdir %DIST_MAC_DIR% 
    mkdir %DIST_LINUX_DIR% 
    mkdir %DIST_WIN_DIR%
    SET GO111MODULE=on
    for /f %%i in ('go list -m') do set MODULE=%%i
    CALL :linux
    CALL :darwin
    CALL :windows
    GOTO DONE
:linux
    SET CGO_ENABLED=0 
    SET GOOS=linux
    SET GOARCH=amd64
    %GOBUILD% -tags release -o %DIST_LINUX_DIR%\%BIN_LINUX% -v %CMD_PATH% && xcopy . %DIST_LINUX_DIR% /E /H /C /I && xcopy ..\config.json %DIST_LINUX_DIR%\..\
    GOTO DONE
:darwin
    SET CGO_ENABLED=
	SET GOOS=darwin
    SET GOARCH=amd64
    %GOBUILD% -tags release -o %DIST_MAC_DIR%\%BIN_MAC% -v %CMD_PATH%  && xcopy . %DIST_MAC_DIR% /E /H /C /I && xcopy ..\config.json %DIST_MAC_DIR%\..\
    GOTO DONE
:windows
    SET CGO_ENABLED=
	SET GOOS=windows
    SET GOARCH=amd64
    %GOBUILD% -tags release -o %DIST_WIN_DIR%\%BIN_WIN% -v %CMD_PATH% && xcopy . %DIST_WIN_DIR% /E /H /C /I && xcopy ..\config.json %DIST_WIN_DIR%\..\
    GOTO DONE
:DONE`
)
