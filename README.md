# Template Manager

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=alert_status&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=coverage&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=sqale_rating&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=duplicated_lines_density&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=ncloc&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=sqale_index&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=code_smells&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=PicPay_picpay-dev-ms-template-manager&metric=bugs&token=78adc83a0fee4bef3522943aeecb260da841c45b)](https://sonarcloud.io/dashboard?id=PicPay_picpay-dev-ms-template-manager)

## O que é?

O ms template-manager tem como objetivo e motivação principal promover a renderização dinamica de screens por meio de um motor da geração de um objeto JSON.

Cada "screen" possuirá um template JSON que será o "esqueleto" que por sua vez será preenchido por valores que serão recuperados a partir de outros micro-serviços, por exemplo ms-store que é o micro-serviço responsavel por armazenar as lojas da Store, a partir do template e das informações recuperadas de outros ms será gerado um JSON final, que é a screen renderizada pelo front-end.

## Features

- Servidor HTTP
- logs no stdout
- banco de dados NoSQL MongoDB
- integraçao com microserviços externos para recuperaçao de informações
- motor de renderização de JSON a partir de um template dinamico

## Configuração

Configuração é feita via variaveis de ambiente ou pelo arquivo `.env` (se ele existir)

```
TEMPLATE_MANAGER_HTTP_ADDRESS=0.0.0.0:8000
TEMPLATE_MANAGER_MONGO_ADDRESS=mongodb://localhost:27017
TEMPLATE_MANAGER_MONGO_TIMEOUT=10
TEMPLATE_MANAGER_MONGO_DATABASE=template-manager

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
