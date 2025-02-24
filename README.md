# Golang AWS Lambda App
This is a project template for a golang application that will be triggered either by an Event Bridge schedule, an SQS queue, or an API Gateway endpoint

# Technology Stack
* Go 1.23.3
* Docker
* Terraform

# Setting Up Your Development Environment

## Clone and Clean the template (if using GitHub)
* Navigate to: https://github.com/NRD-Tech/nrdtech-golang-aws-lambda.git
* Log into your GitHub account (otherwise the "Use this template" option will not show up)
* Click "Use this template" in the top right corner
  * Create a new repository
* Fill in your repository name, description, and public/private setting
* Clone your newly created repository
* If you want to change the license to be proprietary follow these instructions: [Go to Proprietary Licensing Section](#how-to-use-this-template-for-a-proprietary-project)

## Clone and Clean the template (if NOT using GitHub)
```
git clone https://github.com/NRD-Tech/nrdtech-golang-aws-lambda.git my-project
cd my-project
rm -fR .git venv .idea
git init
git add .
git commit -m 'init'
```
* If you want to change the license to be proprietary follow these instructions: [Go to Proprietary Licensing Section](#how-to-use-this-template-for-a-proprietary-project)

# Configuring the App for AWS Deployment

## OIDC Pre-Requisite
* You must have previously set up the AWS Role for OIDC and S3 bucket for the Terraform state files
* The easiest way to do this is to use the NRD-Tech Terraform Bootstrap template
  * https://github.com/NRD-Tech/nrdtech-terraform-aws-account-bootstrap
  * After following the README.md instructions in the bootstrap template project you should have:
    * An AWS Role ARN
    * An AWS S3 bucket for the Terraform state files

## Configure Settings
* Edit .env.global
  * Each config is a little different per application but at a minimum you will need to change:
    * APP_IDENT_WITHOUT_ENV
    * TERRAFORM_STATE_BUCKET
    * AWS_DEFAULT_REGION
    * AWS_ROLE_ARN
* Edit go.mod
  * Set an appropriate module name (likely the same as APP_IDENT_WITHOUT_ENV)
* Choose how your lambda function will be triggered and un-comment the appropriate terraform:
  * Event Bridge Scheduling:
    * Un-comment terraform/main/lambda_eventbridge_schedule.tf
    * Set the schedule that you want as a cron or rate in terraform/main/lambda_eventbridge_schedule.tf
    * Rename cmd/lambda/main_event_bridge.go.tmp to cmd/lambda/main.go
  * SQS Triggered:
    * Un-comment terraform/main/lambda_sqs_trigger.tf
    * Rename cmd/lambda/main_sqs_trigger.go.tmp to cmd/lambda/main.go
  * API Gateway:
    * Un-comment terraform/main/lambda_api_gateway.tf
    * Rename cmd/lambda/main_api_gateway.go.tmp to cmd/lambda/main.go
    * Configure the domain's in .env.prod and .env.staging
* Commit your changes to git
```
git add .
git commit -a -m 'updated config'
```

## (If using Bitbucket) Enable Bitbucket Pipeline (NOTE: GitHub does not require any setup like this for the Actions to work)
* Push your git project up into a new Bitbucket project
* Navigate to your project on Bitbucket
  * Click Repository Settings
  * Click Pipelines->Settings
    * Click Enable Pipelines

## (If using GitHub) Configure the AWS Role
* Edit .github/workflows/main.yml
  * Set the pipeline role for role-to-assume
    * This should be the same as the AWS_ROLE_ARN in your .env.global
  * Set the correct aws-region

## Deploy to Staging
```
git checkout -b staging
git push --set-upstream origin staging
```

## Deploy to Production
```
git checkout -b production
git push --set-upstream origin production
```

## Un-Deploying in Bitbucket
1. Navigate to the Bitbucket project website
2. Click Pipelines in the left nav menu
3. Click Run pipeline button
4. Choose the branch you want to un-deploy
5. Choose the appropriate un-deploy Pipeline
   * un-deploy-staging
   * un-deploy-production
6. Click Run

# Misc How-To's

## How to use this template for a proprietary project
This project's license (MIT License) allows for you to create proprietary code based on this template.

Here are the steps to correctly do this:
1. Replace the LICENSE file with your proprietary license terms if you wish to use your own license.
2. Optionally, include a NOTICE file stating that the original work is licensed under the MIT License and specify the parts of the project that are governed by your proprietary license.

## How To run docker image locally
```
aws ecr get-login-password \
      --region us-west-2 | \
      docker login \
        --username AWS \
        --password-stdin 1234567890.dkr.ecr.us-west-2.amazonaws.com/myapp_lambda_repository

docker run --rm -p 9000:8080 -it 482370276428.dkr.ecr.us-west-2.amazonaws.com/myapp_lambda_repository:latest 

curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'

```

# How To Inspect docker image
```
alias dive="docker run -ti --rm  -v /var/run/docker.sock:/var/run/docker.sock wagoodman/dive"
dive 1234567890.dkr.ecr.us-west-2.amazonaws.com/myapp_lambda_repository:latest
```

# How to view the architecture (and other info) of a docker image
```
export AWS_PROFILE=mycompanyprofile
docker logout 1234567890.dkr.ecr.us-west-2.amazonaws.com
aws ecr get-login-password \
      --region us-west-2 | \
      docker login \
        --username AWS \
        --password-stdin 1234567890.dkr.ecr.us-west-2.amazonaws.com/myapp_lambda_repository
docker pull 1234567890.dkr.ecr.us-west-2.amazonaws.com/myapp_lambda_repository
docker inspect 1234567890.dkr.ecr.us-west-2.amazonaws.com/myapp_lambda_repository
```
