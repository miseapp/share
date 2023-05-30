include .env.dev
include ./Makefile.base.mk

# -- cosmetics --
help-colw = 10

# -- data --
ds-env-d = .env.dev
ds-env-p = .env.prod
ds-crd-d = .creds.dev
ds-crd-p = .creds.prod

df-infra = infra/envs
df-table = $(SHARE_COUNT_NAME)
df-files = $(SHARE_FILES_NAME)
df-item-id = $(df-files)
df-app = ../app
df-app-dcfg = $(df-app)/Mise/InfoDev.plist
df-plan = terraform.tfplan

db-entry = cmd/main.go
db-build = build
db-binary = $(SHARE_ADD_BINARY)
db-archive = $(SHARE_ADD_ARCHIVE)

dr-fn = $(SHARE_ADD_NAME)
dr-input = input.json

# -- tools --
ts-env-d = env $$(grep -h -v "^\#" $(ds-crd-d) $(ds-env-d) | xargs)
ts-env-p = env $$(grep -h -v "^\#" $(ds-crd-p) $(ds-env-p) | xargs)
ts-aws-d = $(ts-env-d) aws --endpoint $(LOCAL_URL)
ts-aws-p = $(ts-env-p) aws

ti-brew = brew

tf-dc = docker-compose --env-file="$(ds-crd-d)" --env-file="$(ds-env-d)"
tf-d = $(ts-env-d) terraform -chdir="$(df-infra)/dev"
tf-p = $(ts-env-p) terraform -chdir="$(df-infra)/prod"
tf-plist = plutil

td-go = GOOS=linux GOARCH=amd64 go
tb-d-go = $(ts-env-d) $(td-go)
tb-p-go = $(ts-env-p) $(td-go)

tt-go = go
tt-d-go = $(ts-env-d) $(tt-go)
tt-p-go = $(ts-env-p) $(tt-go)

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

## build fn [prod]
b/p:
	$(tb-p-go) build -o $(db-binary) $(db-entry)
.PHONY: b/prod

## clean build dir
b/clean:
	rm -rf $(db-build)
.PHONY: b/clean

## -- archive (a) --
$(eval $(call alias, archive, a/0))
$(eval $(call alias, a, a/0))

## build & archive
a/0: b/clean b
	zip $(db-archive) $(db-binary)
.PHONY: a/0

## build & archive [prod]
a/p: b/clean b/p
	zip $(db-archive) $(db-binary)
.PHONY: a/p

## -- run (r) --
$(eval $(call alias, run, r/0))
$(eval $(call alias, r, r/0))

## call local handler fn
r/0:
	$(tr-http) --verbose POST $$($(sf-url)) < $(dr-input)
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

## run infra
f/dev: f/up f/setup f/tail
.PHONY: f/dev

## run infra w/ debug output
f/dbg: f/upv f/setup f/tail
.PHOYN: f/dbg

## run ls container
f/up:
	$(tf-dc) up -d
.PHONY: f/up

## run ls container w/ debug output
f/upv:
	LS_LOG=trace $(tf-dc) up -d
.PHONY: f/upv

## run ls container in foreground
f/upf:
	$(tf-dc) up
.PHONY: f/upf

## tail ls logs
f/tail:
	$(tf-dc) logs -f -t
.PHONY: f/tail

## stop ls
f/down:
	$(tf-dc) down
.PHONY: f/down

## plan, apply, seed
f/setup: f/update f/seed
.PHONY: f/setup

## plan, apply, seed [prod]
f/setup/p: f/update/p f/seed/p
.PHONY: f/setup/p

## plan, apply
f/update: f/plan f/apply f/u/sync
.PHONY: f/update

## plan, apply [prod]
f/update/p: f/plan/p f/apply/p
.PHONY: f/update/p

## make migration plan
f/plan: $(df-infra)/dev/.terraform a
	$(tf-d) plan -out=$(df-plan)
.PHONY: f/plan

# f/plan/p: $(df-infra)/prod/.terraform
## make migration plan [prod]
f/plan/p: a/p
	$(tf-p) plan -out=$(df-plan)
.PHONY: f/plan/p

## apply migration plan
f/apply:
	$(tf-d) apply -auto-approve $(df-plan)
.PHONY: f/apply

## apply migration plan [prod]
f/apply/p:
	$(tf-p) apply $(df-plan)
.PHONY: f/apply/p

## validate config
f/valid:
	$(tf-d) validate
.PHONY: f/validate

## validate config [prod]
f/valid/p:
	$(tf-p) validate
.PHONY: f/valid/p

## seed initial state
f/seed:
	$(ts-aws-d) \
	dynamodb put-item \
	--table-name $(df-table) \
	--item '{ "Id": {"S": $(df-item-id)}, "Count": {"N": "0"} }'
.PHONY: f/seed

## seed initial state [prod]
f/seed/p:
	$(ts-aws-p) \
	dynamodb put-item \
	--table-name $(df-table) \
	--item '{ "Id": {"S": $(df-item-id)}, "Count": {"N": "0"} }'
.PHONY: f/seed/p

## destroy infra
f/clean:
	$(tf-d) destroy
.PHONY: f/clean

## list tables
f/tables:
	$(ts-aws-d) \
	dynamodb list-tables
.PHONY: f/tables

## list files in s3
f/files:
	$(ts-aws-d) \
	s3 ls \
	s3://share-files
.PHONY: f/tables

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
$(df-infra)/dev/.terraform:
	$(tf-d) init

# $(df-infra)/prod/.terraform:
# 	$(tf-p) init