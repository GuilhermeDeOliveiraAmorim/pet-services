# Plano do Módulo de Adoção

## Visão Geral

Esta proposta trata adoção como um novo contexto de domínio dentro do Pet Services, aproveitando autenticação, fotos, localização, notificações e parte do cadastro de pets, mas sem acoplar o fluxo ao módulo atual de serviços.

Hoje o sistema está claramente organizado em torno de serviços, providers, requests e reviews. Por isso, a recomendação é que adoção seja uma vertical própria, com entidades, regras e telas específicas.

## Decisão Arquitetural

### Diretriz principal

1. Serviços continua como domínio independente.
2. Adoção entra como módulo paralelo.
3. O fluxo de adoção não reutiliza Request.
4. O fluxo de adoção não depende de Provider.
5. O cadastro do animal pode reaproveitar Pet como base, mas a publicação para adoção deve ter modelagem própria.

### Motivo

O fluxo atual de request foi desenhado para contratação de serviços entre owner e provider. Já adoção exige outra lógica: qualificação do interessado, histórico de análise, decisão, acompanhamento e fechamento do processo.

## Estratégia de Modelagem

### Decisão central

1. Pet continua sendo a entidade do animal.
2. AdoptionListing vira a entidade de publicação.
3. AdoptionApplication vira a entidade de candidatura.
4. AdoptionGuardianProfile vira a entidade do responsável pela adoção.
5. No MVP, não é necessário criar um novo tipo de usuário.

Essa abordagem evita mexer agora na autenticação, nos guards e na navegação baseada em Owner, Provider e Admin.

## Modelo de Domínio

### 1. AdoptionGuardianProfile

Representa ONG, protetor independente ou tutor que publica um pet para adoção.

Campos sugeridos:

1. id
2. user_id
3. display_name
4. guardian_type: ngo | independent | owner
5. document
6. phone
7. whatsapp
8. about
9. city_id
10. state_id
11. approval_status: pending | approved | rejected
12. approved_by
13. approved_at

### 2. AdoptionListing

Representa o anúncio de adoção vinculado a um pet.

Campos sugeridos:

1. id
2. pet_id
3. guardian_profile_id
4. title
5. description
6. adoption_reason
7. status: draft | published | paused | in_process | adopted | archived
8. sex: male | female
9. size: small | medium | large
10. age_group: puppy | adult | senior
11. vaccinated
12. neutered
13. dewormed
14. special_needs
15. good_with_children
16. good_with_dogs
17. good_with_cats
18. requires_house_screening
19. city_id
20. state_id
21. latitude
22. longitude
23. published_at
24. adopted_at

### 3. AdoptionApplication

Representa a candidatura de um interessado.

Campos sugeridos:

1. id
2. listing_id
3. applicant_user_id
4. status: submitted | under_review | interview | approved | rejected | withdrawn
5. motivation
6. housing_type
7. has_other_pets
8. pet_experience
9. family_members
10. agrees_home_visit
11. contact_phone
12. notes_internal
13. reviewed_by
14. reviewed_at

### 4. AdoptionApplicationEvent

Representa o histórico do processo da candidatura.

Campos sugeridos:

1. id
2. application_id
3. event_type
4. actor_user_id
5. payload_json
6. created_at

### 5. AdoptionFavorite

Opcional para uma evolução posterior do MVP.

Campos sugeridos:

1. id
2. listing_id
3. user_id

## O que reaproveitar do sistema atual

### Reaproveitar diretamente

1. Autenticação e sessão.
2. Cadastro de pet.
3. Upload de fotos.
4. Localização e referências geográficas.
5. Envio de e-mails e notificações.

### Não reaproveitar como base funcional

1. Request.
2. Provider.
3. Reviews.

Esses conceitos pertencem ao marketplace de serviços e não modelam corretamente o processo de adoção.

## Estrutura de Código Sugerida na API

Seguindo o padrão já existente da API, a sugestão é criar:

### Entidades

1. internal/entities/adoption_guardian_profile.go
2. internal/entities/adoption_listing.go
3. internal/entities/adoption_application.go
4. internal/entities/adoption_application_event.go

### Modelos GORM

1. internal/models/adoption_guardian_profile.go
2. internal/models/adoption_listing.go
3. internal/models/adoption_application.go
4. internal/models/adoption_application_event.go

### Repositórios

1. internal/repository_impl/adoption_guardian_profile_repository.go
2. internal/repository_impl/adoption_listing_repository.go
3. internal/repository_impl/adoption_application_repository.go

### Casos de uso

