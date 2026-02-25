# 🔍 Detail View Feature - Visualização e Edição Detalhada de Células

## Visão Geral

A funcionalidade de **Detail View** permite visualizar e editar o conteúdo completo de qualquer célula da tabela em uma interface dedicada, com formatação automática de JSON e validação integrada.

---

## 📋 Índice

1. [Como Usar](#como-usar)
2. [Visualização Detalhada](#visualização-detalhada)
3. [Edição de Conteúdo](#edição-de-conteúdo)
4. [Formatação JSON](#formatação-json)
5. [Validação e Segurança](#validação-e-segurança)
6. [Casos de Uso](#casos-de-uso)
7. [Limitações](#limitações)
8. [Exemplos Práticos](#exemplos-práticos)

---

## Como Usar

### Abrir Detail View

1. Execute uma query ou consulte uma tabela:
   ```bash
   pam run minha_query
   # ou
   pam tables usuarios
   ```

2. Navegue até a célula desejada usando:
   - `h/j/k/l` (vim-style)
   - Setas do teclado (`↑/↓/←/→`)

3. Pressione `Enter` para abrir a visualização detalhada

### Navegação no Detail View

| Tecla | Ação |
|-------|------|
| `↑` ou `k` | Rolar para cima |
| `↓` ou `j` | Rolar para baixo |
| `e` | Editar conteúdo (se editável) |
| `q` | Fechar e voltar para tabela |
| `Esc` | Fechar e voltar para tabela |
| `Enter` | Fechar e voltar para tabela |
| `Ctrl+c` | Sair do aplicativo |

---

## Visualização Detalhada

### O que é mostrado

```
◆ Cell Value - preferences (jsonb)
Row 5, Column 3 • Press 'e' to edit

────────────────────────────────────────

{
  "theme": "dark",
  "notifications": {
    "email": true,
    "push": false
  },
  "language": "pt-BR"
}

────────────────────────────────────────
[1-8 of 8 lines] ↑↓ scroll  e edit  q/esc/enter close
```

### Informações no Cabeçalho

- **Nome da coluna**: Ex: `preferences`
- **Tipo de dados**: Ex: `jsonb`, `text`, `varchar`
- **Posição**: Ex: `Row 5, Column 3`
- **Indicador de edição**: Aparece se a célula pode ser editada

### Scroll Automático

- Conteúdo longo é automaticamente paginado
- Indicador `[1-8 of 15 lines]` mostra a posição atual
- Use `↑/↓` ou `j/k` para navegar

---

## Edição de Conteúdo

### Quando é possível editar

A edição está disponível quando:
- ✅ A tabela tem nome identificado (`tableName`)
- ✅ Existe uma chave primária (`primaryKey`)
- ✅ Você está consultando uma tabela real (não uma VIEW ou JOIN)

### Processo de Edição

1. **No Detail View, pressione `e`**
   - Seu editor (`$EDITOR`) abre com o conteúdo atual
   - JSON é automaticamente formatado para facilitar edição

2. **Edite o conteúdo**
   - Modifique o texto livremente
   - Para JSON, mantenha a sintaxe válida
   - Pode adicionar/remover campos

3. **Salve e feche o editor**
   - Validação automática é executada
   - Se JSON, verifica se está válido
   - UPDATE é executado no banco de dados

4. **Confirmação visual**
   - Tabela atualiza imediatamente
   - Célula editada pisca em verde
   - Volta automaticamente para a visualização da tabela

### Configurar Editor

```bash
# Bash/Zsh
export EDITOR=vim        # ou nano, code, emacs, etc.

# Fish
set -gx EDITOR vim

# Adicionar ao ~/.bashrc ou ~/.zshrc para tornar permanente
echo 'export EDITOR=vim' >> ~/.bashrc
```

---

## Formatação JSON

### Detecção Automática

O sistema detecta JSON quando o conteúdo:
- Começa com `{` (objeto)
- Começa com `[` (array)
- É válido segundo `json.Unmarshal()`

### Formatação para Visualização

**Entrada no banco:**
```json
{"user":{"name":"Alice","roles":["admin","user"],"settings":{"theme":"dark"}}}
```

**Visualização formatada:**
```json
{
  "user": {
    "name": "Alice",
    "roles": [
      "admin",
      "user"
    ],
    "settings": {
      "theme": "dark"
    }
  }
}
```

### Formatação para Edição

Quando você pressiona `e`:
1. JSON é formatado com 2 espaços de indentação
2. Editor abre com conteúdo legível
3. Você edita o JSON formatado
4. Ao salvar, JSON é compactado antes de ir ao banco
5. Banco recebe JSON compacto (sem espaços extras)

### Tipos JSON Suportados

- ✅ Objetos: `{"key": "value"}`
- ✅ Arrays: `[1, 2, 3]`
- ✅ Nested: `{"a": {"b": {"c": "d"}}}`
- ✅ Mixed: `{"arr": [{"obj": "value"}]}`
- ✅ Tipos primitivos: strings, numbers, booleans, null

---

## Validação e Segurança

### Validação de JSON

Antes de salvar no banco:
```
1. Parser JSON verifica sintaxe
2. Se inválido → rejeita e mostra erro
3. Se válido → compacta e salva
```

**Exemplo de JSON inválido:**
```json
{name: "invalid"}  ❌ Faltam aspas nas chaves
```

**Mensagem de erro:**
```
✗ Error: Invalid JSON: invalid character 'n' looking for beginning of object key string
```

### Segurança na Atualização

1. **WHERE com Primary Key**
   ```sql
   UPDATE users 
   SET preferences = '{"theme":"light"}' 
   WHERE id = '123';  -- Sempre usa PK
   ```

2. **Validação de UPDATE Statement**
   - Verifica presença de WHERE clause
   - Requer primary key na condição
   - Previne updates acidentais em massa

3. **Atomic Operations**
   - Cada update é uma transação única
   - Falha não afeta outros dados
   - Rollback automático em caso de erro

### Proteções Implementadas

- ✅ Não permite UPDATE sem WHERE
- ✅ Valida JSON antes de salvar
- ✅ Usa prepared statements quando possível
- ✅ Escapa valores especiais
- ✅ Requer confirmação implícita (salvar editor)

---

## Casos de Uso

### 1. 📊 Inspecionar Dados JSON

**Problema:** JSON no banco está compacto e difícil de ler

**Solução:**
```bash
pam tables api_logs
# Navegar até coluna 'request_body'
# Pressionar Enter
# Ver JSON formatado claramente
```

### 2. ✏️ Editar Configurações JSON

**Problema:** Precisa atualizar preferências de usuário

**Solução:**
```bash
pam tables users
# Navegar até coluna 'preferences'
# Enter → 'e' → editar JSON
# Salvar → atualização automática
```

### 3. 🔍 Ler Textos Longos

**Problema:** Descrições longas truncadas na tabela

**Solução:**
```bash
pam tables products
# Navegar até 'description'
# Enter → ler texto completo
# Scroll com j/k se necessário
```

### 4. 🐛 Debug de API Responses

**Problema:** Verificar payload de API armazenado

**Solução:**
```bash
pam tables webhooks
# Navegar até 'response_payload'
# Enter → ver JSON formatado
# Identificar problema rapidamente
```

### 5. 🔧 Corrigir JSON Inválido

**Problema:** JSON quebrado no banco

**Solução:**
```bash
pam tables configurations
# Enter na célula com JSON inválido
# 'e' → corrigir no editor
# Validação garante que agora está correto
```

---

## Limitações

### Quando NÃO pode editar

❌ **Sem nome de tabela**
```bash
pam run "SELECT col1, col2 FROM (SELECT * FROM users) t"
# JOINs e subqueries complexas
```

❌ **Sem chave primária**
```sql
CREATE TABLE logs (
    message TEXT,
    timestamp TIMESTAMP
);
-- Tabela sem PK não permite edição
```

❌ **VIEWs**
```sql
CREATE VIEW active_users AS SELECT * FROM users WHERE active = true;
-- VIEWs geralmente não são editáveis
```

### Restrições Técnicas

- Tamanho máximo de conteúdo: limitado pela memória
- JSON muito grande pode demorar para formatar
- Editor depende de `$EDITOR` configurado
- Suporte a JSON apenas (XML não formatado ainda)

### Workarounds

**Para tabelas sem PK:**
```bash
# Use a tecla 'u' na tabela principal
# Edita mesmo sem PK (menos seguro)
```

**Para VIEWs:**
```bash
# Consulte a tabela base diretamente
pam tables base_table
```

---

## Exemplos Práticos

### Exemplo 1: Atualizar Preferências de Usuário

```bash
# 1. Conectar ao banco
pam switch production

# 2. Abrir tabela de usuários
pam tables users

# 3. Navegar até linha do usuário específico
# Use j/k para mover entre linhas

# 4. Navegar até coluna 'preferences' (h/l)

# 5. Enter para abrir detail view
# Conteúdo atual:
{
  "theme": "dark",
  "language": "en"
}

# 6. Pressionar 'e' para editar

# 7. No editor, modificar:
{
  "theme": "light",
  "language": "pt-BR",
  "notifications": {
    "email": true,
    "push": false
  }
}

# 8. Salvar e fechar (no vim: :wq)

# 9. ✅ Banco atualizado, célula pisca em verde!
```

### Exemplo 2: Corrigir JSON Quebrado

```bash
# JSON inválido no banco:
{"name":"John,"age":30}
#            ^ faltando aspas

# 1. pam tables users
# 2. Navegar até célula com JSON quebrado
# 3. Enter → ver erro de parsing
# 4. Pressionar 'e'
# 5. Corrigir no editor:
{"name":"John","age":30}
#            ^ adicionado aspas

# 6. Salvar → validação passa → banco atualizado!
```

### Exemplo 3: Inspecionar Logs de API

```bash
# Tabela com logs de requests
pam tables api_requests

# Navegar até coluna 'request_headers'
# Enter para ver JSON formatado:

{
  "Authorization": "Bearer eyJ...",
  "Content-Type": "application/json",
  "User-Agent": "MyApp/1.0",
  "X-Request-ID": "abc-123-def-456"
}

# Rolar com j/k para ver mais
# 'q' para fechar
```

### Exemplo 4: Bulk Edit com Script

```bash
#!/bin/bash
# Script para verificar configs antes de update em massa

for id in $(cat user_ids.txt); do
    echo "Checking user $id..."
    pam run "SELECT preferences FROM users WHERE id = $id"
    # Verificar manualmente com Enter
    # Editar se necessário com 'e'
    read -p "Continue? (y/n) " -n 1 -r
    echo
done
```

---

## Recursos Futuros

### Planejado

- [ ] Syntax highlighting para JSON (cores)
- [ ] Diff view antes de salvar
- [ ] Undo/Redo de edições
- [ ] Suporte para XML formatting
- [ ] Suporte para YAML formatting
- [ ] Copy formatted JSON to clipboard (Ctrl+C no detail view)
- [ ] Search/Find dentro do detail view (/)
- [ ] Line numbers no editor
- [ ] Multiple cell edit (batch update)

### Em Consideração

- [ ] Preview de mudanças antes de UPDATE
- [ ] History de edições por célula
- [ ] JSON schema validation
- [ ] Auto-complete para JSON conhecido
- [ ] Compare com versão anterior
- [ ] Export formatted JSON to file

---

## Troubleshooting

### "Press 'e' to edit" não aparece

**Causa:** Tabela não tem primary key ou não foi detectado

**Solução:**
```bash
# Verificar metadata da tabela
pam info tables <table_name>

# Usar pam tables em vez de query manual
pam tables <table_name>  # Detecta PK automaticamente
```

### Editor não abre ao pressionar 'e'

**Causa:** `$EDITOR` não configurado

**Solução:**
```bash
export EDITOR=vim  # ou nano, code, etc.
pam tables users   # tentar novamente
```

### "Invalid JSON" ao salvar

**Causa:** Sintaxe JSON incorreta

**Solução:**
1. Verificar aspas: `"` não `'`
2. Verificar vírgulas: última propriedade não tem vírgula
3. Usar validador online: https://jsonlint.com
4. Exemplo correto:
   ```json
   {
     "key": "value",
     "array": [1, 2, 3]
   }
   ```

### Conteúdo truncado mesmo no detail view

**Causa:** Limite de caracteres do banco

**Solução:**
1. Aumentar tamanho da coluna no banco
2. Usar tipo TEXT/LONGTEXT em vez de VARCHAR

### Update não funciona

**Possíveis causas:**
- ❌ Sem permissão UPDATE no banco
- ❌ Tabela é uma VIEW
- ❌ Trigger bloqueando UPDATE
- ❌ Constraint violada

**Debug:**
```bash
# Ver query exata executada
# Aparece na tela após update

# Testar query manualmente
pam run "UPDATE ... WHERE ..."
```

---

## Comparação com Outras Features

### vs. Tecla 'u' (Update Cell)

| Feature | Detail View + 'e' | Tecla 'u' Direta |
|---------|-------------------|------------------|
| Visualização | ✅ Mostra conteúdo completo | ❌ Não mostra antes |
| JSON Formatado | ✅ Sim | ❌ Não |
| Preview | ✅ Vê antes de editar | ❌ Edita direto |
| Validação JSON | ✅ Automática | ⚠️ Manual |
| Melhor para | JSON, textos longos | Valores simples |

**Recomendação:** Use Detail View para JSON e textos longos. Use 'u' para valores simples (números, strings curtas).

### vs. Copiar (Tecla 'y')

| Feature | Detail View | Copy 'y' |
|---------|-------------|----------|
| Propósito | Ver/Editar | Copiar |
| Formatação | ✅ JSON formatado | ❌ Valor bruto |
| Edição | ✅ Sim | ❌ Não |
| Clipboard | ❌ Não copia | ✅ Copia |

---

## Conclusão

A funcionalidade de **Detail View** transforma o PAM em uma ferramenta poderosa para:
- 📊 Inspecionar dados complexos (JSON)
- ✏️ Editar conteúdo de forma visual e segura
- 🔍 Ler textos longos sem truncamento
- 🐛 Debug de payloads e configurações

**Lembre-se:**
- `Enter` para abrir
- `e` para editar
- `q` para fechar

🎉 **Aproveite a nova funcionalidade e edite dados com confiança!**

---

**Versão:** 1.0
**Última Atualização:** 2024
**Autor:** PAM Development Team