#App example

This example shows how to perform CRUD operation on cloud foundry app.

This creates a app defined in specified space.After successfull creation it perform update and deletes app.

Example: 

```
go run main.go -org example.com -space test

run main.go  -org example.com  -space test -buildpack nodejs_buildpack -name testapp -routeName testapp -path <path to zip file>

```




