package mocks

import (
	"fmt"
	"sync"
)

// MockTokenManager simula o Setup TokenManager para testes
type MockTokenManager struct {
	mu sync.RWMutex

	// Configuração de comportamento
	GridToken     string
	TokenValid    bool
	RefreshError  error
	GetTokenError error

	// Histórico de chamadas
	GetTokenCalls     []GetTokenCall
	ValidateCalls     []ValidateCall
	RefreshCalls      []RefreshCall
	GetTokenCallCount int
	ValidateCallCount int
	RefreshCallCount  int
}

// GetTokenCall representa uma chamada para GetGridToken
type GetTokenCall struct {
	Timestamp int64
	Result    string
	Error     error
}

// ValidateCall representa uma chamada para ValidateToken
type ValidateCall struct {
	Timestamp int64
	Token     string
	Error     error
}

// RefreshCall representa uma chamada para RefreshToken
type RefreshCall struct {
	Timestamp int64
	Error     error
}

// NewMockTokenManager cria um novo mock do TokenManager
func NewMockTokenManager() *MockTokenManager {
	return &MockTokenManager{
		GridToken:     "mock-grid-token-12345",
		TokenValid:    true,
		GetTokenCalls: make([]GetTokenCall, 0),
		ValidateCalls: make([]ValidateCall, 0),
		RefreshCalls:  make([]RefreshCall, 0),
	}
}

// GetGridToken simula a obtenção do Grid Token
func (m *MockTokenManager) GetGridToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.GetTokenCallCount++
	call := GetTokenCall{
		Timestamp: getCurrentTimestamp(),
		Result:    m.GridToken,
		Error:     m.GetTokenError,
	}
	m.GetTokenCalls = append(m.GetTokenCalls, call)

	if m.GetTokenError != nil {
		return "", m.GetTokenError
	}

	return m.GridToken, nil
}

// ValidateToken simula a validação de um token
func (m *MockTokenManager) ValidateToken(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ValidateCallCount++
	call := ValidateCall{
		Timestamp: getCurrentTimestamp(),
		Token:     token,
		Error:     nil,
	}

	// Validação básica
	if token == "" {
		call.Error = fmt.Errorf("token cannot be empty")
		m.ValidateCalls = append(m.ValidateCalls, call)
		return call.Error
	}

	if token != m.GridToken && m.TokenValid {
		call.Error = fmt.Errorf("invalid token: %s", token)
		m.ValidateCalls = append(m.ValidateCalls, call)
		return call.Error
	}

	m.ValidateCalls = append(m.ValidateCalls, call)
	return nil
}

// RefreshToken simula o refresh do token
func (m *MockTokenManager) RefreshToken() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RefreshCallCount++
	call := RefreshCall{
		Timestamp: getCurrentTimestamp(),
		Error:     m.RefreshError,
	}
	m.RefreshCalls = append(m.RefreshCalls, call)

	return m.RefreshError
}

// Configuração de comportamento do mock

// SetGridToken define o token que será retornado
func (m *MockTokenManager) SetGridToken(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GridToken = token
}

// SetTokenValid define se o token deve ser considerado válido
func (m *MockTokenManager) SetTokenValid(valid bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TokenValid = valid
}

// SetGetTokenError define o erro que será retornado em GetGridToken
func (m *MockTokenManager) SetGetTokenError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GetTokenError = err
}

// SetRefreshError define o erro que será retornado em RefreshToken
func (m *MockTokenManager) SetRefreshError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RefreshError = err
}

// Verificação de chamadas

// GetGetTokenCalls retorna o histórico de chamadas para GetGridToken
func (m *MockTokenManager) GetGetTokenCalls() []GetTokenCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]GetTokenCall, len(m.GetTokenCalls))
	copy(calls, m.GetTokenCalls)
	return calls
}

// GetValidateCalls retorna o histórico de chamadas para ValidateToken
func (m *MockTokenManager) GetValidateCalls() []ValidateCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]ValidateCall, len(m.ValidateCalls))
	copy(calls, m.ValidateCalls)
	return calls
}

// GetRefreshCalls retorna o histórico de chamadas para RefreshToken
func (m *MockTokenManager) GetRefreshCalls() []RefreshCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]RefreshCall, len(m.RefreshCalls))
	copy(calls, m.RefreshCalls)
	return calls
}

