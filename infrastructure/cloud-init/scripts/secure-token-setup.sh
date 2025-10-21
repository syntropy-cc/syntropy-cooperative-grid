#!/bin/bash
set -e

# Script executado durante boot do nó
ENCRYPTED_TOKEN="{{ .GridTokenEncrypted }}"
NODE_ID="{{ .NodeID }}"

# Verificar se as variáveis foram substituídas corretamente
if [ "$ENCRYPTED_TOKEN" = "{{ .GridTokenEncrypted }}" ] || [ "$NODE_ID" = "{{ .NodeID }}" ]; then
    echo "ERROR: Template variables not replaced correctly"
    exit 1
fi

# Criar diretório de configuração
mkdir -p /opt/syntropy/config
mkdir -p /opt/syntropy/bin

# Criar binário para descriptografia do token
cat > /opt/syntropy/bin/decrypt-token << 'EOF'
#!/bin/bash
set -e

# Ler variáveis de ambiente
ENCRYPTED_TOKEN="${ENCRYPTED_TOKEN}"
NODE_ID="${NODE_ID}"

if [ -z "$ENCRYPTED_TOKEN" ] || [ -z "$NODE_ID" ]; then
    echo "ERROR: Missing environment variables"
    exit 1
fi

# Criar script Python para descriptografia (mais simples que Go para o nó)
cat > /tmp/decrypt_token.py << 'PYEOF'
import os
import base64
import hashlib
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend
import pbkdf2

def decrypt_token(encrypted_token_b64, node_id):
    # Derivar chave do NodeID (mesma lógica do cloud-init generator)
    node_seed = hashlib.sha256(node_id.encode()).digest()
    node_key = pbkdf2.PBKDF2(node_seed, b'syntropy-grid', 50000, 32, hashlib.sha256).read()
    
    # Decodificar token criptografado
    ciphertext = base64.b64decode(encrypted_token_b64)
    
    # Extrair nonce e ciphertext
    nonce = ciphertext[:12]  # GCM nonce size
    encrypted_data = ciphertext[12:]
    
    # Descriptografar
    cipher = Cipher(algorithms.AES(node_key), modes.GCM(nonce), backend=default_backend())
    decryptor = cipher.decryptor()
    plaintext = decryptor.update(encrypted_data) + decryptor.finalize()
    
    return plaintext.decode()

if __name__ == "__main__":
    encrypted_token = os.environ.get('ENCRYPTED_TOKEN')
    node_id = os.environ.get('NODE_ID')
    
    if not encrypted_token or not node_id:
        print("ERROR: Missing environment variables")
        exit(1)
    
    try:
        token = decrypt_token(encrypted_token, node_id)
        
        # Salvar token com permissões restritivas
        with open('/opt/syntropy/config/grid_token', 'w') as f:
            f.write(token)
        
        os.chmod('/opt/syntropy/config/grid_token', 0o400)
        
        print("Token decrypted and stored successfully")
        
    except Exception as e:
        print(f"ERROR: Failed to decrypt token: {e}")
        exit(1)
PYEOF

# Executar descriptografia
python3 /tmp/decrypt_token.py

# Limpar arquivos temporários
rm -f /tmp/decrypt_token.py

# Verificar se o token foi salvo corretamente
if [ ! -f "/opt/syntropy/config/grid_token" ]; then
    echo "ERROR: Token file not created"
    exit 1
fi

echo "Grid token decryption completed successfully"
EOF

chmod +x /opt/syntropy/bin/decrypt-token

# Criar serviço systemd para descriptografia
cat > /etc/systemd/system/syntropy-token-decrypt.service << EOF
[Unit]
Description=Syntropy Grid Token Decryption
Before=syntropy-node.service
After=local-fs.target

[Service]
Type=oneshot
Environment=ENCRYPTED_TOKEN=$ENCRYPTED_TOKEN
Environment=NODE_ID=$NODE_ID
ExecStart=/opt/syntropy/bin/decrypt-token
RemainAfterExit=yes
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

systemctl enable syntropy-token-decrypt.service
systemctl start syntropy-token-decrypt.service

# Verificar se o serviço executou com sucesso
if ! systemctl is-failed syntropy-token-decrypt.service; then
    echo "Token decryption service completed successfully"
else
    echo "ERROR: Token decryption service failed"
    exit 1
fi
