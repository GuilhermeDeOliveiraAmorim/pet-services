# 📊 Comparação: API Backend vs Frontend

## ✅ Usecases Implementados no Frontend

### Auth (10/10) ✅
- ✅ login
- ✅ logout
- ✅ register-user
- ✅ refresh-token
- ✅ change-password
- ✅ request-password-reset
- ✅ reset-password
- ✅ verify-email
- ✅ resend-verification-email

### User (5/6) ⚠️
- ✅ get-profile
- ✅ update-user
- ✅ add-user-photo
- ✅ deactivate-user
- ✅ reactivate-user
- ❌ **delete-user** (falta)

### Pet (6/7) ⚠️
- ✅ add-pet
- ✅ get-pet
- ✅ list-pets
- ✅ update-pet
- ✅ delete-pet
- ✅ delete-pet-photo
- ✅ add-pet-photo

### Reference (3/3) ✅
- ✅ list-countries
- ✅ list-states
- ✅ list-cities

### Specie (1/1) ✅
- ✅ list-species

---

## ❌ Módulos FALTANDO no Frontend

### Provider (0/8) 🔴
- ❌ add-provider
- ❌ get-provider
- ❌ update-provider (implícito)
- ❌ delete-provider
- ❌ add-provider-photo
- ❌ delete-provider-photo
- ❌ list-providers (implícito)

### Service (0/11) 🔴
- ❌ add-service
- ❌ get-service
- ❌ update-service
- ❌ delete-service
- ❌ list-services
- ❌ search-services
- ❌ add-service-photo
- ❌ delete-service-photo
- ❌ add-service-category
- ❌ add-service-tag

### Request (0/5) 🔴
- ❌ add-request
- ❌ get-request
- ❌ list-requests
- ❌ accept-request
- ❌ reject-request
- ❌ complete-request

### Review (0/2) 🔴
- ❌ create-review
- ❌ list-reviews

### Category (0/2) 🔴
- ❌ list-categories
- ❌ create-category

### Tag (0/1) 🔴
- ❌ list-tags

### Admin (0/1) 🔴
- ❌ create-admin

### Health (0/2) 🔴
- ❌ health-check-api
- ❌ health-check-db

---

## 📈 Estatísticas

| Módulo | Implementado | Total API | % Completo |
|--------|-------------|-----------|------------|
| **Auth** | 10 | 10 | 100% ✅ |
| **User** | 5 | 6 | 83% ⚠️ |
| **Pet** | 6 | 7 | 86% ⚠️ |
| **Reference** | 3 | 3 | 100% ✅ |
| **Specie** | 1 | 1 | 100% ✅ |
| **Provider** | 0 | 8 | 0% 🔴 |
| **Service** | 0 | 11 | 0% 🔴 |
| **Request** | 0 | 5 | 0% 🔴 |
| **Review** | 0 | 2 | 0% 🔴 |
| **Category** | 0 | 2 | 0% 🔴 |
| **Tag** | 0 | 1 | 0% 🔴 |
| **Admin** | 0 | 1 | 0% 🔴 |
| **Health** | 0 | 2 | 0% 🔴 |
| **TOTAL** | **25** | **58** | **43%** |

---

## 🎯 Prioridades de Implementação

### 🔥 Alta Prioridade (Core Business)
1. **Provider Module** (8 usecases)
   - Essencial para prestadores de serviço se cadastrarem
   - Gerenciar perfil e fotos
   
2. **Service Module** (11 usecases)
   - Principal funcionalidade: listar/buscar serviços
   - Providers criarem/editarem serviços
   - Categorias e tags

3. **Request Module** (5 usecases)
   - Fluxo de agendamento
   - Owners solicitarem serviços
   - Providers aceitarem/rejeitarem

### 🟡 Média Prioridade
4. **Review Module** (2 usecases)
   - Avaliações e feedback
   - Construção de reputação

5. **Category & Tag** (3 usecases)
   - Organização de serviços
   - Filtros de busca

### 🟢 Baixa Prioridade
6. **Admin Module** (1 usecase)
   - Gerenciamento administrativo

7. **Health Module** (2 usecases)
   - Monitoramento (não essencial para usuários)

### ⚠️ Completar Módulos Existentes
- **User**: add delete-user
- **Pet**: add add-pet-photo

---

## 📝 Usecases Faltando (Detalhado)

### User
- `delete-user` - Hard delete de usuário

### Pet
- `add-pet-photo` - Upload de foto do pet

### Provider
- `add-provider` - Criar perfil de provider
- `get-provider` - Obter dados do provider
- `delete-provider` - Deletar provider
- `add-provider-photo` - Upload foto provider
- `delete-provider-photo` - Deletar foto provider

### Service
- `add-service` - Criar serviço
- `get-service` - Obter detalhes do serviço
- `update-service` - Atualizar serviço
- `delete-service` - Deletar serviço
- `list-services` - Listar todos os serviços
- `search-services` - Buscar serviços com filtros
- `add-service-photo` - Upload foto serviço
- `delete-service-photo` - Deletar foto serviço
- `add-service-category` - Associar categoria
- `add-service-tag` - Associar tag

### Request
- `add-request` - Criar requisição de serviço
- `get-request` - Obter requisição específica
- `list-requests` - Listar requisições
- `accept-request` - Provider aceitar
- `reject-request` - Provider rejeitar
- `complete-request` - Marcar como completo

### Review
- `create-review` - Criar avaliação
- `list-reviews` - Listar avaliações

### Category
- `list-categories` - Listar categorias
- `create-category` - Criar categoria (admin)

### Tag
- `list-tags` - Listar tags

### Admin
- `create-admin` - Criar usuário admin

### Health
- `health-check-api` - Status da API
- `health-check-db` - Status do banco

---

## 🚀 Próximos Passos Recomendados

1. **Completar User & Pet** (2 usecases)
   - Fechar módulos já iniciados
   
2. **Implementar Provider** (8 usecases)
   - Módulo crítico para o negócio
   
3. **Implementar Service** (11 usecases)
   - Core da aplicação
   
4. **Implementar Request** (5 usecases)
   - Fluxo de agendamento
   
5. **Implementar Review** (2 usecases)
   - Feedback e reputação

**Total de 28 usecases para completar 100% da API! 🎯**
