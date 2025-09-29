package crypto

import (
	"crypto/rand"
	"fmt"
	"io"
)

// EntropySource fornece uma fonte de entropia criptograficamente segura
type EntropySource struct {
	reader io.Reader
}

// NewEntropySource cria uma nova fonte de entropia
func NewEntropySource() *EntropySource {
	return &EntropySource{
		reader: rand.Reader,
	}
}

// GenerateRandomBytes gera bytes aleatórios criptograficamente seguros
func (es *EntropySource) GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("comprimento deve ser maior que zero")
	}

	if length > 1024*1024 { // 1MB máximo
		return nil, fmt.Errorf("comprimento muito grande: %d bytes", length)
	}

	bytes := make([]byte, length)
	n, err := es.reader.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar bytes aleatórios: %w", err)
	}

	if n != length {
		return nil, fmt.Errorf("gerou apenas %d bytes de %d solicitados", n, length)
	}

	return bytes, nil
}

// GeneratePrivateKey gera uma chave privada de tamanho específico
func (es *EntropySource) GeneratePrivateKey(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("tamanho da chave deve ser maior que zero")
	}

	return es.GenerateRandomBytes(size)
}

// GenerateSeed gera uma seed criptograficamente segura
func (es *EntropySource) GenerateSeed(length int) ([]byte, error) {
	if length < 16 {
		return nil, fmt.Errorf("seed deve ter pelo menos 16 bytes")
	}

	return es.GenerateRandomBytes(length)
}

// GenerateNonce gera um nonce criptograficamente seguro
func (es *EntropySource) GenerateNonce(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("tamanho do nonce deve ser maior que zero")
	}

	return es.GenerateRandomBytes(size)
}

// GenerateSalt gera um salt criptograficamente seguro
func (es *EntropySource) GenerateSalt(length int) ([]byte, error) {
	if length < 16 {
		return nil, fmt.Errorf("salt deve ter pelo menos 16 bytes")
	}

	return es.GenerateRandomBytes(length)
}

// TestEntropy testa a qualidade da fonte de entropia
func (es *EntropySource) TestEntropy() error {
	// Gerar uma amostra de bytes
	sample, err := es.GenerateRandomBytes(1024)
	if err != nil {
		return fmt.Errorf("falha ao gerar amostra para teste: %w", err)
	}

	// Teste básico de entropia - verificar distribuição de bytes
	byteCounts := make(map[byte]int)
	for _, b := range sample {
		byteCounts[b]++
	}

	// Verificar se não há bytes com contagem muito alta (baixa entropia)
	maxCount := 0
	for _, count := range byteCounts {
		if count > maxCount {
			maxCount = count
		}
	}

	// Se algum byte aparece mais de 5% das vezes, pode indicar baixa entropia
	threshold := len(sample) / 20
	if maxCount > threshold {
		return fmt.Errorf("possível baixa entropia detectada: byte mais frequente aparece %d vezes", maxCount)
	}

	return nil
}
