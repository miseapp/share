share
---
for sharing [mise](https://miseapp.co) recipes

# setup

install system dependencies

```sh
$ make init
```

# dev

run the localstack server

```sh
$ make f/start
$ make f/stop # stop it later
```

build the fn if necessary

```sh
$ make f/build
```

create (or update) the infrastructure

```sh
$ make f/plan  # make infra plan, check it for ur changes
$ make f/apply # apply the plan
$ make f/seed  # seed initial state
```

to just do all of that at once

```sh
$ make f/setup
```