# bp-transaction-api

## Transaction Broadcasting and Monitoring Client Module

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Integration](#integration)

## Getting Started

This client module will enable the broadcasting of a transaction and subsequently monitor its status until finalization

### Prerequisites

List any software, tools, or dependencies that need to be installed before using your project.

- Go (version 1.21.1)
- [Gin (HTTP Web framework)](github.com/gin-gonic/gin)
- [Validator (v10)](github.com/go-playground/validator/v10)
- [Viper - Configuration handling](github.com/spf13/viper)
- [Zap Logger](go.uber.org/zap)

### Installation

Step-by-step guide on how to install your project:

1. Clone this repository: `git clone https://github.com/newsrw/bp-transaction-api.git`
2. Navigate to the project directory: `cd bp-transaction-api`
3. Build the project: `make build`
4. Run the project: `./bin/app`

## Configuration

Configuration file is store in /configs/config.yml file
This project use Viper to convert .yml file to go config structure

Config structure:

```
CONFIG:
  SERVER:
    ADDRESS: // Server address
  CLIENT:
    PORT: // client port
    MONITORDELAYDURATION: // Transaction monitoring status delay duration
```

## Usage

Here are use cases:

1. Healthcheck - to check whether our server is ready to use or not

```bash
curl --location 'http://localhost:8080/healthcheck'
```

2. Broadcast Transaction - to broadcast transaction and monitoring

```bash
curl --location 'http://localhost:8080/transactions/broadcast' \
--header 'Content-Type: application/json' \
--data '{
    "symbol": "ETH",
    "price": 4500,
    "timestamp": 1678912345
}'
```

All transaction status:

- `CONFIRMED`: Transaction has been processed and confirmed
- `FAILED`: Transaction failed to process
- `PENDING`: Transaction is awaiting processing
- `DNE`: Transaction does not exist

For monitoring transaction status strategy, the service use tx_hash that got from the broadcast response(POST /broadcast) and store it into a variable
Next, the service'll start calling GET /check/{tx_hash} for getting the status and store it in another variable and logging the status
If the status is CONFIRMED, FAILED, or DNE, the service'll return that response back to let user know that tx_hash is already processed and confirmed, failed, or doesn't exist.

1. To make periodically issue

## Integration

To integration this service, you can build and push docker image onto the server
and pull the image to

Here is the script to build

```bash
docker build -t bp-transaction-api .
```
