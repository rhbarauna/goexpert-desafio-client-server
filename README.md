## USD to BRL Exchange Rate System

## Overview

This system fetches the current exchange rate between the US dollar (USD) and the Brazilian real (BRL). It consists of two components:

- **Server component:**

  - Retrieves exchange rate data from: [https://economia.awesomeapi.com.br/json/last/USD-BRL](https://economia.awesomeapi.com.br/json/last/USD-BRL)
  - Persists the retrieved exchange rate into an SQLite database
  - Exposes an endpoint at: `{server_url}:8080/cotacao`

- **Client component:**
  - Interacts with the server through the `/cotacao` endpoint.
  - Retrieves and displays the received exchange rate.
  - Saves the retrieved rate to a file named `cotacao.txt`.

## Components

- **Server:**
  - Located in the `server/cmd` directory.
  - Main program file: `server.go`.
- **Client:**
  - Located in the `client/cmd` directory.
  - Main program file: `client.go`.

## Prerequisites

- Go language installed ([https://golang.org/dl/](https://golang.org/dl/))

## Running the System

1. **Start the Server:**

   - Navigate to the `server/cmd` directory.
   - Run the command: `go run server.go`

2. **Run the Client:**
   - Open a separate terminal window.
   - Navigate to the `client/cmd` directory.
   - Run the command: `go run client.go`
