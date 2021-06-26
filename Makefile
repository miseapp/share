include .env
include ./Makefile.base.mk

# -- cosmetics --
help-colw = 7

# -- constants --
df-infra = infra
df-tf = $(df-infra)/.terraform
df-table = $(SHARE_COUNT_NAME)
df-item-id = $(SHARE_FILES_NAME)
db-entry = cmd/main.go
db-build = build
db-binary = $(SHARE_ADD_BINARY)
db-archive = $(SHARE_ADD_ARCHIVE)
dr-fn = $(SHARE_ADD_NAME)
dr-payload = payload.json
dr-endpoint = $(AWS_ENDPOINT)

# -- tools --
ti-brew = brew
tf-dc = docker-compose
tf-terraform = . .env && terraform -chdir="$(df-infra)"
tb-go = go
tt-go = go
tr-aws = aws
tr-env = AWS_CONFIG_FILE=.aws/config AWS_SHARED_CREDENTIALS_FILE=.aws/creds

# -- init --
## [i]init dev env
init: i
.PHONY: init

i: i/pre
	$(ti-brew) bundle -v --no-upgrade
.PHONY: i

## updates deps
i/upgr:
	$(ti-brew) bundle -v
.PHONY: i/upadte

# -- i/helpers
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
	GOOS=linux GOARCH=amd64 $(tb-go) build -o $(db-binary) $(db-entry)
.PHONY: b

## clean the build
b/clean:
	rm -rf $(db-build)
.PHONY: b/clean

## build and archive
b/arch: b
	zip $(db-archive) $(db-binary)
.PHONY: b/arch

# -- run --
## [r]un handler fn
run: r
.PHONY: run

r:
	$(tr-env) \
	$(tr-aws) lambda invoke \
	--function-name $(dr-fn) \
	--payload $$(base64 < $(dr-payload)) \
	--endpoint-url=$(dr-endpoint) \
	--debug \
	test.json
.PHONE: r

# -- test --
## run [t]ests
test: t
.PHONY: test

t:
	$(tt-go) test ./... -run "_U"
.PHONY: t

## runs unit & int tests
t/all:
	$(tt-go) test ./...
.PHONY: t/all

# -- infra --
## in[f]ra; aliases f/plan
infra: f
.PHONY: infra

f: f/plan
.PHONY: f

## start localstack
f/start:
	$(tf-dc) up
.PHONY: f/start

## stop localstack
f/stop:
	$(tf-dc) down
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
	$(tf-terraform) apply -auto-approve
.PHONY: f/apply

## init infra state
f/init:
	$(tr-env) \
	$(tr-aws) dynamodb put-item \
	--table-name $(df-table) \
	--item '{ "Id": {"S": "$(df-item-id)"}, "Count": {"N": "0"} }' \
	--endpoint-url=$(dr-endpoint) \
	--debug
.PHOYN: f/init

## destroy infra
f/clean:
	$(tf-terraform) destroy
.PHONY: f/reset

# -- i/helpers
$(df-tf):
	$(tf-terraform) init
