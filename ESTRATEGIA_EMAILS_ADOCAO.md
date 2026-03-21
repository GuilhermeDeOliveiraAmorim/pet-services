# Estratégia de Notificações por Email - Módulo de Adoção

**Data:** 21 de março de 2026  
**Status:** Planejamento de Implementação

---

## 1. Contexto Atual

### 1.1 Padrão Existente no Codebase

O system utiliza `EmailService` interface para notificações:

```go
// Interface padronizada
type EmailService interface {
    SendRequestCreatedEmail(...)
    SendRequestAcceptedEmail(...)
    SendRequestRejectedEmail(...)
    // ... mais métodos
}

// Injeção em factories
func NewRequestFactory(
    db *gorm.DB,
    storageService storage.ObjectStorage,
    mailService mail.EmailService,  // ← Injetado aqui
    logger logging.LoggerInterface,
) *RequestFactory
```

### 1.2 Implementação em UseCases

Exemplo de uso em `AddRequestUseCase`:

```go
// 1. Enviar ao provider
if err := uc.emailService.SendRequestCreatedEmail(
    providerUser.Login.Email,
    provider.BusinessName,
    user.Name,
    pet.Name,
    service.Name,
    requestEntity.ID,
); err != nil {
    return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar email", err)
}

// 2. Confirmação ao owner
if err := uc.emailService.SendRequestCreatedOwnerConfirmationEmail(
    user.Login.Email,
    user.Name,
    provider.BusinessName,
    pet.Name,
    service.Name,
    requestEntity.ID,
); err != nil {
    return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar email", err)
}
```

---

## 2. Análise de Pontos de Notificação por UC

### 2.1 Guardian Profile (5 UCs)

#### UC1: CreateAdoptionGuardianProfile
**Quando:** Perfil criado com sucesso  
**Destinatários:** User  
**Tipo:** Confirmação + Instrução

```
📧 PARA: user@email.com
Assunto: "Seu perfil como responsável foi criado com sucesso!"

Corpo:
  ✓ Perfil criado e aguardando aprovação do administrador
  → Link para editar perfil
  → Próximos passos
```

**Necessário adicionar:** `SendAdoptionGuardianProfileCreatedEmail(to, name string)`

---

#### UC2: GetMyAdoptionGuardianProfile
**Status:** Sem notificação (apenas consulta)

---

#### UC3: UpdateAdoptionGuardianProfile
**Quando:** Perfil atualizado  
**Destinatários:** User  
**Tipo:** Confirmação

```
📧 PARA: user@email.com
Assunto: "Seu perfil foi atualizado"

Corpo:
  ✓ Dados atualizados com sucesso
  → Se pendente: "Aguardando aprovação..."
  → Se aprovado: "Seu perfil aprovado continua ativo"
```

**Necessário adicionar:** `SendAdoptionGuardianProfileUpdatedEmail(to, name string, status string)`

---

#### UC4: ApproveAdoptionGuardianProfile ⭐ CRÍTICO
**Quando:** Admin aprova perfil  
**Destinatários:** User (Guardian)  
**Tipo:** Ação importante + Oportunidade

```
📧 PARA: guardian@email.com
Assunto: "🎉 Seu perfil foi aprovado! Comece a publicar anúncios"

Corpo:
  ✓ Perfil aprovado com sucesso
  ✓ Agora você pode publicar pets para adoção
  → CTA: "Publicar primeiro anúncio"
  → Link para dashboard
  → Dúvidas? Contate suporte
```

**Necessário adicionar:** `SendAdoptionGuardianProfileApprovedEmail(to, name string)`

---

#### UC5: RejectAdoptionGuardianProfile ⭐ CRÍTICO
**Quando:** Admin rejeita perfil  
**Destinatários:** User (Guardian)  
**Tipo:** Desapontamento + Próximos passos

