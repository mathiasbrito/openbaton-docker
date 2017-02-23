openbaton-docker 
================

[![GoDoc](https://godoc.org/github.com/mcilloni/openbaton-docker?status.svg)](https://godoc.org/github.com/mcilloni/openbaton-docker)

This repository hosts several packages and services to enable the [OpenBaton][openbaton] [NFVO][nfvo] to use and orchestrate [Docker containers][docker]. 

## Packages

- [mgmt](https://github.com/mcilloni/openbaton-docker/tree/master/mgmt): provides a fully configurable management protocol via AMQP.
- [pop](https://github.com/mcilloni/openbaton-docker/tree/master/pop): [gRPC] based protocol to handle and administrate a remote Docker instance as a Point-of-Presence.
- [pop/client](https://github.com/mcilloni/openbaton-docker/tree/master/pop/client): provides a full client for Pop, that handles the mapping between Pop and OpenBaton.
- [pop/proto](https://github.com/mcilloni/openbaton-docker/tree/master/pop/proto): Protobuf [gRPC] service definition for Pop. Automatically generated. 
- [docker-pop-server](https://github.com/mcilloni/openbaton-docker/tree/master/pop/server): implements a Pop server that uses a Docker instance as its backend.

## Services
- [cmd/pop](https://github.com/mcilloni/openbaton-docker/tree/master/cmd/pop): CLI client to query Pop daemons.
- [cmd/docker-popd](https://github.com/mcilloni/openbaton-docker/tree/master/cmd/docker-popd): Pop server for Docker.
- [cmd/pop-plugind](https://github.com/mcilloni/openbaton-docker/tree/master/cmd/pop-plugind): OpenBaton Plugin for Pop.
- [cmd/mgmt-gvnfmd](https://github.com/mcilloni/openbaton-docker/tree/master/cmd/mgmt-gvnfmd): Generic VNF Manager, using the mgmt protocol instead of an EMS.

## Issue tracker

Issues and bug reports should be posted to the GitHub Issue Tracker of this project

# What is Open Baton?

Open Baton is an open source project providing a comprehensive implementation of the ETSI Management and Orchestration (MANO) specification and the TOSCA Standard.

Open Baton provides multiple mechanisms for interoperating with different VNFM vendor solutions. It has a modular architecture which can be easily extended for supporting additional use cases. 

It integrates with OpenStack as standard de-facto VIM implementation, and provides a driver mechanism for supporting additional VIM types. It supports Network Service management either using the provided Generic VNFM and Juju VNFM, or integrating additional specific VNFMs. It provides several mechanisms (REST or PUB/SUB) for interoperating with external VNFMs. 

It can be combined with additional components (Monitoring, Fault Management, Autoscaling, and Network Slicing Engine) for building a unique MANO comprehensive solution.

## Source Code and documentation

The Source Code of the other Open Baton projects can be found [here][openbaton-github] and the documentation can be found [here][openbaton-doc] .

## News and Website

Check the [Open Baton Website][openbaton]

## Licensing and distribution
Licensed under the Apache License, Version 2.0. See LICENSE for further details.

[openbaton]: http://openbaton.org
[openbaton-doc]: http://openbaton.org/documentation
[openbaton-github]: http://github.org/openbaton
[nfvo]: https://github.com/openbaton/NFVO
[NFV MANO]:http://docbox.etsi.org/ISG/NFV/Open/Published/gs_NFV-MAN001v010101p%20-%20Management%20and%20Orchestration.pdf
[docker]: http://www.docker.com
[gRPC]: http://grpc.io