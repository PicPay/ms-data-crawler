# Data Crawler

## O que é?

E um micro serviço que foi criado para buscar dados de um micro serviço e apartir desses dados
bater em outro micro serviço para coletar os dados necessarios para popular um componente.
Essa arquitetura veio para resolver uma necessidade de dados, para buscar dados
no digital-goods, a partir dos dados do Midgard.
## Features

- Servidor HTTP
- logs no stdout
- banco de dados NoSQL MongoDB
- integraçao com microserviços externos para recuperaçao de informações
- Concorrência para buscar dados

## Configuração

Configuração é feita via variaveis de ambiente ou pelo arquivo `.env` (se ele existir)

```
DC_HTTP_ADDRESS=0.0.0.0:9003
DC_MONGO_ADDRESS=mongodb://localhost:27019
DC_MONGO_TIMEOUT=10
DC_MONGO_DATABASE=data-crawler

```

## Dependencias

Instala dependencias do projeto (com `go mod download`) e instala no seu `$GOPATH/bin`: [swag](https://github.com/swaggo/swag) e [reflex](https://github.com/cespare/reflex)

`make deps`

## Build

`make`

## Subir o projeto com docker

Caso você não tenha golang instalado, você poderá subir o projeto utilizando docker.
Basta executar o comando `make up` que ele irá subir os três serviços prioritários para você.
Após subir, você pode popular o banco com dados iniciais usando o seed.


## Próximos passos

- Ter a versatilidade para escolher qual campo, que vai ser a chave para iniciar a concorrencia
- Ajustar o mapping para utilizar header no serviço, assim deixando liberado para qualquer outra necessidade.