```
📧 PARA: guardian@email.com
Assunto: "Informações sobre sua candidatura de perfil"

Corpo:
  ❌ Seu perfil foi recusado
  📝 Motivo: [reason opcional]
  → Recursos: Como melhorar
  → Suporte: support@petservices.com
  → Possibilidade: Reenviar após 30 dias
```

**Necessário adicionar:** `SendAdoptionGuardianProfileRejectedEmail(to, name, reason string)`

---

### 2.2 Adoption Listing (7 UCs)

#### UC1: CreateAdoptionListing
**Quando:** Anúncio criado com sucesso  
**Destinatários:** Guardian (Publisher)  
**Tipo:** Confirmação + Próximos passos

```
📧 PARA: guardian@email.com
Assunto: "Seu anúncio de adoção foi criado (ainda em rascunho)"

Corpo:
  ✓ Anúncio criado para ["Pet Name"]
  ⏳ Status: Rascunho
  → CTA: "Publicar agora" (PATCH /listings/{id}/active)
  → Link para editar
  → Dica: Adicione foto e descrição atrativa
```

**Necessário adicionar:** `SendAdoptionListingCreatedEmail(to, guardianName, petName string)`

---

#### UC2: UpdateAdoptionListing
**Quando:** Anúncio atualizado  
**Destinatários:** Guardian (Publisher)  
**Tipo:** Confirmação

```
📧 PARA: guardian@email.com
Assunto: "Seu anúncio foi atualizado"

Corpo:
  ✓ Alterações salvas para ["Pet Name"]
  → Ver anúncio publicado
  
[Se houver candidatos] 
  📬 Você tem 3 novas candidaturas! → Ver
```

**Necessário adicionar:** `SendAdoptionListingUpdatedEmail(to, guardianName, petName string)`

---

#### UC3: ChangeAdoptionListingStatus ⭐ IMPORTANTE
**Quando:** Status mudou (draft → active, active → paused, etc.)  
**Destinatários:** Guardian + possívelmente candidatos ativos  
**Tipo:** Confirmação / Notificação importante

```
📧 PARA: guardian@email.com
Assunto: "Status do seu anúncio foi alterado"

Corpo:
  ✓ Anúncio ["Pet Name"]
  Novo status: ["ATIVO" | "PAUSADO" | ...]
  
[Se ativado]
  🎉 Seu anúncio está visível para todos
  📊 Ver estatísticas
  
[Se pausado]
  ⏸️ Anúncio pausado e não aparecerá nas buscas
  → Reativar em [link]
```

**Quando status = PAUSED também notificar candidatos ativos:**
```
📧 PARA: candidate@email.com
Assunto: "Atualização: Anúncio pausado"

Corpo:
  ⏸️ O anúncio que você se candidatou foi pausado temporariamente
  → Você receberá atualizações quando reativar
```

**Necessário adicionar:**
- `SendAdoptionListingStatusChangedEmail(to, guardianName, petName, newStatus string)`
- `SendAdoptionListingPausedNotificationToApplicants(to, petName string)`

---

#### UC4: ListPublicAdoptionListings
**Status:** Sem notificação (apenas consulta)

---

#### UC5: GetPublicAdoptionListing
**Status:** Sem notificação (apenas consulta)

---

#### UC6: ListMyAdoptionListings
**Status:** Sem notificação (apenas consulta)

---

#### UC7: MarkAdoptionListingAsAdopted ⭐⭐ MUITO IMPORTANTE
**Quando:** Guardian marca como adotado  
**Destinatários:** Guardian, Candidato aprovado, Todos os outros candidatos  
**Tipo:** Celebração + Encerramento

```
📧 [1] PARA: guardian@email.com
Assunto: "🎉 Pet adotado com sucesso!"

Corpo:
  🎉 ["Pet Name"] foi marcado como ADOTADO!
  👋 Obrigado por usar Pet Services
  → Histórico de adoções
  → Próximos passos
  
[Link para feedback/satisfação]
```

