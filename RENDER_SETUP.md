# 🚀 Configuração no Render

## Variáveis de Ambiente

No Render dashboard, você precisa adicionar essas variáveis de ambiente para que a aplicação use Google Cloud Storage em produção.

### 📋 Variáveis Necessárias

| Variável                         | Tipo   | Valor             | Onde Obter                       |
| -------------------------------- | ------ | ----------------- | -------------------------------- |
| `IMAGE_BUCKET_NAME`              | String | `seu-bucket-nome` | Seu bucket GCS criado            |
| `GOOGLE_APPLICATION_CREDENTIALS` | String | Conteúdo do JSON  | Arquivo sa-key.json              |
| `DB_HOST`                        | String | Host do banco     | Render fornece                   |
| `DB_USER`                        | String | Usuário postgres  | Render fornece                   |
| `DB_PASSWORD`                    | String | Senha postgres    | Render fornece                   |
| `DB_PORT`                        | String | `5432`            | Render fornece                   |
| `DB_NAME`                        | String | `pet_services`    | Nome do banco                    |
| `JWT_ACCESS_SECRET`              | String | Chave aleatória   | Gerar com `openssl rand -hex 32` |
| `JWT_REFRESH_SECRET`             | String | Chave aleatória   | Gerar com `openssl rand -hex 32` |

### ⚙️ Setup no Render Dashboard

#### 1️⃣ Variáveis com Conteúdo Direto

```
IMAGE_BUCKET_NAME = seu-bucket-nome
MINIO_USE_SSL = true
RESET_PASSWORD_EXPIRATION_TIME = 3600
EMAIL_CHANGE_EXPIRATION_TIME = 86400
VERIFY_EMAIL_EXPIRATION_TIME = 86400
MAX_CHANGE_EMAIL_ATTEMPTS = 3
```

#### 2️⃣ Variável Crítica: `GOOGLE_APPLICATION_CREDENTIALS`

**⚠️ IMPORTANTE:** Cole o **conteúdo completo do arquivo JSON** (não o caminho)

**Como fazer:**

1. No seu computador, abra o arquivo `static-anchor-406816-8bae79e3e8b3.json`
2. Copie **TODO O CONTEÚDO** (é um JSON de ~2.3KB)
3. No Render dashboard:
   - Vá para **Environment** → **Environment Variables**
   - Clique em **Add Variable**
   - Key: `GOOGLE_APPLICATION_CREDENTIALS`
   - Value: Cole o JSON inteiro aqui
   - Salve

**Exemplo do JSON (primeiras linhas):**

```json
{
  "type": "service_account",
  "project_id": "static-anchor-406816",
  "private_key_id": "...",
  ...
}
```

#### 3️⃣ Variáveis do Banco de Dados

Render fornecerá automaticamente após criar o banco PostgreSQL:

- `DB_HOST` → host fornecido pelo Render
- `DB_USER` → `postgres` (padrão)
- `DB_PASSWORD` → senha gerada
- `DB_PORT` → `5432` (padrão PostgreSQL)
- `DB_NAME` → `pet_services`

#### 4️⃣ Variáveis de JWT

Gere chaves aleatórias seguras:

```bash
openssl rand -hex 32
# Resultado: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6

openssl rand -hex 32
# Resultado: x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5k4
```

Configure no Render:

```
JWT_ACCESS_SECRET = a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
JWT_REFRESH_SECRET = x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5k4
```

### ✅ Checklist

- [ ] Bucket GCS criado (`you-choose` ou outro nome)
- [ ] Service Account criado com Storage Admin
- [ ] Arquivo sa-key.json baixado
- [ ] Conteúdo JSON copiado para `GOOGLE_APPLICATION_CREDENTIALS`
- [ ] `IMAGE_BUCKET_NAME` configurado no Render
- [ ] JWT secrets gerados e configurados
- [ ] Banco de dados PostgreSQL criado no Render
- [ ] Credenciais do banco adicionadas
- [ ] Deploy realizado

### 🔄 Fluxo de Deploy

1. **Você faz push** para GitHub
2. **Render detecta** mudanças
3. **Render executa**:
   ```bash
   cd pet-services-api && go run ./cmd/migrate  # Cria tabelas
   go build -o api ./cmd/api  # Compila
   ./api  # Roda com variáveis de ambiente
   ```
4. **API funciona** com GCS automaticamente

### 🐛 Troubleshooting

**Erro: "bucket not found"**

```
Verificar:
- IMAGE_BUCKET_NAME está correto?
- GOOGLE_APPLICATION_CREDENTIALS foi colado inteiro?
```

**Erro: "permission denied"**

```
Verificar:
- Service Account tem Storage Admin?
- Arquivo sa-key.json é válido?
```

**Erro: "invalid credentials"**

```
Verificar:
- JSON foi colado completamente (sem quebras)?
- Não há espaços/caracteres extras?
```

### 📊 Estrutura Esperada no Render

```
Pet Services API
├── Environment Variables
│   ├── IMAGE_BUCKET_NAME = seu-bucket
│   ├── GOOGLE_APPLICATION_CREDENTIALS = { "type": "service_account", ... }
│   ├── DB_HOST = internal-db.render.com
│   ├── DB_USER = postgres
│   ├── DB_PASSWORD = ****
│   ├── DB_PORT = 5432
│   ├── DB_NAME = pet_services
│   ├── JWT_ACCESS_SECRET = ****
│   └── JWT_REFRESH_SECRET = ****
└── PostgreSQL Database
    └── pet-services-db (Ohio region)
```

### 🚀 Após Deploy

Testar a API:

```bash
# Health check
curl https://seu-render-app.onrender.com/api/health

# Upload de foto (requer autenticação)
curl -X POST https://seu-render-app.onrender.com/api/photos \
  -H "Authorization: Bearer seu-token-jwt" \
  -F "file=@foto.jpg"
```

### 📝 Notas Importantes

1. **Variáveis com `sync: false`** no `render.yaml` = você define no dashboard
2. **Variáveis com `value: "xxx"`** = vêm do arquivo yaml (fixas)
3. **Credenciais GCS** são críticas - não perca o arquivo JSON!
4. **Backup** do arquivo sa-key.json em local seguro
5. **Nunca** commit do arquivo JSON no Git!

### 🔐 Segurança

- ✅ Variáveis sensíveis **não** aparecem nos logs
- ✅ GCS é privado por padrão
- ✅ Service account tem permissões limitadas
- ✅ URLs assinadas expiram em 15 minutos

Tudo pronto! 🎉