// GetCallCounts retorna os contadores de chamadas
func (m *MockTokenManager) GetCallCounts() (getToken, validate, refresh int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.GetTokenCallCount, m.ValidateCallCount, m.RefreshCallCount
}

// Reset reseta o mock para o estado inicial
func (m *MockTokenManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.GridToken = "mock-grid-token-12345"
	m.TokenValid = true
	m.RefreshError = nil
	m.GetTokenError = nil
	m.GetTokenCalls = make([]GetTokenCall, 0)
	m.ValidateCalls = make([]ValidateCall, 0)
	m.RefreshCalls = make([]RefreshCall, 0)
	m.GetTokenCallCount = 0
	m.ValidateCallCount = 0
	m.RefreshCallCount = 0
}

// Verificação de comportamento

// WasGetTokenCalled verifica se GetGridToken foi chamado
func (m *MockTokenManager) WasGetTokenCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.GetTokenCallCount > 0
}

// WasValidateCalled verifica se ValidateToken foi chamado
func (m *MockTokenManager) WasValidateCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.ValidateCallCount > 0
}

// WasRefreshCalled verifica se RefreshToken foi chamado
func (m *MockTokenManager) WasRefreshCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.RefreshCallCount > 0
}

// WasValidateCalledWith verifica se ValidateToken foi chamado com um token específico
func (m *MockTokenManager) WasValidateCalledWith(token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, call := range m.ValidateCalls {
		if call.Token == token {
			return true
		}
	}
	return false
}

// GetLastGetTokenCall retorna a última chamada para GetGridToken
func (m *MockTokenManager) GetLastGetTokenCall() *GetTokenCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.GetTokenCalls) == 0 {
		return nil
	}

	lastCall := m.GetTokenCalls[len(m.GetTokenCalls)-1]
	return &lastCall
}

// GetLastValidateCall retorna a última chamada para ValidateToken
func (m *MockTokenManager) GetLastValidateCall() *ValidateCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.ValidateCalls) == 0 {
		return nil
	}

	lastCall := m.ValidateCalls[len(m.ValidateCalls)-1]
	return &lastCall
}

// GetLastRefreshCall retorna a última chamada para RefreshToken
func (m *MockTokenManager) GetLastRefreshCall() *RefreshCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.RefreshCalls) == 0 {
		return nil
	}

	lastCall := m.RefreshCalls[len(m.RefreshCalls)-1]
	return &lastCall
}

// Utilitários

// getCurrentTimestamp retorna o timestamp atual em nanosegundos
func getCurrentTimestamp() int64 {
	return int64(1234567890123456789) // Mock timestamp para testes determinísticos
}

// MockTokenManagerFactory cria uma factory para mocks do TokenManager
type MockTokenManagerFactory struct {
	defaultToken string
	defaultValid bool
}

// NewMockTokenManagerFactory cria uma nova factory
func NewMockTokenManagerFactory() *MockTokenManagerFactory {
	return &MockTokenManagerFactory{
		defaultToken: "mock-grid-token-12345",
		defaultValid: true,
	}
}

// SetDefaultToken define o token padrão
func (f *MockTokenManagerFactory) SetDefaultToken(token string) {
	f.defaultToken = token
}

// SetDefaultValid define se o token padrão deve ser válido
func (f *MockTokenManagerFactory) SetDefaultValid(valid bool) {
	f.defaultValid = valid
}

// Create cria um novo mock com as configurações padrão
func (f *MockTokenManagerFactory) Create() *MockTokenManager {
	mock := NewMockTokenManager()
	mock.SetGridToken(f.defaultToken)
	mock.SetTokenValid(f.defaultValid)
	return mock
}

// CreateWithError cria um mock que retorna erro em GetGridToken
func (f *MockTokenManagerFactory) CreateWithError(err error) *MockTokenManager {
	mock := f.Create()
	mock.SetGetTokenError(err)
	return mock
}

// CreateWithRefreshError cria um mock que retorna erro em RefreshToken
func (f *MockTokenManagerFactory) CreateWithRefreshError(err error) *MockTokenManager {
	mock := f.Create()
	mock.SetRefreshError(err)
	return mock
}

// CreateInvalidToken cria um mock com token inválido
func (f *MockTokenManagerFactory) CreateInvalidToken() *MockTokenManager {
	mock := f.Create()
	mock.SetTokenValid(false)
	return mock
}
