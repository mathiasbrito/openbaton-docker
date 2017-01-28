# docker-popd

**docker-popd** is a daemon that implements the Pop protocol and handles a Docker instance. See the builtin helper for further informations.

## How to install `docker-popd`

On both *NIX and Windows:
```shell
go get -u github.com/mcilloni/openbaton-docker/cmd/docker-popd
```
The daemon connects to `dockerd` using the official Docker Go client, and thus it works fine with every platform Docker runs on, including Docker Windows containers and `docker-machine`.

## Initialising a new configuration
**docker-popd** uses a config file in TOML format to configure itself. Use the `init` command to create a new configuration file.

## Quick setup
You can quickly set up a local instance of docker-popd by using it with an empty command line; this will use its default parameters.
The server won't work without a configured user; to use a temporary one on the fly set in your environment `POPD_AUTH="username:password"`. 
