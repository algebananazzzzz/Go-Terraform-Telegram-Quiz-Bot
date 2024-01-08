# Go Terraform Telegram Quiz Bot

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/yourusername/Go-Terraform-Telegram-Quiz-Bot/blob/main/LICENSE)

## Table of Contents

- [About](#about)
  - [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Configuration](#configuration)
  - [Quiz Data Configuration](#quiz-data-configuration)
  - [Environment Configuration](#environment-configuration)
  - [VPC Configuration](#vpc-configuration)
- [Local Development Locally](#local-development)
- [Deployment](#deployment)
  - [Using Terraform Locally](#using-terraform-locally)
  - [Setting Webhook](#setting-webhook)
- [License](#license)
- [Contact](#contact)


## About:
Go-Terraform-Telegram-Quiz-Bot is a project that leverages the speed and efficiency of Golang and Terraform Infrastructure as Code (IAC) to provision serverless resources on AWS. The project is delivered with the application code, and users only need to configure environment variables and quiz data to get started.

If you're here only with the interest for the Go code for the bot, feel free to explore the `src` folder!
<!-- This project simplifies the deployment and scaling process, allowing users to focus on crafting engaging quizzes rather than dealing with infrastructure and code complexities. -->

### Features

- **Serverless Architecture:** Leverages AWS Lambda for serverless execution, ensuring scalability (>1000 concurrent executions) and cost efficiency (free-tier compatible) execution.
  
- **Golang for Speed:** The project is built using the Go programming language, known for its speed and performance.

- **Terraform IAC:** Infrastructure is defined and managed as code using Terraform, enabling reproducibility and version control.

- **Easy Customization:** Users can get started quickly by configuring environment variables and updating quiz data to personalize their Telegram quiz bots.

## Getting Started
### Prerequisites

1. **Golang version >=1.21**: You can download Golang from [the official Golang website](https://go.dev/doc/install).

2. **Terraform version >=0.14.0**: You can download Terraform from [the official Terraform website](https://www.terraform.io/).

3. **AWS VPC and Redis Elasticache Cluster**: Create a suitable subnet with connectivity to a Redis (v7) endpoint, and internet connection through NAT (recommended) or IGW.

Follow these steps to set up your own scalable Telegram quiz bots:

### Installation
1. **Clone the Repository:**
```bash
git clone https://github.com/algebananazzzzz/Go-Terraform-Telegram-Quiz-Bot.git
```
   
## Configuration

### Quiz Data Configuration

Under `src/quizdata.yaml`, modify the questions, options, answers using the format: 

```yaml
- id: 0
  quiz_name: Section A
  quiz_length: 5
  questions:
  - id: 0
    title: "1. What is the capital of France?"
    options:
      - "Berlin"
      - "Madrid"
      - "Paris"
      - "Rome"
    answer: 2 # correct answer: Paris
  - id: 1
    title: "2. Which planet is known as the Red Planet?"
    options:
      - "Venus"
      - "Mars"
      - "Jupiter"
      - "Saturn"
    answer: 1 # correct answer: Mars
```

### Environment Variable Configuration

Under `main.tf`, modify the local variable `environment_variables`:

| Variable | Type | Description |
|---|---|---|
| TOKEN | string | API token obtained from botfather | 
| REDIS_ADDR | string | Redis cluster address ending with port number | 
| REDIS_KEY | string | Redis key to store bot data under | 
| START_MESSAGE | string | Message displayed on /start | 
| ERROR_MESSAGE | string | Message displayed on unexpected errors | 
| PASSPHRASE | string (optional) | Configures a passphrase users require to access quizzes | 

### VPC Configuration

Under `data.tf`, modify the following local variables:

| Variable | Type | Description |
|---|---|---|
| subnet_names | string | Names of subnets to launch lambda function in | 
| security_groups | string | Names of security groups to associate with function (e.g. allow nat) | 


## Local Development

1. Comment out `src/production.go` and uncomment `src/local.go`.
2. Set [Environment Variables](#environment-configuration) within the src directory.

```shell
cd src/
export TOKEN=
export REDIS_ADDR=localhost:6379
export REDIS_KEY=
export START_MESSAGE=
export ERROR_MESSAGE=
export PASSPHRASE=
```
3. Run bot with ```go run .```

## Deployment 

### Using Terraform Locally

1. Build Go binary

```shell
cd src/
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o ../build/bootstrap .
```

2. Provision AWS resources using Terraform Locally

```shell
cd infra
terraform init
terraform workspace select -or-create myworkspace # optional
terraform apply
```

### Setting Webhook

Substitute the provided url with the bot api token and output `api_gateway_endpoint` from `terraform apply`. Then, open the url in a new browser window.

```
https://api.telegram.org/bot<TOKEN>/setWebhook?url=<ENDPOINT>
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
