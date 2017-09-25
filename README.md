# faas-dcos
[DCOS](https://dcos.io/) plugin for [OpenFaas](https://github.com/alexellis/faas) 

## Prerequisites: 
1. a running DCOS cluster accessible without authentication 
1. an external [Marathon-LB](https://dcos.io/docs/1.9/networking/marathon-lb/) service running 

Plugin has been tested with DCOS version 1.9.2.

A quick way to have a cluster running is https://github.com/dcos/dcos-vagrant 

## Installation
Install _OpenFaas_ components with:
```
dcos marathon group add faas-dcos.json
```

After a few minutes, _OpenFaaS_ Interface should be available at `http://<public_node_address>:10012/ui/`.

You can now deploy function with values:
```
image: functions/nodeinfo:latest  
name: nodeinfo  
handler: node main.js
```
or using _faas-cli_ by setting the _gateway_ at `http://<public_node_address>:10012`.

## TODO
1. Handle authentication and token expiration (see [#292](https://github.com/gambol99/go-marathon/issues/292))
1. Marathon 1.5 compatibility (remove go-marathon or [#324](https://github.com/gambol99/go-marathon/issues/324))