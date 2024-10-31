# Multi-financial data websocket client on Golang

![IMG](docs/hero.gif)

## Task

1. Connect to the Coinbase cryptocurrency exchange via WebSocket
2. Subscribe to prices for three instruments: ETH-BTC, BTC-USD, BTC-EUR (data will come in through the WebSocket as they appear on the exchange)
3. Create an order table in MySql
4. Connect to MySql and write data received through WebSocket to the `order` table
5. Write data from the three instruments (ETH-BTC, BTC-USD, BTC-EUR) to the database
   in three threads (each instrument has its own write thread)
6. Publish the project repository on Github
7. Use good development practices to enable further application functionality expansion


## Setup for AWS
1. Create a new IAM user with programmatic access
2. Create a new secret in AWS Secrets Manager
3. Add the following keys to the secret:
   - COINBASE_WS_API_KEY
   - COINBASE_WS_API_SECRET
   - COINBASE_WS_API_PASSPHRASE
   - SUPABASE_URL
   - SERVICE_ROLE_KEY
4. Create image to push to ECR
   - `make build-app` creates the coinbase-app image 
   - `make build-analysis` creates the analysis-app image
5. Push images to ECR
   - `aws ecs register-task-definition --cli-input-json file://task-definition.json` creates a task definition (task-definition.json is in the root of the project)
   - `aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 886436943186.dkr.ecr.us-east-2.amazonaws.com` logs in to AWS ECR to push images
   - `docker push 886436943186.dkr.ecr.us-east-2.amazonaws.com/bz-findata-analysis` pushes the analysis-app image to ECR
   - `docker push 886436943186.dkr.ecr.us-east-2.amazonaws.com/bz-findata-coinbase` pushes the coinbase-app image to ECR


## Security Roles and Policies
1. Add the the following roles to the IAM user:
   - ecsTaskExecutionRole
   - ecsTaskRole

2. Add the following policies to the ecsTaskExecutionRole:
   - AmazonECSTaskExecutionRolePolicy
   - AmazonElasticFileSystemClientFullAccess
   - CloudWatchLogsFullAccess
   - ecsTaskExecutionPolicy
   - SecretsManagerReadWrite

3. Add the following policies to the ecsTaskRole:
   - ecsTaskPolicy
   - SecretsManagerReadWrite

## How to use

Project based on Clean architecture principles.

Requirements:

- Go 1.23 installed (https://go.dev/dl/)
- Docker installed (to run docker-compose)
- AWS CLI installed (https://awscli.amazonaws.com/v2/download/awscli2.exe)


## TODO

- [x] Logger points
- [ ] Rate limiter
- [ ] Prometheus metrics
- [ ] Testing


## Resources

- [Coinbase websocket overview](https://docs.cloud.coinbase.com/exchange/docs/websocket-overview)
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

