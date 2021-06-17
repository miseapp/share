include ./Makefile.base.mk

# -- cosmetics --
help-column-width = 7

# -- constants --
d-infra = infra
d-tf = $(d-infra)/.terraform

# -- tools --
t-docker = docker
t-localstack = localstack
t-tf = terraform -chdir="$(d-infra)"

# -- init --
## setup dev env
init: i
.PHONY: init

i: i/pre
	brew bundle -v
.PHONY: i

i/pre:
ifeq ("$(shell command -v brew)", "")
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
	GOOS=linux go build cmd/...
.PHONY: b

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
	$(t-localstack) start
.PHONY: f/start

## stop localstack
f/stop:
	$(t-docker) stop localstack_main
.PHONY: f/stop

## plan dev infra
f/plan: $(d-tf)
	$(t-tf) plan
.PHONY: f

## validate infra
f/valid:
	$(t-tf) validate
.PHONY: f/validate

## apply planned infra
f/apply:
	$(t-tf) apply
.PHONY: f/apply

## destroy infra
f/reset:
	$(t-tf) destroy
.PHONY: f/reset

# -- i/helpers
$(d-tf):
	$(t-tf) init