1. internal/usecases/create_adoption_guardian_profile.go
2. internal/usecases/create_adoption_listing.go
3. internal/usecases/publish_adoption_listing.go
4. internal/usecases/list_adoption_listings.go
5. internal/usecases/create_adoption_application.go
6. internal/usecases/review_adoption_application.go
7. internal/usecases/mark_adoption_listing_as_adopted.go

### Handlers

1. internal/handlers/adoption_guardian_handler.go
2. internal/handlers/adoption_listing_handler.go
3. internal/handlers/adoption_application_handler.go

### Rotas

1. internal/routes/adoption_routes.go

### Migrações e seeds

1. internal/database/migrations.go
2. internal/database/migrate.go

## Casos de Uso Necessários

Os casos de uso abaixo foram derivados diretamente das entidades já definidas para o domínio de adoção. A ideia é organizar a aplicação em torno dos fluxos reais do produto e evitar lacunas entre entidade, regra de negócio, endpoint e interface do front.

### MVP essencial

Esses são os casos de uso que formam o núcleo funcional do módulo.

1. Criar perfil de responsável pela adoção.
2. Consultar meu perfil de responsável.
3. Atualizar perfil de responsável.
4. Submeter perfil para aprovação.
5. Criar anúncio de adoção.
6. Consultar meu anúncio de adoção.
7. Listar meus anúncios de adoção.
8. Atualizar anúncio de adoção.
9. Publicar anúncio de adoção.
10. Pausar anúncio de adoção.
11. Marcar anúncio como em processo.
12. Marcar anúncio como adotado.
13. Arquivar anúncio.
14. Listar anúncios públicos.
15. Consultar detalhe de anúncio público.
16. Criar candidatura de adoção.
17. Listar minhas candidaturas.
18. Consultar detalhe da minha candidatura.
19. Desistir da candidatura.
20. Listar candidaturas de um anúncio.
21. Colocar candidatura em análise.
22. Avançar candidatura para entrevista.
23. Aprovar candidatura.
24. Rejeitar candidatura.
25. Registrar evento da candidatura.
26. Listar histórico de eventos da candidatura.

### Administração e moderação

Esses casos de uso entram logo depois do núcleo, porque são importantes para governança, auditoria e operação segura do módulo.

1. Listar perfis de responsável pendentes de aprovação.
2. Aprovar perfil de responsável.
3. Rejeitar perfil de responsável.
4. Listar anúncios para moderação.
5. Bloquear ou desativar anúncio, se necessário.
6. Consultar candidaturas de qualquer anúncio em contexto administrativo.
7. Consultar trilha de eventos para auditoria.

### Casos de uso de suporte

Esses casos de uso não são necessariamente a primeira entrega, mas costumam aparecer cedo porque complementam o fluxo principal.

1. Ver perfil público do responsável.
2. Validar elegibilidade para publicar anúncio.
3. Validar se usuário já se candidatou ao anúncio.
4. Validar se anúncio aceita candidatura no estado atual.
5. Notificar responsável sobre nova candidatura.
6. Notificar candidato sobre mudança de status.
7. Fechar automaticamente outras candidaturas quando um anúncio for marcado como adotado.

### Casos de uso de fase 2

Esses casos são relevantes, mas não precisam estar no primeiro ciclo de implementação.

1. Favoritar anúncio de adoção.
2. Remover favorito.
3. Listar favoritos.
4. Buscar anúncios por geolocalização e raio.
5. Anexar documentos do candidato.
6. Agendar visita ou entrevista.
7. Reabrir anúncio pausado ou arquivado com regras mais refinadas.
8. Exportar relatório operacional de adoção.

### Recorte recomendado para implementação imediata

Se a ideia for começar a pasta de use cases agora, a ordem mais pragmática é:

1. CreateAdoptionGuardianProfile.
2. GetMyAdoptionGuardianProfile.
3. UpdateAdoptionGuardianProfile.
4. CreateAdoptionListing.
5. UpdateAdoptionListing.
6. PublishAdoptionListing.
7. ListPublicAdoptionListings.
8. GetPublicAdoptionListing.
9. CreateAdoptionApplication.
10. ListMyAdoptionApplications.
11. ListAdoptionApplicationsByListing.
12. ReviewAdoptionApplication.
13. WithdrawAdoptionApplication.
14. MarkAdoptionListingAsAdopted.
15. ApproveAdoptionGuardianProfile.
16. RejectAdoptionGuardianProfile.

### Recomendação de agrupamento por fase

Para evitar abrir casos de uso demais cedo, o agrupamento recomendado é:

