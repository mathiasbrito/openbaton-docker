# mgmt-based Generic VNFM (AMQP)
`mgmt-gvnfmd` is a generic Virtual Network Function Manager implementation for [OpenBaton][openbaton], written using Go and [go-openbaton] and designed to use the [mgmt] protocol to manage Virtual Network Functions, instead of an Element Management System.

## Technical Requirements

You need a fully working NFVO to use the Pop VNFM, plus a mgmt-enabled consumer, like the one provided by [pop-plugind].

## How to install `mgmt-gvnfmd`

On both *NIX and Windows:
```shell
go get -u github.com/mcilloni/openbaton-docker/cmd/mgmt-gvnfmd
```

The `go` tool will automatically fetch and build both the sources and their dependencies, and a `mgmt-gvnfmd` binary will be generated in `$GOPATH/bin` (`%GOPATH%\bin` on Windows CMD).

## How to use `mgmt-gvnfmd`

 ```bash
 mgmt-gvnfmd --cfg /path/to/config.toml
 ```

The VNFM must be configured using a configuration file, specified through the `--cfg` flag (see [the sample configuration for further details][sample-conf]).

In case no such flag is specified, the default behaviour is to search in the current directory for a file named `config.toml`.

## How to configure `mgmt-gvnfmd`

The sample configuration should work straight out of the box with a standard local setup of OpenBaton.

`vnfm.logger.level` can be used to change the default verbosity of the logger, choosing a value between `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL` and `PANIC`.

## Issue tracker

Issues and bug reports should be posted to the GitHub Issue Tracker of this project.

## What is Open Baton?

Open Baton is an open source project providing a comprehensive implementation of the ETSI Management and Orchestration (MANO) specification and the TOSCA Standard.

Open Baton provides multiple mechanisms for interoperating with different VNFM vendor solutions. It has a modular architecture which can be easily extended for supporting additional use cases. 

It integrates with OpenStack as its standard de-facto VIM implementation, and provides a driver mechanism for supporting additional VIM types. It supports Network Service management either using the provided Generic VNFM and Juju VNFM, or integrating additional specific VNFMs. It provides several mechanisms (REST or PUB/SUB) for interoperating with external VNFMs. 

It can be combined with additional components (Monitoring, Fault Management, Autoscaling, and Network Slicing Engine) for building a unique MANO comprehensive solution.

## Source Code and documentation

The Source Code of the other Open Baton projects can be found [on their GitHub page][openbaton-github], and the documentation can be found [on the official website][openbaton-doc].

## News and Website

Check the [Open Baton Website][openbaton]

## Licensing and distribution
Apache License, Version 2.0. See LICENSE for further details.

[openbaton]: http://openbaton.org
[openbaton-doc]: http://openbaton.org/documentation
[openbaton-github]: http://github.org/openbaton
[sample-conf]: https://raw.githubusercontent.com/mcilloni/openbaton-docker/master/cmd/mgmt-gvnfmd/config.toml
[go-openbaton]: http://github.com/openbaton/go-openbaton
[mgmt]: https://github.com/mcilloni/openbaton-docker/tree/master/mgmt
[pop-plugind]: https://github.com/mcilloni/openbaton-docker/tree/master/cmd/pop-plugind