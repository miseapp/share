include ./Makefile.base.mk

# -- cosmetics --
help-colw = 7

# -- constants --
df-infra = infra
df-tf = $(df-infra)/.terraform
db-build = build
db-binary = $(db-build)/share
db-archive = $(db-binary).zip
dr-fn = share.add
dr-endpoint = http://localhost:4566

# -- tools --
ti-brew = brew
tf-docker = docker
tf-localstack = localstack
tf-terraform = terraform -chdir="$(df-infra)"
tr-aws = AWS_CONFIG_FILE=.aws/config AWS_SHARED_CREDENTIALS_FILE=.aws/creds aws

# -- init --
## setup dev env
init: i
.PHONY: init

i: i/pre
	$(ti-brew) bundle -v
.PHONY: i

i/pre:
ifeq ("$(shell command -v $(ti-brew))", "")
	$(info ✘ brew is not installed, please see:)
	$(info - https://brew.sh)
	$(error 1)
endif
.PHONY: i/pre

# -- build --
## [b]uild handler fn
build: b
.PHONY: build

b:
	GOOS=linux go build -o $(db-binary) cmd/main.go
.PHONY: b

## build and archive
b/arch: b
	zip $(db-archive) $(db-binary)
.PHONY: b/arch

# -- run --
## [r]un handler fn
run: r
.PHONY: run

r:
	$(tr-aws) lambda get-function --function-name $(dr-fn) --endpoint-url=$(dr-endpoint)
.PHONE: r

# -- test --
## run tests
test: t
.PHONY: test

t:
	go test pkg/...
.PHONY: t

# -- infra --
## in[f]ra; aliases f/start
infra: f
.PHONY: infra

f: f/start
.PHONY: f

## start localstack
f/start:
	$(tf-localstack) start
.PHONY: f/start

## stop localstack
f/stop:
	$(tf-docker) stop localstack_main
.PHONY: f/stop

## plan dev infra
f/plan: $(df-tf)
	$(tf-terraform) plan
.PHONY: f

## validate infra
f/valid:
	$(tf-terraform) validate
.PHONY: f/validate

## apply planned infra
f/apply:
	$(tf-terraform) apply
.PHONY: f/apply

## destroy infra
f/reset:
	$(tf-terraform) destroy
.PHONY: f/reset

# -- i/helpers
$(df-tf):
	$(tf-terraform) init
