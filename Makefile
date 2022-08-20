include .env
include ./Makefile.base.mk

# -- cosmetics --
help-colw = 8

# -- constants --
df-infra = infra
df-tf = $(df-infra)/.terraform
df-plan = terraform.tfplan
df-table = $(SHARE_COUNT_NAME)
df-item-id = $(SHARE_FILES_NAME)
df-app = ../app
df-app-cfg-dev = $(df-app)/Mise/InfoDev.plist
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
tf-plist = plutil
tb-go = go
tt-go = go
tr-aws = aws
tr-env = AWS_CONFIG_FILE=.aws/config AWS_SHARED_CREDENTIALS_FILE=.aws/creds

# -- state --
sf-url = \
	$(tf-terraform) output share_add_url \
	| tr -d '"' \
	| sed 's/us-east-1.amazonaws.com/localhost.localstack.cloud:4566/'

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
r/0:
	curl \
	$$($(sf-url))
.PHONY: r/0

# r/0:
# 	$(tr-env) \
# 	$(tr-aws) lambda invoke \
# 	--function-name $(dr-fn) \
# 	--payload $$(base64 < $(dr-payload)) \
# 	--endpoint-url=$(dr-endpoint) \
# 	--log-type Tail \
# 	--debug \
# 	test.json
# .PHONY: r/0

## read logs
r/logs:
	$(tr-env) \
	$(tr-aws) logs get-log-events \
	--log-group-name /aws/lambda/$(dr-fn) \
	--log-stream-name $$(/bin/cat out) \
	--limit 5
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

## prepare & run dev stack
f/dev: f/up f/setup f/tail
.PHONY: f/dev

## run localstack
f/up:
	$(tf-dc) up -d
.PHONY: f/up

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
	$(tf-terraform) plan -out=$(df-plan)
.PHONY: f

## apply migration plan
f/apply:
	$(tf-terraform) apply -auto-approve $(df-plan)
.PHONY: f/apply

## validate configuration
f/valid:
	$(tf-terraform) validate
.PHONY: f/validate

## seed initial state
f/seed:
	$(tr-env) \
	AWS_ACCESS_KEY_ID=test \
	AWS_SECRET_ACCESS_KEY=test \
	AWS_PAGER="" \
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

## show the dev share url
f/url:
	echo "$$($(sf-url))"
.PHONY: f/url

## sync the dev share url
f/u/sync:
	# $(tf-plist) \
	# -replace "Share-URL" \
	# -string $$($(sf-url)) \
	# $(df-app-cfg-dev)
.PHONY: f/u/sync

# -- i/helpers
$(df-tf):
	$(tf-terraform) init