```
📧 [2] PARA: accepted_candidate@email.com (se houver)
Assunto: "🎉 Parabéns! Seu novo companheiro"

Corpo:
  🎉 ["Pet Name"] foi adotado por você!
  💝 Bem-vindo à família!
  → Dicas de incorporação do pet
  → Suporte veterinário
  → Comunidade de adotantes
  
[Link para deixar foto/histórias]
```

```
📧 [3] PARA: other_candidates@email.com
Assunto: "Informação: Anúncio foi fechado"

Corpo:
  ❌ Infelizmente, ["Pet Name"] foi adotado
  💪 Não desista! Existem muitos pets esperando
  → Ver outros anúncios disponíveis
  → Novas publicações por email?
```

**Necessário adicionar:**
- `SendPetAdoptedGuardianEmail(to, guardianName, petName string)`
- `SendPetAdoptedApplicantEmail(to, applicantName, petName string)`
- `SendPetAdoptedRejectedApplicantsEmail(to, petName string)`

---

### 2.3 Adoption Application (5 UCs)

#### UC1: CreateAdoptionApplication ⭐ IMPORTANTE
**Quando:** Candidatura criada  
**Destinatários:** Candidato + Guardian  
**Tipo:** Confirmação + Notificação

```
📧 [1] PARA: candidate@email.com
Assunto: "✓ Sua candidatura foi enviada"

Corpo:
  ✓ Você se candidatou a ["Pet Name"]
  ⏳ Aguardando resposta do responsável
  → Acompanhar candidatura
  → Ver outros anúncios enquanto espera
```

```
📧 [2] PARA: guardian@email.com
Assunto: "📬 Nova candidatura para ["Pet Name"]"

Corpo:
  📬 Uma nova pessoa se candidatou a ["Pet Name"]!
  👤 Candidato: [Name]
  → CTA: "Ver candidatura" 
  → Link para revisar no dashboard
  🕐 Responda em até 7 dias para manter a candidatura ativa
```

**Necessário adicionar:**
- `SendAdoptionApplicationSubmittedEmail(to, applicantName, petName string)`
- `SendAdoptionApplicationReceivedGuardianEmail(to, guardianName, applicantName, petName string)`

---

#### UC2: ListMyAdoptionApplications
**Status:** Sem notificação (apenas consulta)

---

#### UC3: ListAdoptionApplicationsByListing
**Status:** Sem notificação (apenas consulta)

---

#### UC4: ReviewAdoptionApplication ⭐⭐ MUITO IMPORTANTE
**Quando:** Guardian revisa (aceita/rejeita)  
**Destinatários:** Candidato  
**Tipo:** Decisão importante

```
[SE ACEITO]

📧 PARA: candidate@email.com
Assunto: "🎉 Parabéns! Sua candidatura foi aprovada!"

Corpo:
  🎉 Excelentes notícias! Sua candidatura para ["Pet Name"] foi aceita!
  👋 O responsável entrará em contato em breve
  
  👤 Contato do responsável: [phone/email]
  
  📋 Próximos passos:
  1. Entre em contato para agendar encontro
  2. Conheça o pet pessoalmente
  3. Assine documentos de adoção
  4. Leve seu novo companheiro para casa! 🐾
```

```
[SE REJEITADO]

📧 PARA: candidate@email.com
Assunto: "Atualização sobre sua candidatura"

Corpo:
  ❌ Sua candidatura para ["Pet Name"] foi recusada
  
  Isso não significa que você não seja uma boa família!
  → Existem muitos outros pets esperando por você
  → Ver anúncios similares
  → Dicas para próximas candidaturas
  
  💌 Não desanime!
```

**Necessário adicionar:**
- `SendAdoptionApplicationApprovedEmail(to, applicantName, petName, guardianContact string)`
- `SendAdoptionApplicationRejectedEmail(to, applicantName, petName string)`

---

#### UC5: WithdrawAdoptionApplication
**Quando:** Candidato cancela inscrição  
**Destinatários:** Guardian (se candidata estava em andamento)  
**Tipo:** Notificação

