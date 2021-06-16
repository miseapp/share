include ./Makefile.base.mk

# -- cosmetics --
help-column-width = 7

# -- tools --
t-docker = docker
t-localstack = localstack
t-terraform = terraform

# -- consants --
d-terraform = .terraform

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

f: $(d-terraform)
	$(t-terraform) plan
.PHONY: f

## validate infra
f/valid:
	$(t-terraform) validate
.PHONY: f/validate

## apply planned infra
f/apply:
	$(t-terraform) apply
.PHONY: f/apply

## destroy infra
f/destroy:
	$(t-terraform) destroy
.PHONY: f/destroy

# -- i/helpers
$(d-terraform):
	$(t-terraform) init
