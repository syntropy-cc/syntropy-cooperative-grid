module syntropy-cc/cooperative-grid/interfaces/cli

go 1.24.0

toolchain go1.24.7

require (
	github.com/spf13/cobra v1.8.0
	syntropy-cc/cooperative-grid/infrastructure v0.0.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace syntropy-cc/cooperative-grid/infrastructure => ../../infrastructure
