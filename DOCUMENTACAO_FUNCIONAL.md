# Documentação Funcional - Pet Services API

**Versão:** 1.0  
**Data:** Fevereiro de 2026  
**Público-alvo:** Product Owner e equipe não-técnica

---

## Índice

1. [Autenticação e Segurança](#1-autenticação-e-segurança)
2. [Usuários](#2-usuários)
3. [Pets](#3-pets)
4. [Provedores de Serviço](#4-provedores-de-serviço)
5. [Serviços](#5-serviços)
6. [Solicitações/Agendamentos](#6-solicitaçõesagendamentos)
7. [Avaliações](#7-avaliações)
8. [Categorias e Tags](#8-categorias-e-tags)
9. [Referências (Países, Estados, Cidades)](#9-referências-países-estados-cidades)
10. [Espécies](#10-espécies)

11. [Adoção & Guardian](#12-adoção--guardian)

---

## 1. Autenticação e Segurança

### 1.1 Cadastro de Novo Usuário

**O que faz:** Permite que uma nova pessoa se registre na plataforma.

**Quem pode usar:** Qualquer pessoa (sem necessidade de login).

**Informações necessárias:**

- Nome completo
- Email
- Senha
- Telefone (código do país, DDD e número)

**O que retorna:** Confirmação de que o cadastro foi criado com sucesso.

**Regras de negócio importantes:**

- Email deve ser único (não pode ter duas pessoas com mesmo email)
- Senha deve atender requisitos mínimos de segurança
- Após cadastro, usuário precisa verificar o email antes de fazer login
- Perfil inicial fica incompleto até que usuário informe tipo (dono de pet ou prestador) e endereço

---

### 1.2 Login (Entrar no Sistema)

**O que faz:** Permite que um usuário cadastrado entre na plataforma.

**Quem pode usar:** Usuários já cadastrados.

**Informações necessárias:**

- Email
- Senha

**O que retorna:**

- Dados do usuário logado
- Tokens de acesso (credenciais temporárias para usar o sistema)
- Tempo de validade dos tokens

**Regras de negócio importantes:**

- Email precisa estar verificado
- Conta precisa estar ativa (não desativada)
- Senha deve estar correta
- Login revoga todos os tokens anteriores (desloga de outros dispositivos)

---

### 1.3 Logout (Sair do Sistema)

**O que faz:** Permite que o usuário saia da plataforma com segurança.

**Quem pode usar:** Usuários logados.

**Informações necessárias:**

- Pode sair apenas do dispositivo atual OU
- Pode sair de todos os dispositivos de uma vez

**O que retorna:** Confirmação de logout.

**Regras de negócio importantes:**

- Revoga os tokens de acesso (credenciais)
- Usuário precisará fazer login novamente para acessar

---

### 1.4 Renovar Token de Acesso

**O que faz:** Renova as credenciais de acesso sem precisar fazer login novamente.

**Quem pode usar:** Usuários que já fizeram login.

**Informações necessárias:**

- Token de renovação (gerado no login)

**O que retorna:**

- Novos tokens de acesso
- Tempo de validade

**Regras de negócio importantes:**

- Token de renovação deve estar válido e não expirado
- Usuário deve estar com conta ativa
- Email deve estar verificado

---

### 1.5 Verificar Email

**O que faz:** Confirma que o email do usuário é válido através de um link/código enviado por email.

**Quem pode usar:** Usuários que acabaram de se cadastrar.

**Informações necessárias:**

- Código de verificação (recebido por email)

**O que retorna:** Confirmação de que email foi verificado.

**Regras de negócio importantes:**

- Código tem validade de 24 horas
- Após verificar email, usuário pode fazer login
- Se já verificado, retorna mensagem informando

---

### 1.6 Reenviar Email de Verificação

**O que faz:** Envia novamente o email de verificação caso usuário não tenha recebido.

**Quem pode usar:** Usuários que ainda não verificaram o email.

**Informações necessárias:**

- Email cadastrado

**O que retorna:** Confirmação de que novo email foi enviado.

**Regras de negócio importantes:**

- Revoga códigos anteriores
- Só envia se email ainda não foi verificado
- Novo código tem validade de 24 horas

---

### 1.7 Solicitar Redefinição de Senha

**O que faz:** Envia um link por email para o usuário redefinir sua senha esquecida.

**Quem pode usar:** Qualquer pessoa com email cadastrado.

**Informações necessárias:**

- Email cadastrado

**O que retorna:** Confirmação de que instruções foram enviadas por email.

**Regras de negócio importantes:**

- Revoga todos os códigos anteriores de redefinição
- Código tem validade de 1 hora
- Por segurança, sempre retorna sucesso mesmo se email não existir

---

### 1.8 Redefinir Senha

**O que faz:** Permite criar uma nova senha usando o código recebido por email.

**Quem pode usar:** Usuários que solicitaram redefinição de senha.

**Informações necessárias:**

- Código de redefinição (recebido por email)
- Nova senha

**O que retorna:** Confirmação de que senha foi alterada.

**Regras de negócio importantes:**

- Código deve estar válido e não expirado
- Nova senha deve atender requisitos de segurança
- Código é revogado após uso

---

### 1.9 Alterar Senha (Usuário Logado)

**O que faz:** Permite que usuário logado troque sua senha.

**Quem pode usar:** Usuários logados.

**Informações necessárias:**

- Senha atual
- Nova senha

**O que retorna:** Confirmação de que senha foi alterada.

**Regras de negócio importantes:**

- Senha atual deve estar correta
- Nova senha deve atender requisitos de segurança
- Nova senha não pode ser igual à atual

---

### 1.10 Criar Administrador do Sistema

**O que faz:** Cria uma nova conta de administrador da plataforma.

**Quem pode usar:** Apenas outros administradores.

**Informações necessárias:**

- Nome completo
- Email
- Senha
- Telefone
- Endereço completo

**O que retorna:** Confirmação de que administrador foi criado.

**Regras de negócio importantes:**

- Apenas administradores podem criar outros administradores
- Email deve ser único
- Administrador tem acesso completo ao sistema

---

## 2. Usuários

### 2.1 Ver Meu Perfil

**O que faz:** Mostra todos os dados do usuário logado.

**Quem pode usar:** Usuários logados.

**Informações necessárias:** Nenhuma (usa dados do usuário logado).

**O que retorna:**

- Nome
- Email
- Tipo de usuário (dono de pet, prestador ou admin)
- Telefone
- Endereço
- Foto de perfil
- Status da conta (ativa/inativa)
- Status do email (verificado/não verificado)

**Regras de negócio importantes:**

- Foto de perfil tem link temporário (válido por 15 minutos)

---

### 2.2 Ver Perfil de Outro Usuário

**O que faz:** Permite ver informações públicas de outros usuários.

**Quem pode usar:** Usuários logados.

**Informações necessárias:**

- ID do usuário que deseja visualizar

**O que retorna:** Dados públicos do usuário.

**Regras de negócio importantes:**

- Administradores podem ver qualquer perfil
- Donos de pets só podem ver perfis de prestadores
- Outros tipos não podem ver perfis alheios
- Usuário pode sempre ver seu próprio perfil

---

### 2.3 Listar Todos os Usuários

**O que faz:** Lista todos os usuários cadastrados na plataforma.

**Quem pode usar:** Administradores.

**Informações necessárias:**

- Número da página
- Quantidade de itens por página (máximo 100)

**O que retorna:**

- Lista de usuários
- Total de usuários
- Informações de paginação

**Regras de negócio importantes:**

- Paginação obrigatória
- Limite máximo de 100 usuários por página

---

### 2.4 Atualizar Meu Perfil

**O que faz:** Permite que usuário atualize seus dados cadastrais.

**Quem pode usar:** Usuários logados.

**Informações necessárias (opcionais):**

- Nome
- Tipo de usuário (dono de pet ou prestador)
- Telefone
- Endereço completo

**O que retorna:**

- Dados atualizados do usuário
- Confirmação de sucesso

**Regras de negócio importantes:**

- Não é possível tornar-se administrador por esta função
- Tipo deve ser "dono de pet" (owner) ou "prestador" (provider)
- Perfil é marcado como completo quando todos os dados obrigatórios estão preenchidos
- Coordenadas de localização são obrigatórias no endereço

---

### 2.5 Adicionar/Atualizar Foto de Perfil

**O que faz:** Permite enviar ou trocar a foto do perfil do usuário.

**Quem pode usar:** Usuários logados.

**Informações necessárias:**

- Arquivo de imagem

**O que retorna:**

- Link da foto enviada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas imagens são aceitas
- Substitui foto anterior se existir
- Foto antiga é removida do sistema

---

### 2.6 Desativar Minha Conta

**O que faz:** Desativa temporariamente a conta do usuário.

**Quem pode usar:** Usuários logados.

**Informações necessárias:** Nenhuma.

**O que retorna:** Confirmação de desativação.

**Regras de negócio importantes:**

- Revoga todos os tokens de acesso (faz logout de todos os dispositivos)
- Conta desativada não pode fazer login
- Dados não são excluídos (conta pode ser reativada)

---

### 2.7 Reativar Minha Conta

**O que faz:** Reativa uma conta que foi desativada anteriormente.

**Quem pode usar:** Usuários com contas desativadas.

**Informações necessárias:** Nenhuma (usuário precisa estar autenticado).

**O que retorna:** Confirmação de reativação.

**Regras de negócio importantes:**

- Usuário pode fazer login novamente após reativação
- Se conta já estiver ativa, apenas informa

---

### 2.8 Remover Usuário

**O que faz:** Desativa permanentemente um usuário (exclusão suave).

**Quem pode usar:** Administradores.

**Informações necessárias:**

- ID do usuário

**O que retorna:** Confirmação de remoção.

**Regras de negócio importantes:**

- Na verdade desativa o usuário (não exclui do banco)
- Usuário já inativo retorna erro
- Dados permanecem no sistema por questões de auditoria

---

### 2.9 Verificar se Email Existe

**O que faz:** Verifica se um email já está cadastrado no sistema.

**Quem pode usar:** Qualquer pessoa.

**Informações necessárias:**

- Email a verificar

**O que retorna:**

- Se email existe ou não
- Mensagem informativa

**Regras de negócio importantes:**

- Email deve ter formato válido
- Útil antes de tentar cadastrar

---

### 2.10 Verificar se Telefone Existe

**O que faz:** Verifica se um telefone já está cadastrado no sistema.

**Quem pode usar:** Qualquer pessoa.

**Informações necessárias:**

- Código do país
- DDD
- Número

**O que retorna:**

- Se telefone existe ou não
- Mensagem informativa

**Regras de negócio importantes:**

- Todos os componentes do telefone são obrigatórios
- Validação de formato

---

## 3. Pets

### 3.1 Cadastrar Novo Pet

**O que faz:** Permite que dono cadastre um novo pet.

**Quem pode usar:** Apenas donos de pets.

**Informações necessárias:**

- Nome do pet
- Espécie (ID da espécie cadastrada)
- Idade
- Peso
- Observações (opcional)

**O que retorna:**

- Dados completos do pet cadastrado
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas donos de pets podem cadastrar
- Espécie deve existir no sistema
- Idade e peso devem ser positivos

---

### 3.2 Listar Meus Pets

**O que faz:** Lista todos os pets cadastrados pelo dono.

**Quem pode usar:** Apenas donos de pets.

**Informações necessárias:**

- Número da página
- Quantidade por página

**O que retorna:**

- Lista de pets do usuário
- Fotos dos pets (se houver)
- Informações de paginação

**Regras de negócio importantes:**

- Apenas o próprio dono vê seus pets
- Máximo 100 pets por página
- Links de fotos válidos por 15 minutos

---

### 3.3 Ver Detalhes de um Pet

**O que faz:** Mostra todos os dados de um pet específico.

**Quem pode usar:** Apenas o dono do pet.

**Informações necessárias:**

- ID do pet

**O que retorna:**

- Nome, espécie, idade, peso
- Observações
- Fotos do pet

**Regras de negócio importantes:**

- Apenas dono pode ver seus pets
- Pet deve estar ativo
- Links de fotos válidos por 15 minutos

---

### 3.4 Atualizar Dados do Pet

**O que faz:** Permite alterar informações de um pet cadastrado.

**Quem pode usar:** Apenas o dono do pet.

**Informações necessárias (opcionais):**

- Nome
- Espécie
- Idade
- Peso
- Observações

**O que retorna:**

- Dados atualizados do pet
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas dono pode atualizar seus pets
- Pet deve estar ativo
- Valores não podem ser negativos
- Observações limitadas a 500 caracteres

---

### 3.5 Remover Pet

**O que faz:** Remove (desativa) um pet do cadastro.

**Quem pode usar:** Apenas o dono do pet.

**Informações necessárias:**

- ID do pet

**O que retorna:** Confirmação de remoção.

**Regras de negócio importantes:**

- Apenas dono pode remover seus pets
- Pet é desativado (não excluído)
- Pet já inativo retorna erro

---

### 3.6 Adicionar Foto ao Pet

**O que faz:** Permite enviar fotos do pet.

**Quem pode usar:** Apenas o dono do pet.

**Informações necessárias:**

- ID do pet
- Arquivo de imagem

**O que retorna:**

- Link da foto enviada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas imagens são aceitas
- Limite de 10 fotos por pet
- Apenas dono do pet pode adicionar fotos

---

### 3.7 Remover Foto do Pet

**O que faz:** Remove uma foto específica do pet.

**Quem pode usar:** Apenas o dono do pet.

**Informações necessárias:**

- ID do pet
- ID da foto

**O que retorna:** Confirmação de remoção.

**Regras de negócio importantes:**

- Apenas dono pode remover fotos dos seus pets
- Foto é excluída do armazenamento
- Foto deve pertencer ao pet especificado

---

## 4. Provedores de Serviço

### 4.1 Cadastrar Perfil de Prestador

**O que faz:** Cria o perfil profissional/comercial do prestador de serviços.

**Quem pode usar:** Usuários com tipo "prestador".

**Informações necessárias:**

- Nome do negócio/empresa
- Descrição dos serviços
- Faixa de preço (ex: "$$", "$$$")
- Endereço completo com coordenadas

**O que retorna:**

- Dados do prestador criado
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas usuários tipo "prestador" podem criar
- Cada usuário pode ter apenas um perfil de prestador
- Endereço deve incluir coordenadas para busca por localização
- Se já existe prestador para o usuário, retorna erro

---

### 4.2 Ver Perfil Completo do Prestador

**O que faz:** Mostra todas as informações de um prestador e seus serviços.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias:**

- ID do prestador
- Número da página de serviços
- Quantidade de serviços por página

**O que retorna:**

- Dados completos do prestador
- Lista de serviços oferecidos
- Fotos do prestador
- Avaliações e nota média
- Paginação dos serviços

**Regras de negócio importantes:**

- Prestador deve estar ativo
- Máximo 100 serviços por página
- Links de fotos válidos por 15 minutos
- Mostra apenas serviços ativos

---

### 4.3 Remover Perfil de Prestador

**O que faz:** Desativa o perfil de prestador.

**Quem pode usar:** Apenas o próprio prestador.

**Informações necessárias:**

- ID do prestador

**O que retorna:** Confirmação de remoção.

**Regras de negócio importantes:**

- Apenas o próprio prestador pode remover
- Prestador é desativado (não excluído)
- Prestador já inativo retorna erro
- Serviços associados permanecem mas inativos

---

### 4.4 Adicionar Foto ao Prestador

**O que faz:** Permite enviar fotos do estabelecimento/trabalho.

**Quem pode usar:** Apenas o próprio prestador.

**Informações necessárias:**

- Arquivo de imagem

**O que retorna:**

- Link da foto enviada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas imagens são aceitas
- Limite de 10 fotos por prestador
- Apenas o próprio prestador pode adicionar

---

### 4.5 Remover Foto do Prestador

**O que faz:** Remove uma foto específica do prestador.

**Quem pode usar:** Apenas o próprio prestador.

**Informações necessárias:**

- ID da foto

**O que retorna:** Confirmação de remoção.

**Regras de negócio importantes:**

- Foto é excluída do armazenamento
- Apenas o próprio prestador pode remover
- Foto deve pertencer ao prestador

---

## 5. Serviços

### 5.1 Cadastrar Novo Serviço

**O que faz:** Permite que prestador cadastre um serviço que oferece.

**Quem pode usar:** Apenas prestadores.

**Informações necessárias:**

- Nome do serviço
- Descrição
- Preço fixo OU faixa de preço (mínimo e máximo)
- Duração estimada em minutos

**O que retorna:**

- Dados completos do serviço
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas prestadores podem cadastrar
- Deve informar preço fixo OU faixa de preço (não ambos)
- Se usar faixa, preço mínimo deve ser menor que máximo
- Duração deve ser positiva
- Serviço é criado como ativo

---

### 5.2 Listar Serviços

**O que faz:** Lista serviços com filtros opcionais.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias (opcionais):**

- ID do prestador (para ver serviços de um prestador específico)
- ID da categoria
- ID da tag
- Preço mínimo
- Preço máximo
- Paginação

**O que retorna:**

- Lista de serviços
- Nome do prestador
- Fotos, categorias e tags de cada serviço
- Paginação

**Regras de negócio importantes:**

- Mostra apenas serviços ativos
- Máximo 100 serviços por página
- Links de fotos válidos por 15 minutos

---

### 5.3 Buscar/Pesquisar Serviços

**O que faz:** Busca avançada de serviços com múltiplos filtros.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias (opcionais):**

- Texto de busca (nome ou descrição)
- Categoria
- Tag
- Localização (latitude, longitude e raio em km)
- Faixa de preço
- Paginação

**O que retorna:**

- Lista de serviços que atendem aos critérios
- Distância do usuário (se informou localização)
- Dados do prestador
- Fotos do serviço

**Regras de negócio importantes:**

- Busca por proximidade usa raio padrão de 10km
- Combina todos os filtros informados (busca E filtros)
- Máximo 100 resultados por página
- Links de fotos válidos por 15 minutos

---

### 5.4 Ver Detalhes de um Serviço

**O que faz:** Mostra todas as informações de um serviço específico.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias:**

- ID do serviço

**O que retorna:**

- Nome, descrição, preços, duração
- Fotos do serviço
- Categorias e tags
- Dados do prestador

**Regras de negócio importantes:**

- Serviço deve estar ativo
- Links de fotos válidos por 15 minutos

---

### 5.5 Atualizar Serviço

**O que faz:** Permite que prestador altere dados de seu serviço.

**Quem pode usar:** Apenas o prestador dono do serviço.

**Informações necessárias (opcionais):**

- Nome
- Descrição
- Preços
- Duração

**O que retorna:**

- Dados atualizados do serviço
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas prestador dono pode atualizar
- Serviço deve estar ativo
- Validações de preço e duração mantêm-se
- Descrição limitada a 1000 caracteres

---

### 5.6 Remover Serviço

**O que faz:** Desativa um serviço oferecido pelo prestador.

**Quem pode usar:** Apenas o prestador dono do serviço.

**Informações necessárias:**

- ID do serviço

**O que retorna:** Confirmação de remoção.

**Regras de negócio importantes:**

- Apenas prestador dono pode remover
- Serviço é desativado (não excluído)
- Serviço já inativo retorna erro

---

### 5.7 Adicionar Foto ao Serviço

**O que faz:** Permite enviar fotos do serviço oferecido.

**Quem pode usar:** Apenas o prestador dono do serviço.

**Informações necessárias:**

- ID do serviço
- Arquivo de imagem

**O que retorna:**

- Link da foto enviada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas imagens são aceitas
- Limite de 10 fotos por serviço
- Serviço deve estar ativo

---

### 5.8 Remover Foto do Serviço

**O que faz:** Remove uma foto específica do serviço.

**Quem pode usar:** Apenas o prestador dono do serviço.

**Informações necessárias:**

- ID do serviço
- ID da foto

**O que retorna:** Confirmação de remoção.

**Regras de negócio importantes:**

- Foto é excluída do armazenamento
- Serviço deve estar ativo
- Foto deve pertencer ao serviço

---

### 5.9 Adicionar Categoria ao Serviço

**O que faz:** Vincula uma categoria ao serviço (ex: banho, tosa, veterinário).

**Quem pode usar:** Apenas o prestador dono do serviço.

**Informações necessárias:**

- ID do serviço
- ID da categoria

**O que retorna:**

- Categoria vinculada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Categoria deve existir e estar ativa
- Serviço deve estar ativo
- Não pode vincular categoria duplicada
- Apenas prestador dono pode adicionar

---

### 5.10 Adicionar Tag ao Serviço

**O que faz:** Adiciona palavras-chave/tags ao serviço para facilitar busca.

**Quem pode usar:** Apenas o prestador dono do serviço.

**Informações necessárias:**

- ID do serviço
- ID da tag OU nome da tag (cria se não existir)

**O que retorna:**

- Tag vinculada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Se tag não existir, cria automaticamente
- Tag deve estar ativa
- Serviço deve estar ativo
- Não pode vincular tag duplicada

---

## 6. Solicitações/Agendamentos

### 6.1 Criar Nova Solicitação

**O que faz:** Dono de pet solicita um serviço para seu pet.

**Quem pode usar:** Apenas donos de pets.

**Informações necessárias:**

- ID do serviço desejado
- ID do pet
- Observações (opcional)

**O que retorna:**

- Dados completos da solicitação
- Confirmação de envio

**Regras de negócio importantes:**

- Apenas donos de pets podem solicitar
- Pet deve pertencer ao dono
- Serviço e prestador devem estar ativos
- Não pode ter solicitação pendente duplicada (mesmo pet + mesmo serviço)
- Solicitação criada com status "pendente"

---

### 6.2 Listar Minhas Solicitações

**O que faz:** Lista solicitações do usuário (como dono ou como prestador).

**Quem pode usar:** Donos de pets e prestadores.

**Informações necessárias (opcionais):**

- Filtro por status (pendente, aceita, rejeitada, concluída)
- Paginação

**O que retorna:**

- Lista de solicitações
- Dados do pet
- Nome do serviço
- Nome do dono e do prestador
- Paginação

**Regras de negócio importantes:**

- Dono vê solicitações que fez
- Prestador vê solicitações recebidas
- Máximo 100 por página
- Links de fotos dos pets válidos por 15 minutos

---

### 6.3 Ver Detalhes de uma Solicitação

**O que faz:** Mostra todas as informações de uma solicitação específica.

**Quem pode usar:** Dono que fez a solicitação OU prestador que recebeu.

**Informações necessárias:**

- ID da solicitação

**O que retorna:**

- Dados completos da solicitação
- Informações do pet com fotos
- Dados do serviço
- Status atual
- Motivo de rejeição (se aplicável)

**Regras de negócio importantes:**

- Apenas participantes podem ver (dono ou prestador)
- Links de fotos válidos por 15 minutos

---

### 6.4 Aceitar Solicitação

**O que faz:** Prestador aceita uma solicitação recebida.

**Quem pode usar:** Apenas o prestador que recebeu a solicitação.

**Informações necessárias:**

- ID da solicitação

**O que retorna:**

- Dados atualizados da solicitação
- Confirmação de aceite

**Regras de negócio importantes:**

- Apenas prestador dono do serviço pode aceitar
- Solicitação deve estar com status "pendente"
- Status muda para "aceita"
- Dono será notificado (futuramente)

---

### 6.5 Rejeitar Solicitação

**O que faz:** Prestador recusa uma solicitação recebida.

**Quem pode usar:** Apenas o prestador que recebeu a solicitação.

**Informações necessárias:**

- ID da solicitação
- Motivo da rejeição (obrigatório)

**O que retorna:**

- Dados atualizados da solicitação
- Confirmação de rejeição

**Regras de negócio importantes:**

- Apenas prestador dono do serviço pode rejeitar
- Solicitação deve estar com status "pendente"
- Motivo é obrigatório (máximo 500 caracteres)
- Status muda para "rejeitada"

---

### 6.6 Concluir Solicitação

**O que faz:** Prestador marca solicitação como concluída após realizar o serviço.

**Quem pode usar:** Apenas o prestador que recebeu a solicitação.

**Informações necessárias:**

- ID da solicitação

**O que retorna:**

- Dados atualizados da solicitação
- Confirmação de conclusão

**Regras de negócio importantes:**

- Apenas prestador dono do serviço pode concluir
- Solicitação deve estar com status "aceita"
- Status muda para "concluída"
- Permite que dono faça avaliação

---

## 7. Avaliações

### 7.1 Criar Avaliação

**O que faz:** Dono de pet avalia um prestador após serviço concluído.

**Quem pode usar:** Apenas donos de pets.

**Informações necessárias:**

- ID do prestador
- Nota (0 a 5)
- Comentário (opcional)

**O que retorna:**

- Dados da avaliação criada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas donos podem avaliar
- Deve ter pelo menos uma solicitação concluída com o prestador
- Nota entre 0 e 5
- Comentário opcional
- Uma avaliação por solicitação concluída

---

### 7.2 Listar Avaliações

**O que faz:** Lista avaliações de um prestador ou feitas por um usuário.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias:**

- ID do prestador OU ID do usuário
- Paginação

**O que retorna:**

- Lista de avaliações
- Nota e comentário
- Nome do avaliador
- Data da avaliação
- Paginação

**Regras de negócio importantes:**

- Deve informar pelo menos um filtro (prestador ou usuário)
- Máximo 100 avaliações por página
- Avaliações ordenadas da mais recente para mais antiga

---

## 8. Categorias e Tags

### 8.1 Criar Nova Categoria

**O que faz:** Cria uma nova categoria de serviços (ex: Banho e Tosa, Veterinário).

**Quem pode usar:** Apenas administradores.

**Informações necessárias:**

- Nome da categoria

**O que retorna:**

- Dados da categoria criada
- Confirmação de sucesso

**Regras de negócio importantes:**

- Apenas administradores podem criar
- Nome deve ser único
- Categoria criada como ativa

---

### 8.2 Listar Categorias

**O que faz:** Lista todas as categorias disponíveis.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias (opcionais):**

- Filtro por nome
- Paginação

**O que retorna:**

- Lista de categorias ativas
- Total de categorias

**Regras de negócio importantes:**

- Mostra apenas categorias ativas
- Permite busca por nome

---

### 8.3 Listar Tags

**O que faz:** Lista todas as tags disponíveis.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias (opcionais):**

- Filtro por nome
- Paginação

**O que retorna:**

- Lista de tags ativas
- Total de tags

**Regras de negócio importantes:**

- Mostra apenas tags ativas
- Permite busca por nome
- Tags são criadas automaticamente ao adicionar em serviços

---

## 9. Referências (Países, Estados, Cidades)

### 9.1 Listar Países

**O que faz:** Lista todos os países disponíveis no sistema.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias:** Nenhuma.

**O que retorna:** Lista completa de países.

**Regras de negócio importantes:**

- Lista pré-cadastrada no sistema
- Usado para preenchimento de endereços

---

### 9.2 Listar Estados

**O que faz:** Lista todos os estados/províncias disponíveis.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias:** Nenhuma.

**O que retorna:** Lista completa de estados.

**Regras de negócio importantes:**

- Lista pré-cadastrada no sistema
- Usado para preenchimento de endereços

---

### 9.3 Listar Cidades

**O que faz:** Lista cidades, opcionalmente filtradas por estado.

**Quem pode usar:** Qualquer usuário.

**Informações necessárias (opcional):**

- ID do estado (para filtrar cidades daquele estado)

**O que retorna:** Lista de cidades.

**Regras de negócio importantes:**

- Se informar estado, retorna apenas cidades daquele estado
- Sem filtro, retorna todas as cidades
- Usado para preenchimento de endereços

---

## 10. Espécies

### 10.1 Listar Espécies de Pets

**O que faz:** Lista todas as espécies de pets disponíveis (cachorro, gato, etc).

**Quem pode usar:** Qualquer usuário.

**Informações necessárias:** Nenhuma.

**O que retorna:** Lista completa de espécies.

**Regras de negócio importantes:**

- Lista pré-cadastrada no sistema
- Usado ao cadastrar pets
- Inclui nome comum da espécie

---

## 12. Adoção & Guardian

### 12.1 Criar/Editar Perfil de Responsável (Guardian)

**O que faz:** Permite que ONGs, protetores ou tutores publiquem pets para adoção criando um perfil de responsável.

**Quem pode usar:** Usuários autenticados (não admin, não provider).

**Informações necessárias:**

- Nome de exibição
- Tipo de responsável (ONG, protetor, tutor)
- Documento (CNPJ/CPF)
- Telefone/Whatsapp
- Sobre (opcional)
- Cidade/Estado

**O que retorna:** Dados do perfil criado/atualizado, status de aprovação.

**Regras de negócio importantes:**

- Cada usuário pode ter apenas 1 perfil de responsável
- Perfil criado fica com status "pendente" até aprovação admin
- Admin pode aprovar/rejeitar perfis
- Apenas perfis aprovados podem publicar pets para adoção

---

### 12.2 Publicar/Listar Pets para Adoção

**O que faz:** Permite que responsáveis publiquem pets disponíveis para adoção.

**Quem pode usar:** Responsáveis (guardian) com perfil aprovado.

**Informações necessárias:**

- Nome do pet, espécie, idade, peso, descrição, fotos
- Cidade/Estado

**O que retorna:** Dados da listagem criada, status (ativo/adotado)

**Regras de negócio importantes:**

- Apenas guardian aprovado pode publicar
- Listagens podem ser editadas/removidas pelo próprio responsável
- Admin pode gerenciar todas as listagens

---

### 12.3 Candidatar-se à Adoção

**O que faz:** Permite que usuários interessados se candidatem para adotar um pet.

**Quem pode usar:** Usuários autenticados (não admin, não guardian do pet).

**Informações necessárias:**

- Mensagem de apresentação

**O que retorna:** Confirmação de candidatura, status (pendente)

**Regras de negócio importantes:**

- Não pode se candidatar ao próprio pet
- Não pode candidatar-se mais de uma vez para o mesmo pet
- Guardian recebe notificação de nova candidatura

---

### 12.4 Revisar/Alterar Status de Candidatura

**O que faz:** Permite que o responsável aprove ou rejeite candidaturas recebidas.

**Quem pode usar:** Guardian responsável pela listagem.

**Informações necessárias:**

- Aprovar ou rejeitar (com motivo)

**O que retorna:** Status atualizado da candidatura

**Regras de negócio importantes:**

- Apenas guardian pode aprovar/rejeitar candidaturas de seus pets
- Candidatura aprovada: contato do responsável é enviado ao adotante
- Candidatura rejeitada: notificação ao candidato

---

### 12.5 Marcar Pet como Adotado

**O que faz:** Permite que o responsável marque o pet como adotado, encerrando a listagem.

**Quem pode usar:** Guardian responsável pela listagem.

**Informações necessárias:**

- ID da listagem

**O que retorna:** Status atualizado da listagem (adotado)

**Regras de negócio importantes:**

- Apenas guardian pode marcar como adotado
- Todos os candidatos não aprovados recebem notificação de encerramento

---

### 12.6 Aprovação/Rejeição de Perfis (Admin)

**O que faz:** Permite que administradores aprovem ou rejeitem perfis de responsável (guardian).

**Quem pode usar:** Administradores

**Informações necessárias:**

- ID do perfil, decisão (aprovar/rejeitar), motivo (se rejeitar)

**O que retorna:** Status atualizado do perfil

**Regras de negócio importantes:**

- Apenas admin pode aprovar/rejeitar
- Guardian recebe notificação da decisão

---

### 12.7 Permissões e Notificações

- Apenas guardian aprovado pode publicar/adotar
- Admin pode gerenciar todos os perfis e listagens
- Notificações por e-mail para todos os eventos críticos (candidatura, aprovação, rejeição, adoção)

---

### 11.1 Verificar Status da API

**O que faz:** Verifica se a API está funcionando corretamente.

**Quem pode usar:** Qualquer pessoa (sem necessidade de login).

**Informações necessárias:** Nenhuma.

**O que retorna:** Status "ok" se sistema está funcionando.

**Regras de negócio importantes:**

- Endpoint público
- Usado para monitoramento

---

### 11.2 Verificar Status do Banco de Dados

**O que faz:** Verifica se a conexão com o banco de dados está funcionando.

**Quem pode usar:** Qualquer pessoa (sem necessidade de login).

**Informações necessárias:** Nenhuma.

**O que retorna:**

- Status "ok" se banco está acessível
- Status "error" se houver problemas

**Regras de negócio importantes:**

- Endpoint público
- Usado para monitoramento e diagnóstico

---

## Glossário de Termos

**Dono de pet (Owner):** Usuário que possui pets e busca serviços para eles.

**Prestador (Provider):** Usuário que oferece serviços relacionados a pets.

**Administrador (Admin):** Usuário com permissões especiais para gerenciar o sistema.

**Token:** Credencial temporária que permite acesso às funcionalidades após login.

**Desativar:** Tornar inativo, mas manter os dados no sistema (não é exclusão definitiva).

**Status Pendente:** Solicitação aguardando resposta do prestador.

**Status Aceita:** Prestador concordou em realizar o serviço.

**Status Rejeitada:** Prestador recusou a solicitação.

**Status Concluída:** Serviço foi realizado.

**Paginação:** Divisão de resultados em páginas para facilitar visualização.

**Tag:** Palavra-chave que ajuda na busca e categorização.

**Categoria:** Agrupamento principal de tipos de serviços.

---

## Observações Importantes

1. **Fotos:** Todos os links de fotos têm validade de 15 minutos por questões de segurança.

2. **Perfil Completo:** Usuário precisa informar tipo (dono/prestador) e endereço completo para ter perfil completo e acessar todas as funcionalidades.

3. **Localização:** Sistema usa coordenadas geográficas (latitude e longitude) para busca por proximidade.

4. **Segurança:** Senhas são criptografadas e nunca retornadas nas consultas.

5. **Email:** Verificação de email é obrigatória antes de fazer login.

6. **Exclusões:** Sistema usa "soft delete" (desativação) ao invés de exclusão definitiva para manter histórico.

7. **Permissões:** Cada funcionalidade valida cuidadosamente quem pode executá-la.

8. **Limites:** Há limites de caracteres em textos e quantidade de fotos para manter qualidade e performance.

---

**Última atualização:** Fevereiro de 2026