#### Fase 1

1. Perfis de responsável.
2. Criação e publicação de anúncio.
3. Catálogo público.
4. Candidatura.
5. Revisão de candidatura.
6. Aprovação administrativa de responsável.

#### Fase 2

1. Histórico detalhado.
2. Notificações.
3. Favoritos.
4. Moderação ampliada.

## Endpoints Sugeridos

### Públicos

1. GET /adoption/listings
   Lista anúncios publicados com filtros.

2. GET /adoption/listings/:id
   Retorna o detalhe do anúncio.

3. GET /adoption/guardians/:id
   Retorna o perfil público resumido do responsável.

### Autenticados

1. POST /adoption/applications
   Cria candidatura.

2. GET /adoption/applications/me
   Lista minhas candidaturas.

3. POST /adoption/applications/:id/withdraw
   Permite desistir da candidatura.

### Guardian ou Admin

1. POST /adoption/guardian-profile
   Cria ou atualiza perfil do responsável.

2. GET /adoption/guardian-profile/me
   Consulta o próprio perfil.

3. POST /adoption/listings
   Cria anúncio.

4. PATCH /adoption/listings/:id
   Edita anúncio.

5. POST /adoption/listings/:id/publish
   Publica anúncio.

6. POST /adoption/listings/:id/pause
   Pausa anúncio.

7. POST /adoption/listings/:id/mark-in-process
   Marca anúncio como em processo.

8. POST /adoption/listings/:id/mark-adopted
   Finaliza anúncio como adotado.

9. GET /adoption/listings/:id/applications
   Lista candidaturas do anúncio.

10. POST /adoption/applications/:id/review
    Avança ou altera o status da candidatura.

### Admin

1. GET /adoption/guardian-profiles
   Lista fila de aprovação.

2. POST /adoption/guardian-profiles/:id/approve
   Aprova perfil de responsável.

3. POST /adoption/guardian-profiles/:id/reject
   Rejeita perfil de responsável.

## Regras de Negócio do MVP

1. Um anúncio publicado só aceita candidatura se estiver com status published.
2. Um usuário não pode se candidatar duas vezes ao mesmo anúncio.
3. Quando o anúncio vai para in_process, novas candidaturas podem ser bloqueadas.
4. Quando o anúncio vai para adopted, ele fica fechado e auditável.
5. Apenas o responsável pelo anúncio ou admin pode alterar seus status.
6. Apenas usuários aprovados como guardian podem publicar anúncios.
7. Dados de adoção devem morar em AdoptionListing, e não em Pet.
8. Toda mudança relevante de candidatura deve gerar evento no histórico.

## Front-end Web

O módulo de adoção deve ter área pública e área autenticada.

### Rotas sugeridas

1. /adoption
   Catálogo público de pets para adoção.

2. /adoption/[id]
   Página de detalhe do anúncio.

3. /adoption/applications
   Área do interessado para acompanhar candidaturas.

4. /adoption/guardian
   Painel do responsável pela adoção.

5. /adoption/guardian/listings/new
   Criação de anúncio.

6. /adoption/guardian/listings/[id]/edit
   Edição de anúncio.

7. /admin/adoption/guardians
   Painel administrativo de aprovação.

### Componentes principais do MVP

1. Card de adoção.
2. Filtros por espécie, porte, sexo, idade e localização.
3. Galeria de fotos do animal.
4. Formulário de candidatura.
5. Painel de status da candidatura.
6. Tabela de candidaturas para o guardian.
7. Badge de status do anúncio.
8. Badge de status da candidatura.

## Mobile

Para a primeira entrega, a recomendação é não implementar o módulo completo no mobile.

Prioridade sugerida no mobile:

1. Catálogo público.
2. Detalhe do pet.
3. Minhas candidaturas.

O painel do guardian pode começar no web, onde o custo de implementação e operação será menor.

## Posicionamento no Produto

Na navegação principal, faz sentido adicionar um novo item:

1. Adote um Pet

Isso convive bem com os pilares já existentes:

1. Encontre Serviços
2. Seja um Parceiro

## Métricas Mínimas

Desde o MVP, vale acompanhar:

1. anúncios publicados
2. candidaturas enviadas
3. taxa de conversão por anúncio
4. tempo médio até adoção
5. origem do tráfego do catálogo de adoção

## Backlog MVP

### Fase 1: Fundamento de domínio e banco