```
📧 PARA: guardian@email.com
Assunto: "Candidato retirou sua candidatura"

Corpo:
  ⚠️ [Candidate Name] retirou sua candidatura de ["Pet Name"]
  
  → Ver outras candidaturas
  → Publicar em redes sociais
  → Contato de suporte
```

**Necessário adicionar:** `SendAdoptionApplicationWithdrawnGuardianEmail(to, guardianName, applicantName, petName string)`

---

## 3. Mapeamento de Métodos a Adicionar

### Na interface `EmailService`:

```go
type EmailService interface {
    // Existing methods...
    
    // === ADOPTION MODULE ===
    
    // Guardian Profile
    SendAdoptionGuardianProfileCreatedEmail(to, name string) error
    SendAdoptionGuardianProfileUpdatedEmail(to, name, status string) error
    SendAdoptionGuardianProfileApprovedEmail(to, name string) error
    SendAdoptionGuardianProfileRejectedEmail(to, name, reason string) error
    
    // Adoption Listing
    SendAdoptionListingCreatedEmail(to, guardianName, petName string) error
    SendAdoptionListingUpdatedEmail(to, guardianName, petName string) error
    SendAdoptionListingStatusChangedEmail(to, guardianName, petName, newStatus string) error
    SendAdoptionListingPausedNotificationToApplicants(to, petName string) error
    SendPetAdoptedGuardianEmail(to, guardianName, petName string) error
    SendPetAdoptedApplicantEmail(to, applicantName, petName string) error
    SendPetAdoptedRejectedApplicantsEmail(to, petName string) error
    
    // Adoption Application
    SendAdoptionApplicationSubmittedEmail(to, applicantName, petName string) error
    SendAdoptionApplicationReceivedGuardianEmail(to, guardianName, applicantName, petName string) error
    SendAdoptionApplicationApprovedEmail(to, applicantName, petName, guardianContact string) error
    SendAdoptionApplicationRejectedEmail(to, applicantName, petName string) error
    SendAdoptionApplicationWithdrawnGuardianEmail(to, guardianName, applicantName, petName string) error
}
```

**Total: 15 novos métodos de email**

---

## 4. Modificações Necessárias

### 4.1 Factories

#### `adoption_guardian_factory.go`

```go
type AdoptionGuardianFactory struct {
    CreateAdoptionGuardianProfile   *usecases.CreateAdoptionGuardianProfileUseCase
    GetMyAdoptionGuardianProfile    *usecases.GetMyAdoptionGuardianProfileUseCase
    UpdateAdoptionGuardianProfile   *usecases.UpdateAdoptionGuardianProfileUseCase
    ApproveAdoptionGuardianProfile  *usecases.ApproveAdoptionGuardianProfileUseCase
    RejectAdoptionGuardianProfile   *usecases.RejectAdoptionGuardianProfileUseCase
}

func NewAdoptionGuardianFactory(db *gorm.DB, mailService mail.EmailService, logger logging.LoggerInterface) *AdoptionGuardianFactory {
    // Adicionar mailService às UCs que precisam
    return &AdoptionGuardianFactory{
        CreateAdoptionGuardianProfile:  usecases.NewCreateAdoptionGuardianProfileUseCase(userRepo, guardianRepo, mailService, logger),
        UpdateAdoptionGuardianProfile:  usecases.NewUpdateAdoptionGuardianProfileUseCase(guardianRepo, mailService, logger),
        ApproveAdoptionGuardianProfile: usecases.NewApproveAdoptionGuardianProfileUseCase(guardianRepo, mailService, logger),
        RejectAdoptionGuardianProfile:  usecases.NewRejectAdoptionGuardianProfileUseCase(guardianRepo, mailService, logger),
    }
}
```

#### `adoption_listing_factory.go`

