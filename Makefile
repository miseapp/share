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
## init dev environment
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

# -- start --
## start localstack
stack: s
.PHONY: start

s:
	$(t-localstack) start
.PHONY: s

## stop localstack
s/stop:
	$(t-docker) stop localstack_main
.PHONY: s/stop

# -- infra --
## plan dev in[f]ra
infra: i
.PHONY: infra

f: $(d-tf)
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
f/destroy:
	$(t-tf) destroy
.PHONY: f/destroy

# -- i/helpers
$(d-tf):
	$(t-tf) init
