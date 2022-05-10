# Cluster example V2

This example shows how to create an Kubernetes cluster with boot volume encryption enabled.

Example: 

```
go run main.go -kmsid <kmsinstanceid> -crkid <rootkeyinyourkms> -vpcid <vpcid> -subnetid <subnetid>
```