```go
func NewAdoptionListingFactory(db *gorm.DB, mailService mail.EmailService, logger logging.LoggerInterface) *AdoptionListingFactory {
    return &AdoptionListingFactory{
        CreateAdoptionListing:        usecases.NewCreateAdoptionListingUseCase(petRepo, listingRepo, mailService, logger),
        UpdateAdoptionListing:        usecases.NewUpdateAdoptionListingUseCase(listingRepo, mailService, logger),
        ChangeAdoptionListingStatus:  usecases.NewChangeAdoptionListingStatusUseCase(listingRepo, mailService, logger),
        MarkAdoptionListingAsAdopted: usecases.NewMarkAdoptionListingAsAdoptedUseCase(listingRepo, appRepo, userRepo, mailService, logger),
    }
}
```

#### `adoption_application_factory.go`

```go
func NewAdoptionApplicationFactory(db *gorm.DB, mailService mail.EmailService, logger logging.LoggerInterface) *AdoptionApplicationFactory {
    return &AdoptionApplicationFactory{
        CreateAdoptionApplication:  usecases.NewCreateAdoptionApplicationUseCase(..., mailService, logger),
        ReviewAdoptionApplication:  usecases.NewReviewAdoptionApplicationUseCase(..., mailService, logger),
        WithdrawAdoptionApplication: usecases.NewWithdrawAdoptionApplicationUseCase(..., mailService, logger),
    }
}
```

### 4.2 Handlers Factory

#### `handlers.go`

```go
func NewHandlerFactory(...) *HandlerFactory {
    // ...
    adoptionGuardianFactory := factories.NewAdoptionGuardianFactory(inputFactory.DB, mailService, logger)
    adoptionListingFactory := factories.NewAdoptionListingFactory(inputFactory.DB, mailService, logger)
    adoptionApplicationFactory := factories.NewAdoptionApplicationFactory(inputFactory.DB, mailService, logger)
    // ...
}
```

### 4.3 Use Cases

Adicionar email em 9 dos 17 UCs:

```
✅ CreateAdoptionGuardianProfile
✅ UpdateAdoptionGuardianProfile
✅ ApproveAdoptionGuardianProfile
✅ RejectAdoptionGuardianProfile
✅ CreateAdoptionListing
✅ UpdateAdoptionListing
✅ ChangeAdoptionListingStatus
✅ MarkAdoptionListingAsAdopted
✅ CreateAdoptionApplication
✅ ReviewAdoptionApplication
✅ WithdrawAdoptionApplication
```

---

## 5. Priorização de Implementação

### Fase 1 (Sprint atual) - CRÍTICO 🔴
Implementar notificações que são essenciais para fluxo de negócio:

1. `SendAdoptionGuardianProfileApprovedEmail` ← Guardian recebe permissão
2. `SendAdoptionGuardianProfileRejectedEmail` ← Guardian recebe feedback
3. `SendAdoptionApplicationSubmittedEmail` ← Candidato confirmação
4. `SendAdoptionApplicationReceivedGuardianEmail` ← Guardian nova candidatura
5. `SendAdoptionApplicationApprovedEmail` ← Candidato ACEITO
6. `SendAdoptionApplicationRejectedEmail` ← Candidato REJEITADO
7. `SendPetAdoptedGuardianEmail` ← Guardian celebração
8. `SendPetAdoptedApplicantEmail` ← Candidato celebração

**Impacto:** Estes 8 emails fecham os fluxos críticos. User não sabe o que fazer sem eles.

---

### Fase 2 (Sprint seguinte) - IMPORTANTE 🟡
Melhorar experiência com notificações de suporte:

1. `SendAdoptionGuardianProfileCreatedEmail`
2. `SendAdoptionListingCreatedEmail`
3. `SendAdoptionListingUpdatedEmail`
4. `SendAdoptionListingStatusChangedEmail`
5. `SendAdoptionListingPausedNotificationToApplicants`
6. `SendPetAdoptedRejectedApplicantsEmail`
7. `SendAdoptionApplicationWithdrawnGuardianEmail`

**Impacto:** Melhor acompanhamento + feedback contínuo.

---

### Fase 3 (Backlog) - ENHANCEMENT 🟢
Email avançados (analytics, recomendações):

