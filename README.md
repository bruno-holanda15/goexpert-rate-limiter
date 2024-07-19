# Rate Limiter

Projeto desenvolvido na linguagem Go proposto pelo desafio técnico da Pós em Go da FullCycle.

## Descrição

Inicia um servidor http permitindo request para o endpoint http://localhost:8080/ping.
Contudo o mesmo possui um limite de requests X por IP e/ou Token (enviado pelo Header API_KEY) e por Y tempo que deve ser configurado pelo usuário que for clonar o projeto.
O mesmo contém um arquivo .env preenchido, caso deseje alterar os valores e tipos de periodicidade para limite, deve alterar as seguintes variáveis.

- RATE_LIMIT_IP é a quantidade limite para request por um IP antes de bloqueá-lo.

- RATE_LIMIT_TOKEN é a quantidade limite para request por um Token antes de bloqueá-lo.

- TIME_LIMIT_TYPE é o tipo de tempo para limite de requests, entre ele second, minute e hour.

- TIME_BLOCK_TYPE é o tipo de tempo para bloqueio de requests, entre ele second, minute e hour.

- BLOCK_LIMIT_TIME_DURATION é a duração de tempo que será bloqueado para executar as requests.

## Pré-requisitos
```
Docker
Docker Compose
Go
```

## Como rodar?
O projeto possui um arquivo Makefile para facilitar a execução.

Inicia o projeto executando os containers.
```
make up
```

Exibe os logs do container rate-limiter.
```
make logs
```

Encerra a execução dos containers.
```
make down
```

Executa os testes do projeto Go e cria um arquivo html coverage.html na raíz do projeto.
```
make tests
```