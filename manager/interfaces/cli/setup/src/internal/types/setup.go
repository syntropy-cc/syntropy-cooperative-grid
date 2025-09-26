package types

import (
	"errors"
	"time"
)

// SetupOptions defines the options for the setup process
type SetupOptions struct {
	Force          bool   // Force setup even if validations fail
	InstallService bool   // Install system service
	ConfigPath     string // Custom configuration file path
	HomeDir        string // Custom home directory
}

// SetupResult contains the result of the setup process
type SetupResult struct {
	Success     bool         // Indicates if the setup was successful
	StartTime   time.Time    // Setup start time
	EndTime     time.Time    // Setup end time
	ConfigPath  string       // Configuration file path
	Environment string       // Environment (windows, linux, darwin)
	Options     SetupOptions // Options used in the setup
	Error       error        // Error, if any
	Message     string       // Human-readable message
}

// ErrNotImplemented is returned when a functionality is not implemented
var ErrNotImplemented = errors.New("functionality not implemented for this operating system")