- Digest semanal de novos anúncios
- Matchmaking automático
- Recordatório de candidaturas expiradas
- Feedback de satisfação pós-adoção

---

## 6. Implementação de Email SMTP

### Já existe:

```go
// internal/mail/email_service.go
type SMTPEmailService struct {
    host          string  // SMTP_HOST
    port          string  // SMTP_PORT
    user          string  // SMTP_USER
    password      string  // SMTP_PASSWORD
    from          string  // SMTP_FROM
    verifyBaseURL string
    resetBaseURL  string
}
```

### Adicionar para Adoção:

```go
type SMTPEmailService struct {
    // ... existing
    adoptionBaseURL string // NEW: URL base para links de adoção
}

// E.g.: https://petservices.com/adoption/listings/...
```

---

## 7. Exemplo de Implementação - UC: ApproveAdoptionGuardianProfile

### Passo 1: Adicionar interface

```go
// internal/mail/email_service.go
type EmailService interface {
    // ... existing methods
    SendAdoptionGuardianProfileApprovedEmail(to, name string) error
}
```

### Passo 2: Implementar no SMTP

```go
// internal/mail/email_service.go
func (s *SMTPEmailService) SendAdoptionGuardianProfileApprovedEmail(to, name string) error {
    subject := "🎉 Seu perfil foi aprovado! Comece a publicar anúncios"
    
    htmlBody := fmt.Sprintf(`
        <h1>Bem-vindo, %s!</h1>
        <p>Sua candidatura para ser responsável por adoção foi <strong>APROVADA</strong>!</p>
        <p>Agora você pode:</p>
        <ul>
            <li>Publicar pets para adoção</li>
            <li>Gerenciar candidaturas</li>
            <li>Acompanhar adoções</li>
        </ul>
        <p><a href="%s/adoption/listings/new" style="background: #4CAF50; color: white; padding: 10px 20px; text-decoration: none;">
            Publicar Primeiro Anúncio
        </a></p>
    `, html.EscapeString(name), s.adoptionBaseURL)
    
    return s.send(to, subject, htmlBody)
}

func (s *SMTPEmailService) send(to, subject, htmlBody string) error {
    // Use existing SMTP configuration
    auth := smtp.PlainAuth("", s.user, s.password, s.host)
    // ... send logic
}
```

### Passo 3: Adicionar ao UC

```go
// internal/usecases/approve_adoption_guardian_profile.go

type ApproveAdoptionGuardianProfileUseCase struct {
    guardianProfileRepo entities.AdoptionGuardianProfileRepository
    userRepository      entities.UserRepository  // NEW
    emailService        mail.EmailService        // NEW
    logger              logging.LoggerInterface
}

func (u *ApproveAdoptionGuardianProfileUseCase) Execute(ctx context.Context, input ApproveAdoptionGuardianProfileInput) (*ApproveAdoptionGuardianProfileOutput, []exceptions.ProblemDetails) {
    // ... existing logic ...
    
    // Update profile
    profile.Approve(input.ApprovedBy)
    if err := u.guardianProfileRepo.Update(profile); err != nil {
        // ... error handling
    }
    
    // NEW: Get user email
    user, err := u.userRepository.FindByID(profile.UserID)
    if err != nil {
        u.logger.LogError(ctx, "ApproveAdoptionGuardianProfileUseCase", "Erro ao buscar usuário", err)
        // ⚠️ Decide: return error or continue without email?
        // Recomendação: log error but continue (graceful degradation)
    } else {
        // NEW: Send approval email
        if err := u.emailService.SendAdoptionGuardianProfileApprovedEmail(user.Login.Email, user.Name); err != nil {
            u.logger.LogError(ctx, "ApproveAdoptionGuardianProfileUseCase", "Erro ao enviar email", err)
            // ⚠️ Don't fail the approval if email fails
            // Just log it for monitoring
        }
    }
    
    u.logger.LogInfo(ctx, "ApproveAdoptionGuardianProfileUseCase", "Perfil "+input.ProfileID+" aprovado")
    
    return &ApproveAdoptionGuardianProfileOutput{
        ID:             profile.ID,
        ApprovalStatus: profile.ApprovalStatus,
    }, nil
}
```

