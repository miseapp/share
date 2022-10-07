include .env-dev
include ./Makefile.base.mk

# -- cosmetics --
help-colw = 8

# -- data --
ds-endpoint = $(AWS_ENDPOINT)
ds-denv = .env-dev
ds-penv = .env-prod

df-infra = infra
df-tf = $(df-infra)/.terraform
df-plan = terraform.tfplan
df-table = $(SHARE_COUNT_NAME)
df-item-id = $(SHARE_FILES_NAME)
df-app = ../app
df-app-dcfg = $(df-app)/Mise/InfoDev.plist

db-entry = cmd/main.go
db-build = build
db-binary = $(SHARE_ADD_BINARY)
db-archive = $(SHARE_ADD_ARCHIVE)

dr-fn = $(SHARE_ADD_NAME)
dr-input = input.json

# -- tools --
ts-denv = env $$(grep -v "^\#" $(ds-denv) | xargs)
ts-penv = env $$(grep -v "^\#" $(ds-penv) | xargs)
ts-aws-d = $(ts-denv) aws --endpoint $(ds-endpoint)
ts-aws-p = $(ts-penv) aws

ti-brew = brew

tf-dc = docker-compose --env-file=$(ds-denv)
tf-d = $(ts-denv) terraform -chdir="$(df-infra)"
tf-p = $(ts-penv) terraform -chdir="$(df-infra)"
tf-plist = plutil

td-go = GOOS=linux GOARCH=amd64 go
tb-d-go = $(ts-denv) $(td-go)
tb-p-go = $(ts-penv) $(td-go)

tt-go = go
tt-d-go = $(ts-denv) $(tt-go)
tt-p-go = $(ts-penv) $(tt-go)

tr-http = http

# -- state --
sf-url = $(tf-d) output -raw share_add_url

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

## build fn
b/0:
	$(tb-d-go) build -o $(db-binary) $(db-entry)
.PHONY: b/0

## build & archive
b/arch: b
	zip $(db-archive) $(db-binary)
.PHONY: b/arch

## build prod
b/p:
	$(tb-p-go) build -o $(db-binary) $(db-entry)
.PHONY: b/prod

## build & archive prod
b/arch/p: b/p
	$(tb-p-go) build -o $(db-binary) $(db-entry)
.PHONY: b/p/arch

## clean build dir
b/clean:
	rm -rf $(db-build)
.PHONY: b/clean

## -- run (r) --
$(eval $(call alias, run, r/0))
$(eval $(call alias, r, r/0))

## call local handler fn
r/0:
	$(tr-http) POST $$($(sf-url)) < $(dr-input)
.PHONY: r/0

## read logs
r/logs:
	$(ts-aws-d) \
	logs tail \
	/aws/lambda/$(dr-fn)
.PHONY: r/logs

## -- test (t) --
$(eval $(call alias, test, t/0))
$(eval $(call alias, t, t/0))

## run tests
t/0:
	$(tt-d-go) test ./... -run "_U"
.PHONY: t/0

## run unit & int tests
t/all:
	$(tt-d-go) test ./...
.PHONY: t/all

## -- infra (f) --
$(eval $(call alias, infra, f/0))
$(eval $(call alias, f, f/0))

## alias for f/dev
f/0: f/dev
.PHONY: f/0

## prepare & run dev stack
f/dev: f/up f/setup f/tail
.PHONY: f/dev

## prepare & run dev stack w/ debug output
f/dbg: f/upv f/setup f/tail
.PHOYN: f/dbg

## run localstack
f/up:
	$(tf-dc) up -d
.PHONY: f/up

## run localstack w/ debug output
f/upv:
	LS_LOG=trace $(tf-dc) up -d
.PHONY: f/upv

## run localstack in foreground
f/upf:
	$(tf-dc) up
.PHONY: f/upf

## tail localstack logs
f/tail:
	$(tf-dc) logs -f -t
.PHONY: f/tail

## stop localstack
f/down:
	$(tf-dc) down
.PHONY: f/down

## run plan->apply->seed
f/setup: f/update f/seed
.PHONY: f/scaffold

## run plan->apply
f/update: f/plan f/apply f/u/sync
.PHONY: f/update

## create migration plan
f/plan: $(df-tf)
	$(tf-d) plan -out=$(df-plan)
.PHONY: f

## apply migration plan
f/apply:
	$(tf-d) apply -auto-approve $(df-plan)
.PHONY: f/apply

## validate configuration
f/valid:
	$(tf-d) validate
.PHONY: f/validate

## seed initial state
f/seed:
	$(ts-aws-d) \
	dynamodb put-item \
	--table-name $(df-table) \
	--item '{ "Id": {"S": $(df-item-id)}, "Count": {"N": "0"} }'
.PHONY: f/seed

## list tables
f/tables:
	$(ts-aws-d) \
	dynamodb list-tables
.PHONY: f/tables

## destroy infra
f/clean:
	$(tf-d) destroy
.PHONY: f/reset

## show the dev share url
f/url:
	echo "$$($(sf-url))"
.PHONY: f/url

## sync the dev share url
f/u/sync:
	$(tf-plist) \
	-replace "Share-URL" \
	-string $$($(sf-url)) \
	$(df-app-dcfg)
.PHONY: f/u/sync

# -- i/helpers
$(df-tf):
	$(tf-d) init