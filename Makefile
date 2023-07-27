env := $(or $(ENV), dev)
env-file := .env.$(env)
crd-file := .creds.$(env)

include $(env-file)
include ./Makefile.base.mk

# -- cosmetics --
help-colw = 10

# -- data --
df-infra = infra/envs/$(env)
df-table = $(SHARE_COUNT_NAME)
df-files = $(SHARE_FILES_NAME)
df-item-id = $(df-files)
df-app = ../app
df-app-dcfg = $(df-app)/Mise/InfoDev.plist
df-plan = terraform.tfplan

db-builds = build/bin
db-archives = build/archive
db-entry = cmd/main.go
db-binary = $(SHARE_ADD_BINARY)
db-archive = $(SHARE_ADD_ARCHIVE)

dr-fn = $(SHARE_ADD_NAME)
dr-input = input.json

# -- tools --
ts-env = env $$(grep -h -v "^\#" $(crd-file) $(env-file) | xargs)
ts-aws = $(ts-env) aws $(if $(LOCAL_URL),--endpoint $(LOCAL_URL),)

ti-brew = brew

tf-dc = docker-compose --env-file="$(crd-file)" --env-file="$(env-file)"
tf-tf = $(ts-env) terraform -chdir="$(df-infra)"
tf-plist = plutil

tb-go = $(ts-env) GOOS=linux GOARCH=amd64 go
tt-go = $(ts-env) go

tr-http = http

# -- state --
sf-url = $(tf-tf) output -raw share_add_url

# -- helpers --
dev-only:
ifneq ($(env), dev)
	$(info ✘ this can only be run in dev)
	$(error 1)
endif
.PHONY: dev-only

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
ifeq ($(shell command -v $(ti-brew)),)
	$(info ✘ brew is not installed, please see:)
	$(info - https://brew.sh)
	$(error 2)
endif
.PHONY: i/pre

## -- build (b) --
$(eval $(call alias, build, b/0))
$(eval $(call alias, b, b/0))

## build fn
b/0: $(db-builds)
	$(tb-go) build -o $(db-builds)/$(db-binary) $(db-entry)
.PHONY: b/0

## clean build dir
b/clean:
	rm -rf $(db-builds)
.PHONY: b/clean

$(db-builds):
	mkdir -p $(db-builds)

## -- archive (a) --
$(eval $(call alias, archive, a/0))
$(eval $(call alias, a, a/0))

## build & archive
a/0: $(db-archives) b/clean b
	rm -f $(db-archives)/$(db-archive)
	cd $(db-builds) && zip $(db-archive) $(db-binary)
	mv $(db-builds)/$(db-archive) $(db-archives)/$(db-archive)
.PHONY: a/0

$(db-archives):
	mkdir -p $(db-archives)

## -- run (r) --
$(eval $(call alias, run, r/0))
$(eval $(call alias, r, r/0))

## call local handler fn
r/0:
	$(tr-http) --verbose POST $$($(sf-url)) < $(dr-input)
.PHONY: r/0

## read logs
r/logs:
	$(ts-aws) \
	logs tail \
	/aws/lambda/$(dr-fn)
.PHONY: r/logs

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

## alias for f/dev
f/0: f/dev
.PHONY: f/0

## run infra
f/dev: dev-only f/up f/setup f/tail
.PHONY: f/dev

## run infra w/ debug output
f/dbg: dev-only f/upv f/setup f/tail
.PHOYN: f/dbg

## run ls container
f/up: dev-only
	$(tf-dc) up -d
.PHONY: f/up

## run ls container w/ debug output
f/upv: dev-only
	LS_LOG=trace $(tf-dc) up -d
.PHONY: f/upv

## run ls container in foreground
f/upf: dev-only
	$(tf-dc) up
.PHONY: f/upf

## tail ls logs
f/tail: dev-only
	$(tf-dc) logs -f -t
.PHONY: f/tail

## stop ls
f/down: dev-only
	$(tf-dc) down
.PHONY: f/down

## init tf provider
f/init:
	$(tf-tf) init
.PHONY: f/init

$(df-infra)/.terraform:
	$(tf-tf) init

## plan, apply, seed
f/setup: dev-only f/update f/seed
.PHONY: f/setup

## plan, apply, sync url
f/update: dev-only f/plan/a f/apply f/u/sync
.PHONY: f/update

## build plan
f/plan: $(df-infra)/.terraform
	$(tf-tf) plan -out=$(df-plan)
.PHONY: f/plan

## build plan & archive
f/plan/a: $(df-infra)/.terraform a f/plan
.PHONY: f/plan

## apply plan
f/apply:
	$(tf-tf) apply $(TF_APPLY_OPT) $(df-plan)
.PHONY: f/apply

## validate config
f/valid:
	$(tf-tf) validate
.PHONY: f/valid

## seed initial state
f/seed:
	$(ts-aws) \
	dynamodb put-item \
	--table-name $(df-table) \
	--item '{ "Id": {"S": $(df-item-id)}, "Count": {"N": "0"} }'
.PHONY: f/seed

## destroy infra
f/clean: dev-only
	$(tf-tf) destroy
.PHONY: f/clean

## list tables
f/tables:
	$(ts-aws) \
	dynamodb list-tables
.PHONY: f/tables

## list files in s3
f/files:
	$(ts-aws) \
	s3 ls \
	s3://share-files
.PHONY: f/tables

## show the dev share url
f/url: dev-only
	echo "$$($(sf-url))"
.PHONY: f/url

## sync the dev share url
f/u/sync: dev-only
	$(tf-plist) \
	-replace "Share-URL" \
	-string $$($(sf-url)) \
	$(df-app-dcfg)
.PHONY: f/u/sync