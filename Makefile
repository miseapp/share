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

## -- init (i) --
$(eval $(call alias, init, i/0))
$(eval $(call alias, i, i/0))

## init dev env
i/0: i/pre
	$(ti-brew) bundle -v --no-upgrade
.PHONY: i/0

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

## -- build (b) --
$(eval $(call alias, build, b/0))
$(eval $(call alias, b, b/0))

## build handler fn
b/0:
	GOOS=linux GOARCH=amd64 $(tb-go) build -o $(db-binary) $(db-entry)
.PHONY: b/0

## build & archive
b/arch: b
	zip $(db-archive) $(db-binary)
.PHONY: b/arch

## clean build dir
b/clean:
	rm -rf $(db-build)
.PHONY: b/clean

## -- run (r) --
$(eval $(call alias, run, r/0))
$(eval $(call alias, r, r/0))

## call local handler fn
r/0: r
	$(tr-env) \
	$(tr-aws) lambda invoke \
	--function-name $(dr-fn) \
	--payload $$(base64 < $(dr-payload)) \
	--endpoint-url=$(dr-endpoint) \
	--debug \
	test.json
.PHONY: r/0

## -- test (t) --
$(eval $(call alias, test, t/0))
$(eval $(call alias, t, t/0))

## run tests
t/0:
	$(tt-go) test ./... -run "_U"
.PHONY: t/0

## run unit & int tests
t/all:
	$(tt-go) test ./...
.PHONY: t/all

## -- infra (f) --
$(eval $(call alias, infra, f/0))
$(eval $(call alias, f, f/0))

## alias for f/plan
f/0: f/plan
.PHONY: f/0

## run localstack
f/dev:
	$(tf-dc) up
.PHONY: f/start

## run plan->apply->seed
f/setup: f/plan f/apply f/seed
.PHONY: f/scaffold

## create infra migration plan
f/plan: $(df-tf)
	$(tf-terraform) plan
.PHONY: f

## validate plan
f/valid:
	$(tf-terraform) validate
.PHONY: f/validate

## apply migraition plan
f/apply:
	$(tf-terraform) apply -auto-approve
.PHONY: f/apply

## seed initial state
f/seed:
	$(tr-env) \
	$(tr-aws) dynamodb put-item \
	--table-name $(df-table) \
	--item '{ "Id": {"S": $(df-item-id)}, "Count": {"N": "0"} }' \
	--endpoint-url=$(dr-endpoint) \
	--debug
.PHOYN: f/seed

## destroy infra
f/clean:
	$(tf-terraform) destroy
.PHONY: f/reset

# -- i/helpers
$(df-tf):
	$(tf-terraform) init
