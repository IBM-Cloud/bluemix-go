# Authorization Policy

This example creates an authorization to allow a service instance access to another service instance.

Details of the API function implemented can be found in the IBM CLoud API docs: https://cloud.ibm.com/apidocs/iam-policy-management#create-a-policy

Example:

```
go run main.go -org new@in.ibm.com -source_service_name "cloud-object-storage" -target_service_name kms -roles Reader -source_service_instance_id 123123 -target_service_instance_id 456456
```

# Create authorization policy with service instance ID

```
go run main.go -org new@in.ibm.com -source_service_name <source_service_name> -target_service_name <target_service_name> -roles Reader -source_service_instance_id <source_service_instance_id> -target_service_instance_id <target_service_instance_id>
```

# Create authorization policy with resource group ID

```
go run main.go -org new@in.ibm.com -source_service_name <source_service_name> -target_service_name <target_service_name> -roles Reader -source_resource_group_id <source_resource_group_id> -target_resource_group_id <target_resource_group_id>
```