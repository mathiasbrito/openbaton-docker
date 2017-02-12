# Pop VIM Driver plugin
`pop-plugind` is a VIM Driver Plugin that enables OpenBaton to use a containerisation solution (like Docker) through a Pop daemon, like [docker-popd][docker-popd].

This plugin is still in its very early phase, and many features are yet to be implemented. 

## How to install `pop-plugind`
`pop-plugind` can be built using `go`, with 

```shell
go get -u github.com/mcilloni/openbaton-docker/cmd/pop-plugind
```

This will fetch and build the source and its dependencies, creating a `pop-plugind` in the `bin` directory of your _GOPATH_.

## How to launch the Docker plugin

```shell
pop-plugind --log "logfile" "name" "rabbit host" "port" "# of workers to be spawned" "username" "password"
```

If `--log` is omitted, the plugin will log on `/var/log/openbaton/type-plugin.log` on *NIX and in the Event Logger on Windows.
Use `--log -` to log only on `stdout`.

An empty command line is equivalent with the invocation below:

```shell
pop-plugind --log "-" openbaton localhost 5672 10 admin openbaton
```

## How to use the Pop plugin

After launching the plugin (see above) and the Pop server, use a VIM instance JSON descriptor similar to this:

```json
{
  "name":"docker-pop-vim-instance",
  "authUrl":"http://localhost:60000",
  "tenant":"tenantName",
  "username":"user_name",
  "password":"password_value",
  "keyPair":"keyName",
  "securityGroups": [
    "securityName"
  ],
  "type":"docker-pop",
  "location":{
    "name":"Berlin",
    "latitude":"52.525876",
    "longitude":"13.314400"
  }
}
```

## How to use `thunks` to make a Jar

OpenBaton expects its plugins to be contained in `.jar` files. If you wish to autolaunch this plugin together with the NFVO, install [thunks] using go, and then invoke

```shell
# set the environment variable GOOS = linux
go build github.com/mcilloni/openbaton-docker/cmd/pop-plugind
thunks pop-plugind
``` 

This will create a self extracting Jar named _pop-plugind.jar_ that will extract and launch the VimDriver, executing it with the arguments provided by the NFVO.

## Issue tracker

Issues and bug reports should be posted to the GitHub Issue Tracker of this project.

## What is Open Baton?

Open Baton is an open source project providing a comprehensive implementation of the ETSI Management and Orchestration (MANO) specification and the TOSCA Standard.

Open Baton provides multiple mechanisms for interoperating with different VNFM vendor solutions. It has a modular architecture which can be easily extended for supporting additional use cases. 

It integrates with OpenStack as standard de-facto VIM implementation, and provides a driver mechanism for supporting additional VIM types. It supports Network Service management either using the provided Generic VNFM and Juju VNFM, or integrating additional specific VNFMs. It provides several mechanisms (REST or PUB/SUB) for interoperating with external VNFMs. 

It can be combined with additional components (Monitoring, Fault Management, Autoscaling, and Network Slicing Engine) for building a unique MANO comprehensive solution.

## Source Code and documentation

The Source Code of the other Open Baton projects can be found [on their GitHub page][openbaton-github], and the documentation can be found [on the official website][openbaton-doc].

## News and Website

Check the [Open Baton Website][openbaton]!

## Licensing and distribution
Licensed under the Apache License. See LICENSE for further details.

[openbaton]: http://openbaton.org
[openbaton-doc]: http://openbaton.org/documentation
[openbaton-github]: http://github.org/openbaton
[docker-popd]: https://github.com/mcilloni/openbaton-docker/tree/master/cmd/docker-popd
[thunks]: https://github.com/mcilloni/thunks