1. Criar modelos e migrações de AdoptionGuardianProfile.
2. Criar modelos e migrações de AdoptionListing.
3. Criar modelos e migrações de AdoptionApplication.
4. Criar modelo de histórico AdoptionApplicationEvent.
5. Adicionar seeds básicos com pets para adoção.
6. Definir enums e validações de status.

Critério de pronto:

1. Banco sobe com migração limpa.
2. Seeds carregam corretamente.
3. Entidade de anúncio suporta filtros básicos.

### Fase 2: API pública de catálogo

1. Implementar listagem pública de anúncios.
2. Implementar detalhe do anúncio.
3. Adicionar filtros por espécie, porte, sexo, idade e cidade.
4. Adicionar paginação e total_items seguindo o padrão atual.
5. Expor dados resumidos do responsável.
6. Documentar endpoints no Swagger.

Critério de pronto:

1. Catálogo funciona anonimamente.
2. Detalhe carrega fotos e dados essenciais.
3. Filtros respondem bem.

### Fase 3: candidatura

1. Criar endpoint de candidatura.
2. Validar uma candidatura por usuário por anúncio.
3. Criar endpoint de minhas candidaturas.
4. Criar endpoint de desistência.
5. Disparar e-mail de candidatura recebida para o guardian.
6. Disparar e-mail de confirmação para o candidato.

Critério de pronto:

1. Usuário autenticado consegue se candidatar.
2. Histórico de candidatura fica rastreável.
3. Notificações são disparadas corretamente.

### Fase 4: painel do guardian

1. Criar cadastro de perfil de guardian.
2. Criar fluxo de criação de anúncio.
3. Criar fluxo de edição de anúncio.
4. Criar ação de publicar e pausar.
5. Criar listagem das candidaturas por anúncio.
6. Criar transições de status da candidatura.
7. Criar ação de marcar pet como adotado.

Critério de pronto:

1. Guardian gerencia anúncios sem apoio manual no banco.
2. Processo básico de triagem funciona.

### Fase 5: moderação e admin

1. Criar fila de perfis pendentes.
2. Criar aprovação e rejeição de guardians.
3. Criar bloqueio de publicação para perfis não aprovados.
4. Criar visibilidade administrativa de anúncios e candidaturas.
5. Criar trilha mínima de auditoria.

Critério de pronto:

1. Existe governança mínima do processo.
2. O risco operacional diminui.

### Fase 6: front web público

1. Criar página /adoption.
2. Criar página /adoption/[id].
3. Criar componente de filtro.
4. Criar card de anúncio.
5. Adicionar item de navegação no menu principal.
6. Adicionar chamada de ação de adoção na home, se fizer sentido.

Critério de pronto:

1. Descoberta pública está navegável.
2. O catálogo tem identidade própria no produto.

### Fase 7: front web autenticado

1. Criar formulário de candidatura.
2. Criar página de minhas candidaturas.
3. Criar painel do guardian.
4. Criar tela de criação e edição de anúncio.
5. Criar tela de gestão de candidaturas.
6. Criar painéis administrativos mínimos.

Critério de pronto:

1. Fluxo ponta a ponta funciona no navegador.

### Fase 8: qualidade

1. Adicionar testes dos casos de uso da API.
2. Adicionar testes de handlers críticos.
3. Validar transições inválidas de candidatura.
4. Testar listagem com filtros.
5. Adicionar testes básicos do front para candidatura e gestão de anúncio.

Critério de pronto:

1. Estados principais estão cobertos.
2. O risco de regressão diminui.

## Priorização Realista

### Prioridade 1

1. Catálogo público.
2. Detalhe do anúncio.
3. Candidatura autenticada.
4. Painel básico do guardian.
5. Aprovação manual por admin.

### Prioridade 2

1. Filtros geográficos mais refinados.
2. Histórico detalhado da candidatura.
3. Favoritos.
4. Notificações mais completas.

### Prioridade 3

1. Mobile.
2. Matching inteligente.
3. Documentos anexos.
4. Agenda de visita.

## Decisões para evitar retrabalho

1. Não criar novo user type no MVP.
2. Não usar Request para adoção.
3. Não misturar anúncio de adoção com cadastro de serviço.
4. Não colocar pagamento no primeiro ciclo.
5. Não automatizar entrevista ou visita domiciliar no MVP.
6. Manter moderação obrigatória para quem publica adoção.

## Ordem Recomendada de Implementação

1. API e banco.
2. Front público.
3. Fluxo autenticado de candidatura.
4. Painel do guardian.
5. Admin.
6. Mobile.

Essa ordem reduz risco e permite validar o interesse no módulo de adoção antes de expandir o investimento para todas as superfícies do produto.
