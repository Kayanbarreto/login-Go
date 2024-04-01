# API de Autenticação de Usuário

Esta é uma api de Autenticação de Usuário, criado com o intuito de...

## Endpoints

### `POST /v1/login`

Este endpoint permite realizar login na API via dados JSON.

#### Body Envio

```json
{
	"email": "zeca.spr3@gmail.com",
	"password": "91951891ab"
}
```

#### Body Retorno - Código 200 -OK

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1hdGV1cy5zcHIzQGdtYWlsLmNvbSIsImV4cCI6MTcxMTkxNzg5MX0.YEYL1WaMsdSRWcn5Fwr-Xw2AW1G6UaYKANX8hXXc8vE"
}
```

### `POST /v1/user/register`

Este endpoint permite realizar cadastro de usuário na API via dados JSON.

#### Body Envio

```json
{
  "email": "mateus.spr3@gmail.com",
  "password": "91951891a",
  "birthDate": "1990-01-01T00:00:00Z",
  "name": "sas",
  "active": true
}
```

#### Body Retorno - Código 201 - CREATED

```json
{
  "id": "06326436-01d3-45d8-bd9b-3bfff952d1e6",
  "email": "mateus.spr3@gmail.com",
  "password": "$2a$10$76dAIOvRv2nF7MUv649WvOdnM8h9qkHE17AXrMj0JgC/rOJ7swroi",
  "birthDate": "1990-01-01T00:00:00Z",
  "name": "sas",
  "active": true
}
```

### `PUT /v1/user/updateUser`

Este endpoint permite realizar atualização de usuário na API via dados JSON.

#### Body Envio

```json
{
  "email": "mateus.spr3@gmail.com",
  "password": "91951891ab",
  "birthDate": "1990-01-01T00:00:00Z",
  "name": "sas",
  "active": true
}
```

#### Body Retorno - Código 200 - OK

```json
{
  "message": "Usuário atualizado com sucesso"
}
```

### `POST /v1/user/generate-token-recover-password`

Este endpoint permite gerar token de recuperação de senha.

#### Body Envio

```json
{
  "email": "mateus.spr10@gmail.com"
}
```

#### Body Retorno - Código 200 - OK

```json
{
  "message": "Token de recuperação enviado para o email",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1hdGV1cy5zcHIxMEBnbWFpbC5jb20iLCJleHAiOjE3MTE5MTc4NTR9.-eyQ_KAND8EmoHn8WUf41lhAkMC59CKGMfXHef-gcys"
}
```

### `POST /v1/user/recover-password`

Este endpoint permite a recuperação de senha.

#### Body Envio

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1hdGV1cy5zcHIxMEBnbWFpbC5jb20iLCJleHAiOjE3MTE4OTcyMTJ9.-o1p1oM5JSek1uv1VVcaV-hKM8wZIhWi9ymyWUvdPOk",
  "new_password": "456uia"
}

```

#### Body Retorno - Código 200 - OK

```json
{
  "message": "Senha recuperada com sucesso"
}
```

### Para rodar a aplicação
Rodar o seguinte comando na raiz do projeto.
```
docker-compose up
go run . 
```

### Pré-Requisitos
```
Golang instalado na máquina.
Docker e Docker compose instalado na máquina. 
```