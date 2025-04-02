# Aplicação de Cotação de Dólar em Go

Este projeto consiste em dois sistemas escritos em Go:

- `server.go`: Um servidor HTTP que busca a cotação do dólar e a armazena em um banco SQLite.
- `client.go`: Um cliente que solicita a cotação ao servidor e salva o valor em um arquivo de texto.

## Requisitos

- Go instalado ([https://go.dev/dl/](https://go.dev/dl/))
- SQLite3 instalado
- SQLite Viewer extensão instalada no VSCode, para ser possível a visualização dos dados da cotação no banco de dados.

## Instalação

Clone este repositório e navegue até o diretório do projeto:

```sh
git clone https://github.com/AvelinoRSN/desafio-client-server-api.git
cd seu-repositorio
```

## Executando o Servidor

O servidor busca a cotação do dólar na API `https://economia.awesomeapi.com.br/json/last/USD-BRL`, salva no SQLite e retorna a cotação via HTTP.

1. Execute o seguinte comando para iniciar o servidor:

```sh
go run server.go
```

O servidor rodará na porta `8080`.

## Executando o Cliente

O cliente solicita a cotação ao servidor e salva o valor em um arquivo `cotacao.txt`.

1. Em um novo terminal, execute:

```sh
go run client.go
```

Se tudo correr bem, um arquivo `cotacao.txt` será criado com o seguinte formato:

```
Dólar: {valor}
```

## Considerações

- O servidor tem um timeout de `200ms` para chamar a API externa e `10ms` para salvar no banco.
- O cliente tem um timeout de `300ms` para receber a resposta do servidor.
- Caso o tempo estoure, erros serão registrados nos logs.

## Autor

Avelino Ramos

