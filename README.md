# Auto Voter

Programa em Go para automatizar votos no formulário.

## Descrição

Este programa envia múltiplos votos automaticamente para a primeira opção de um formulário especifico.

## Requisitos

- Go 1.21 ou superior instalado
- Conexão com a internet

## Como Usar

### 1. Executar o programa

```powershell
go run auto_voter.go
```

### 2. Ou compilar e executar

```powershell
# Compilar
go build auto_voter.go

# Executar
.\auto_voter.exe
```

## Configurações

Você pode ajustar os seguintes parâmetros no código:

- **`totalVotes`** (linha 66): Número total de votos a enviar (padrão: 100)
- **`time.Sleep`** (linha 81): Intervalo entre votos em milissegundos (padrão: 500ms)

### Para votar em outra opção:

Altere a constante `voteValue` (linha 17), analise o HTML e informe a opção.


## Exemplo de Saída

```
==========================================
      ("  Auto Voter - Votação")
==========================================

Enviando 100 votos...

✓ Voto #1 enviado com sucesso! (Status: 200)
✓ Voto #2 enviado com sucesso! (Status: 200)
...
--- Progresso: 10/100 votos enviados ---
...
==========================================
  Resumo Final
==========================================
Total de votos tentados: 100
✓ Votos bem-sucedidos:   100
✗ Votos com erro:        0
Taxa de sucesso:         100.0%
==========================================
```

## Observações

⚠️ **Aviso**: Este código é apenas para fins educacionais. Usar scripts para manipular votações pode violar os termos de serviço do Google Forms.

## Funcionalidades

- ✅ Envio automático de votos
- ✅ Contador de progresso
- ✅ Tratamento de erros
- ✅ Intervalo configurável entre votos
- ✅ Relatório final de estatísticas
- ✅ User-Agent realista
- ✅ Timeout de requisições

## Estrutura do Código

- **`submitVote()`**: Função que envia um único voto
- **`main()`**: Função principal que coordena o envio em loop
- Headers HTTP apropriados para simular navegador
- Gerenciamento de redirects do Google Forms
