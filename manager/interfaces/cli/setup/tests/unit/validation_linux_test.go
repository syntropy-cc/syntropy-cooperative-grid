//go:build linux
// +build linux

package tests

import (
	"testing"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup"
)

func TestValidateLinuxEnvironment(t *testing.T) {
	// Test basic validation
	result, err := setup.ValidateLinuxEnvironment(false)
	if err != nil {
		t.Fatalf("ValidateLinuxEnvironment failed: %v", err)
	}

	// Check if result is not nil
	if result == nil {
		t.Fatal("Expected non-nil validation result")
	}

	// Check if environment information is populated
	if result.Environment.OS == "" {
		t.Error("Expected OS to be populated")
	}

	if result.Environment.Architecture == "" {
		t.Error("Expected Architecture to be populated")
	}

	if result.Environment.HomeDir == "" {
		t.Error("Expected HomeDir to be populated")
	}
}

func TestValidateLinuxEnvironmentWithForce(t *testing.T) {
	// Test validation with force flag
	result, err := setup.ValidateLinuxEnvironment(true)
	if err != nil {
		t.Fatalf("ValidateLinuxEnvironment with force failed: %v", err)
	}

	// Check if result is not nil
	if result == nil {
		t.Fatal("Expected non-nil validation result")
	}

	// With force flag, validation should pass even with warnings
	if !result.Valid && len(result.Warnings) > 0 {
		t.Error("Expected validation to pass with force flag despite warnings")
	}
}
