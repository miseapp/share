# in your makefile:
#   1. include ./Makefile.base.mk".
#   2. help-colw = <num> # the len of ur longest target
#
# to make stuff show up in the help, write magic comments starting
# w/ "##", for example..
#
# to define a section:
#   ## -- <name> <desc...> --
#
# to define a target:
#   ## <desc>
#   <name>: ...

# -- config --
.DEFAULT_GOAL := help

# -- cosmetics --
rs = \033[0m
ul = \033[4;37m
bd = \033[1;37m
rd = \033[1;31m
gr = \033[0;90m
gl = \033[4;90m

# -- functions --
# creates a phony alias for another target, use like:
# 	$(eval $(call alias, <alias>, <original>))
define alias
$(1): $(2)
.PHONY: $(1)
endef

# -- help --
$(eval $(call alias, h, help))

help:
	@awk "$$HELP" $(MAKEFILE_LIST)
.PHONY: help

define HELP
BEGIN {
	# set props
	k_nsec = 0

	# print header
	print "$(gl)usage$(rs)"
	print "  $(bd)make <target>$(rs)\n"
	print "$(gl)targets$(rs)"
}

# match and print sections:
/^## -- .* --$$/ {
	# parse name/desc
	name = $$3
	nlen = length(name)
	desc = substr($$0, 7 + nlen, length($$0) - 9 - nlen)

	# print section name
	printf "%s  $(ul)%s$(rs)%-*s$(gr) %s$(rs)\n",
		(k_nsec == 0) ? "" : "\n",
		name,
		$(help-colw) - nlen,
		"",
		desc

	# incr section counter
	k_nsec++

	# consume line so next regex doesn't match
	getline
}

# match and print targets
/^## .*$$/ {
	$$1 = ""
	desc = $$0
	getline;
	sub(/:/, "", $$1)
	name = $$1
	printf "  $(rd)%-$(help-colw)s$(gr) %s$(rs)\n", name, desc
}
endef
export HELP