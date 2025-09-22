module github.com/syntropy-cc/cooperative-grid/interfaces/cli

go 1.21

require (
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	github.com/syntropy-cc/cooperative-grid/core v0.0.0
	github.com/sirupsen/logrus v1.9.3
)

replace github.com/syntropy-cc/cooperative-grid/core => ../../core
