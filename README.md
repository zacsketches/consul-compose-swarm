# Microservice Registry

## Environment
I run single node docker swarm for on-premises microservice deployments in 
systems that do not have a high availability requirement.  This allows
the benefits of DevOps for creating the microservices even though I'm deploying
to a single server that is not in the cloud.  

During the development process our team likes to use Docker compose invoked via
`docker-compose up` to stand up the system for testing.  In deployment we prefer
`docker stack deploy` in order to take advantage of service health monitoring 
that comes with docker swarm.

## Problem
One of our developers reported that there is a service discovery problem when
standing up a collection of services that is caused by naming convention 
differences between swarm and compose. This repo attempts to recreate that 
problem, but after testing with both compose and swarm **it appears that in both
instances the webserver service can successfull connect** to the mongo database
at the hostname `db` with the URI string `mongodb://db:27017`.

Our team reported that both swarm and compose offer a DNS service that supports 
inter-service communication.  However, there are subtle differences in the way 
the two systems create DNS names.  Compose appends `_1` to the end of each 
service name in the compose file and swarm does not append this suffix.  For 
example, this snippet of a docker-compose.yml defines a service named `broker`.

``` docker
version: '3.1'
services:
  broker:
    healthcheck:
      test: "exit 0"
    image:  fathom5/broker:swarm
    networks: 
      - grace
    depends_on: 
      - mongo
    ports:
      - "8001:8001"

    ...
```

If this service and the other services in the file were stood up using 
`docker-compose up` the resulting DNS name for the broker would be 
`gracev0_broker_1`.  Whereas if the system were launched with the command 
`docker stack deploy -c docker-compose.yml grace-v0` then the resulting
DNS name for broker would be `gracev0_broker`.  

**This seemed to be a problem** because it changes the DNS names hard 
coded into all of the other microservices.  The normal answer to avoid hard coding
service DNS names is to provide a service registry in the system so that all of 
the microservices look up the DNS names for servcies they want to communicate with
in that registry.  

Assume for a moment that there was another service, like [Consul](https://www.consul.io/discovery.html) to provide a service registry that was also 
set up in the `docker-compose.yml` file.

```
version: '3.1'
services:
  registry:
    image:  consul
    networks: 
      - grace
    ports:
      - "8500:8500"
```

Then all the other microservices could cosult the registry for `<address>:<port>`
information about the other services. 

**But, how do people get the right address for the registry service?** The registry
will solve the problem of service identification for all other services, but we have
not addressed the underlying problem that *consul* is still going to have two
different DNS names depending on how the mesh was stood up.

   * Compose - `gracev0_consul_1`
   * Swarm - `gracev0_consul`

## Proposed solution
This repo provides a worked example for providing a service registry with Consul and
providing a small package to resolve the registry DNS name.
