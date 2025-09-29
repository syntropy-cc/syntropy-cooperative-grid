package unit

import (
	"testing"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/mocks"
)

// TestNewSetupManager testa a criação de um novo SetupManager
func TestNewSetupManager(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should create setup manager successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := src.NewSetupManager()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSetupManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && manager == nil {
				t.Error("NewSetupManager() returned nil manager")
			}
		})
	}
}

// TestSetupManager_Setup testa o método Setup do SetupManager
func TestSetupManager_Setup(t *testing.T) {
	tests := []struct {
		name             string
		options          *types.SetupOptions
		mockValidator    *mocks.MockValidator
		mockConfigurator *mocks.MockConfigurator
		mockKeyManager   *mocks.MockKeyManager
		mockStateManager *mocks.MockStateManager
		mockLogger       *mocks.MockSetupLogger
		wantErr          bool
		expectedError    string
	}{
		{
			name:    "should setup successfully with valid options",
			options: helpers.CreateValidSetupOptions(),
			mockValidator: &mocks.MockValidator{
				ValidateEnvironmentFunc: func() (*types.EnvironmentInfo, error) {
					env := helpers.CreateValidEnvironmentInfo()
					env.CanProceed = true
					return env, nil
				},
			},
			mockConfigurator: &mocks.MockConfigurator{
				CreateStructureFunc: func() error {
					return nil
				},
				GenerateConfigFunc: func(options *types.ConfigOptions) error {
					return nil
				},
			},
			mockKeyManager: &mocks.MockKeyManager{
				GenerateKeyPairFunc: func(algorithm string) (*types.KeyPair, error) {
					return helpers.CreateValidKeyPair(), nil
				},
			},
			mockStateManager: &mocks.MockStateManager{
				SaveStateFunc: func(state *types.SetupState) error {
					return nil
				},
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:    "should fail when validation fails and force is false",
			options: helpers.CreateValidSetupOptions(),
			mockValidator: &mocks.MockValidator{
				ValidateEnvironmentFunc: func() (*types.EnvironmentInfo, error) {
					env := helpers.CreateValidEnvironmentInfo()
					env.CanProceed = false
					return env, nil
				},
			},
			mockConfigurator: &mocks.MockConfigurator{},
			mockKeyManager:   &mocks.MockKeyManager{},
			mockStateManager: &mocks.MockStateManager{},
			mockLogger:       &mocks.MockSetupLogger{},
			wantErr:          true,
			expectedError:    "validation_failed",
		},
		{
			name: "should proceed when validation fails but force is true",
			options: &types.SetupOptions{
				Force:        true,
				ValidateOnly: false,
				Verbose:      true,
				Quiet:        false,
				ConfigPath:   "",
				CustomSettings: map[string]string{
					"owner_name":  "Force User",
					"owner_email": "force@example.com",
				},
			},
			mockValidator: &mocks.MockValidator{
				ValidateEnvironmentFunc: func() (*types.EnvironmentInfo, error) {
					env := helpers.CreateValidEnvironmentInfo()
					env.CanProceed = false
					return env, nil
				},
			},
			mockConfigurator: &mocks.MockConfigurator{
				CreateStructureFunc: func() error {
					return nil
				},
				GenerateConfigFunc: func(options *types.ConfigOptions) error {
					return nil
				},
			},
			mockKeyManager: &mocks.MockKeyManager{
				GenerateKeyPairFunc: func(algorithm string) (*types.KeyPair, error) {
					return helpers.CreateValidKeyPair(), nil
				},
			},
			mockStateManager: &mocks.MockStateManager{
				SaveStateFunc: func(state *types.SetupState) error {
					return nil
				},
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar SetupManager com mocks
			manager := &src.SetupManager{
				Validator:    tt.mockValidator,
				Configurator: tt.mockConfigurator,
				StateManager: tt.mockStateManager,
				KeyManager:   tt.mockKeyManager,
				Logger:       tt.mockLogger,
			}

			err := manager.Setup(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedError != "" {
				helpers.AssertErrorContains(t, err, tt.expectedError)
			}
		})
	}
}