### Passo 4: Adicionar à Factory

```go
// internal/factories/adoption_guardian_factory.go

func NewAdoptionGuardianFactory(db *gorm.DB, mailService mail.EmailService, logger logging.LoggerInterface) *AdoptionGuardianFactory {
    userRepo := repository_impl.NewUserRepository(db)
    guardianProfileRepo := repository_impl.NewAdoptionGuardianProfileRepository(db)
    
    return &AdoptionGuardianFactory{
        // ... existing
        ApproveAdoptionGuardianProfile: usecases.NewApproveAdoptionGuardianProfileUseCase(
            guardianProfileRepo,
            userRepo,       // NEW
            mailService,    // NEW
            logger,
        ),
    }
}
```

---

## 8. Considerações Técnicas

### 8.1 Tratamento de Erro de Email

**Cenário 1: Email falha, deve a operação falhar?**

❌ **NÃO.** A adoção foi aprovada. Usuario pode reenviar email depois.

```go
// ✅ Padrão recomendado
if err := uc.emailService.SendEmail(...); err != nil {
    uc.logger.LogError(ctx, "UCName", "Email send failed", err)
    // Continue com sucesso - não retorna erro
}
```

**Cenário 2: Email está desabilitado (mock durante testes)**

```go
// Interface permite mock
type MockEmailService struct{}
func (m *MockEmailService) SendAdoptionGuardianProfileApprovedEmail(to, name string) error {
    return nil // Silently succeed in tests
}
```

### 8.2 Dados Necessários no Email

Para enviar email, UC precisa:

1. **Email do usuário** → Repository `FindByID(userID)`
2. **Nome do usuário** → Repository
3. **Detalhes contexto** → Já tem (pet name, listing title, etc.)

**Impacto em UC:** Ajuste de queries para trazer dados necessários.

---

### 8.3 URLs no Email

Links de ação devem apontar para frontend:

```
❌ Errado:
href="/adoption/listings/new"

✅ Correto:
href="https://petservices.com/adoption/listings/new"

📝 Configuração:
- Adicionar `adoptionBaseURL` em SMTPEmailService
- Carregar de env: `ADOPTION_BY_EMAIL_BASE_URL`
```

### 8.4 Localização (i18n)

Emails atualmente são em português. Para escalabilidade:

```go
// Future: suportar múltiplos idiomas
SendAdoptionGuardianProfileApprovedEmail(to, name, locale string) error
// locale = "pt-BR" | "en-US" | "es-ES"
```

---

## 9. Roadmap de Implementação

```
WEEK 1:
├─ Adicionar 15 métodos à interface EmailService
├─ Implementar 8 emails (Fase 1 - CRÍTICO)
├─ Atualizar ApproveAdoptionGuardianProfile UC (exemplo completo)
└─ Testes unitários dos novos UCs

WEEK 2:
├─ Adicionar email para 3 outros UCs críticos
├─ Testes de integração (E2E com SMTP mock)
└─ Documentation

WEEK 3:
├─ Implementar 7 emails (Fase 2 - IMPORTANTE)
├─ E2E testing em staging
└─ Preparar para produção

WEEK 4:
├─ Monitoring e alertas de falha de email
├─ Retry logic para emails falhados
└─ Deploy em produção
```

---

## 10. Conclusão

O módulo de adoção depende de **notificações por email** para:
- ✅ Guiar users através dos fluxos
- ✅ Prover feedback de ações importantes
- ✅ Criar senso de comunidade
- ✅ Manter engagement

**Recomendação:** Implementar Fase 1 (8 emails críticos) antes de qualquer release em produção.

**Esforço estimado:** 3-4 dias de desenvolvimento (experiência anterior com EmailService)

