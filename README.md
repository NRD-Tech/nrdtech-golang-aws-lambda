# Golang AWS Lambda App

Template for a Go Lambda triggered by **EventBridge schedule**, **SQS**, or **API Gateway**. One trigger is active at a time via `trigger_type` in config; no Terraform files are commented out.

## Quick Start

- **Run tests:** `go test -v ./...`
- **Build:** `go build -v ./cmd/lambda`
- **Deploy staging:** push to `main` (GitHub Actions runs test then deploy).
- **Deploy production:** tag with `v*` (e.g. `v1.0.0`) and push the tag.
- **Destroy:** push tag `destroy-staging-*` or `destroy-prod-*` to tear down that environment.

## Technology Stack

- Go 1.24.3
- Docker
- Terraform

---

# Configuring the App for AWS Deployment

## OIDC Pre-Requisite

- An AWS role for OIDC and an S3 bucket for Terraform state must exist (e.g. via [nrdtech-terraform-aws-account-bootstrap](https://github.com/NRD-Tech/nrdtech-terraform-aws-account-bootstrap)).
- You need: **AWS Role ARN** and **Terraform state bucket name**.

## Configure Settings

1. **Edit `config.global`** — set at least:
   - `APP_IDENT_WITHOUT_ENV`
   - `TERRAFORM_STATE_BUCKET`
   - `AWS_DEFAULT_REGION`
   - `AWS_ROLE_ARN`

2. **Set `trigger_type`** in `config.global` (or override in `config.staging` / `config.prod`):
   - `trigger_type=sqs` — SQS queue (default)
   - `trigger_type=api_gateway` — API Gateway HTTP API
   - `trigger_type=eventbridge` — EventBridge Scheduler (cron/rate)

   All three Terraform trigger files stay in the repo; only the one matching `trigger_type` is applied. To switch triggers, change this value and re-deploy (no uncommenting or manual state fixes).

3. **Ensure `cmd/lambda/main.go` matches your trigger** — use the corresponding template as your entrypoint:
   - **SQS:** `main_sqs_trigger.go.tmpl` → build as `main.go`
   - **API Gateway:** `main_api_gateway.go.tmpl` → build as `main.go`
   - **EventBridge:** `main_event_bridge.go.tmpl` → build as `main.go`

4. **If using API Gateway:** set `API_DOMAIN` and `API_ROOT_DOMAIN` in `config.staging` and `config.prod` (optional; only needed for custom domain). `API_ROOT_DOMAIN` must exist in Route53.

5. **If using EventBridge:** adjust `schedule_expression` in `terraform/main/lambda_eventbridge_schedule.tf` if needed (cron or rate).

6. **Edit `go.mod`** — set the module name (e.g. same as `APP_IDENT_WITHOUT_ENV`).

## GitHub Actions (GitHub Flow)

- Workflow: `.github/workflows/github_flow.yml`
- Ensure `role-to-assume` and `aws-region` in the workflow match `AWS_ROLE_ARN` and `AWS_DEFAULT_REGION` in `config.global` (or use the same config source).
- **Staging:** push to `main` → test job runs, then deploy to staging.
- **Production:** create a version tag (e.g. `git tag v1.0.0`) and push it (`git push origin v1.0.0`) → test then deploy to prod.
- **Destroy staging:** push a tag matching `destroy-staging-*` (e.g. `destroy-staging-20250227`).
- **Destroy production:** push a tag matching `destroy-prod-*` (e.g. `destroy-prod-20250227`).

## Bitbucket

- To use Bitbucket Pipelines instead, enable Pipelines in repository settings. The project includes `bitbucket-pipelines.yml` for that flow.

---

# Setting Up Your Development Environment

## Clone and clean the template (GitHub)

- Use “Use this template” on [NRD-Tech/nrdtech-golang-aws-lambda](https://github.com/NRD-Tech/nrdtech-golang-aws-lambda) to create a new repository, then clone it.
- For proprietary use: [Proprietary licensing](#how-to-use-this-template-for-a-proprietary-project).

## Clone and clean the template (without GitHub)

```bash
git clone https://github.com/NRD-Tech/nrdtech-golang-aws-lambda.git my-project
cd my-project
rm -fR .git .idea
git init
git add .
git commit -m 'init'
```

---

# Deploy and destroy

## Deploy to staging

```bash
git checkout main
git add .
git commit -m 'your message'
git push origin main
```

## Deploy to production

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Destroy (GitHub Actions)

- **Staging:** `git tag destroy-staging-$(date +%Y%m%d)` then `git push origin destroy-staging-$(date +%Y%m%d)`
- **Production:** `git tag destroy-prod-$(date +%Y%m%d)` then `git push origin destroy-prod-$(date +%Y%m%d)`

---

# Misc How-To's

## How to use this template for a proprietary project

This project's license (MIT License) allows you to create proprietary code based on this template. Replace the LICENSE file with your terms and optionally add a NOTICE file stating that the original work is licensed under the MIT License.

## Run the Lambda Docker image locally

```bash
aws ecr get-login-password --region us-west-2 | \
  docker login --username AWS --password-stdin <account>.dkr.ecr.<region>.amazonaws.com

docker run --rm -p 9000:8080 -it <account>.dkr.ecr.<region>.amazonaws.com/<repository>:latest

curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
```

## Inspect the Docker image (dive)

```bash
alias dive="docker run -ti --rm -v /var/run/docker.sock:/var/run/docker.sock wagoodman/dive"
dive <account>.dkr.ecr.<region>.amazonaws.com/<repository>:latest
```

## View image architecture

```bash
export AWS_PROFILE=yourprofile
aws ecr get-login-password --region us-west-2 | \
  docker login --username AWS --password-stdin <account>.dkr.ecr.<region>.amazonaws.com
docker pull <account>.dkr.ecr.<region>.amazonaws.com/<repository>:latest
docker inspect <account>.dkr.ecr.<region>.amazonaws.com/<repository>:latest
```
