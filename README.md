# faas-dcos
[DC/OS](https://dcos.io/) plugin for [OpenFaas](https://github.com/openfaas/faas) 

## Prerequisites: 
1. a running DC/OS cluster accessible without authentication 
1. an external [Marathon-LB](https://dcos.io/docs/1.9/networking/marathon-lb/) service running
1. DC/OS CLI installed and configured (installation instructions can be found [here](https://dcos.io/docs/1.9/cli/install/) or in the Dashboard top-left corner)

A quick way to have a local DCOS cluster running is https://github.com/dcos/dcos-vagrant. Do not forget to:
1. add option ```oauth_enabled: 'false'``` in _etc/config-1.9.yaml_ config file to disable authentication 
1. set environment variable ```export DCOS_VERSION=1.9.2``` in shell before running _vagrant up_.

Marathon-LB can be easily installed from DCOS Universe packages (the default configuration is okay, just click on INSTALL button):

![Marathon-LB in Universe](docs/images/mlb.png?raw=true "Marathon-LB in Universe")

Plugin has been tested with DC/OS version 1.9.2.

## Installation

Once you have your cluster running, you can easily install _OpenFaas_ components. From _faas-dcos_ project root run the following command :
```
dcos marathon group add faas-dcos.json
```

You should see services being deployed and, after a few minutes, you should have something like this:

![OpenFaas running](docs/images/install.png?raw=true "OpenFaas running")

_OpenFaaS_ Interface should be now be available at `http://<public_node_address>:10012/ui/` where _<public_node_address>_ is the cluster node accessible from outside, is the one running Marathon-LB by the way...

You can now deploy functions using the web interface, for instance using these values:
```
image: functions/nodeinfo:latest  
name: nodeinfo  
handler: node main.js
```

or alternatively you can use the [CLI for OpenFaaS](https://github.com/openfaas/faas-cli) with the following YAML stack file (```functions.yaml```):

```yaml
provider:
  name: faas
  gateway: http://<public_node_address>:10012

functions:
  nodeinfo:
    handler: node main.js
    image: functions/nodeinfo:latest
```
Remember to set the _gateway_ to `http://<public_node_address>:10012`!

Then run this command to deploy your function:

```
$ faas-cli deploy -f ./functions.yml
```

Once the function has been created, you should see a new service running in DC/OS

![Function running](docs/images/function.png?raw=true "Function running")

and it will be available to be executed

![Function invoked](docs/images/invoke.png?raw=true "Function invoked")

## TODO
1. Handle authentication and token expiration (see [#292](https://github.com/gambol99/go-marathon/issues/292))
1. Marathon 1.5 compatibility (remove go-marathon or [#324](https://github.com/gambol99/go-marathon/issues/324))
1. Functions memory limits cannot be configured (see [faas#239](https://github.com/openfaas/faas/issues/239))
