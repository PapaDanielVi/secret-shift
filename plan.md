# SecretShift

secret-shift is a cli tool for migrating/syncing secrets/env from sources like gitlab/github/etcd/hashicorp vault to
gitlab/github/etcd/hashicorp vault/file/kuebernetes secret and configmap.

this cli follows pipeline concept to read and process and write secrets/envs.

pipeline flow: read -> process -> write

this cli provides two way to do it's purpose:

1. a TUI flow that will ask user information step by step
2. json config file that has all configs in it and user will pass it's config using -c flag.

## steps details

### read

    user will provide url and credentials of data source
    like repo (id or url), url, github token, gitlab token, vault username and password etc...

### process

    user can provide some changes to secrets/env in process step
    like (add prefix, add suffix, regex match for including anf excluding, type for including anf excluding, etc...)

### write

 user will provide url and credentials of destination. like repo (id or url), url, github token, gitlab token, vault username and password. also some other flags to handle conflicts, merges, replaces and report already exists

## JSON file config flow

    this way user will provide a path to json config file.
    this flow will allow user to also use ths cli in cronjobs to syn their secrets/env on time bases.
    also user can pass --periodically and --frequency 10m to tell the cli that I want you to sync my secrets/envs every 10 minutes.

## TUI flow

    this flow we will use bubbletea to create interactive terminal ui and then guide user to provide credentials and
    information and then execute the flow at the end.
    we must also suggest main github.com and gitlab.com and custom server to their repo url. like gh cli.

## Must haves

    we must also support environment variable for configurations. so user can send env instead of passing args to cli.
    like SECRET_SHIFT_GITHUB_TOKEN, SECRET_SHIFT_GITLAB_TOKEN, SECRET_SHIFT_VAULT_USERNAME, etc...
