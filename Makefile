# !!!MAKE SURE YOUR GOPATH ENVIRONMENT VARIABLE IS SET FIRST!!!
# Any issues with this file, let me know / make a PR, I haven't tested it completely but it should be close enough.

IMPLANT=GwopImplant
CLITOOL=GwopCli
DIR=out
LDFLAGS=-ldflags "-s -w"
D64=Darwin-x64
L64=Linux-x64
W64=Windows-x64

# Make Directory to store executables
$(shell mkdir -p ${DIR})

# Change default to just make for the host OS and add MAKE ALL to do this
default: cli-linux cli-darwin implant-windows implant-linux

all: default

# Compile Darwin binaries
darwin: cli-darwin implant-darwin

# Compile Linux Binaries
linux: cli-linux implant-linux

windows: implant-windows

# Compile Implant - Linux x64
implant-linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${IMPLANT}-${L64} cmd/implant/main.go

# Compile Implant - Windows x64     REPLACE LDFLAGS + WINAGENTLDFLAGS for actual release!!
implant-windows:
	export GOOS=windows GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${IMPLANT}-${W64}.exe cmd/implant/main.go

# Compile Implant - MacOS
implant-darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${IMPLANT}-${D64} cmd/implant/main.go

# Compile Cli - MacOS
cli-darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${CLITOOL}-${D64} cmd/clitool/main.go

# Compile Cli - Linux x64
cli-linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${CLITOOL}-${L64} cmd/clitool/main.go

# Compile Cli - Windows x64
cli-windows:
	export GOOS=windows;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${CLITOOL}-${W64}.exe cmd/clitool/main.go

clean:
	rm -rf ${DIR}*