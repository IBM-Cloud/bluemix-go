# Invite User

This example invites user to an account.

Example: 

```
go run main.go -userEmail new@in.ibm.com -accountID <account-id> 
```

# Invite user with acces groups

This example invites user to an account with the comma separated list of access groups

Example: 

```
go run main.go -userEmail new@in.ibm.com -accountID <account-id> -accessGroups "<AccessGroupId-****>,<AccessGroupId-***>"
```

# Invite user with IAM policy

This example invites user to an account with an IAM Policy

Example: 

```
go run main.go -userEmail new@in.ibm.com -accountID <account-id> -roles "Opera
tor,Writer" -service "<service>" -resourceGroupID "<resourceGroupID>" 
```

# Invite user with Classic infrastructure Permissions

This example invites user to an account with Comma separated list of classic infrastructure permissions

Example: 

```
go run main.go -userEmail new@in.ibm.com -accountID <account-id> -infraPermission "LOCKBOX_MANAGE,I
P_ADD,FIREWALL_RULE_MANAGE,LOADBALANCER_MANAGE"
```

# Invite user with CloudFoundry roles 

This example invites user to an account with Comma separated list of orgnization and space roles

Example: 

```
go run main.go -userEmail new@in.ibm.com -accountID <account-id> -org <org-name> -space <space-name> -region <region> -orgRoles "BillingManager,Manager" -spaceRoles "Developer,Manager"
```

# Grant a permission to a user from a set of whitelisted IPs to IBM Cloud console 

This example configure a user with list of white-listed IPs from which he/she can access the IBM Cloud console. White-listed ip is a string in which all the IPs are comma seperated.

Example: 

```
go run main.go -userEmail new@in.ibm.com -accountID <account-id> -allowedIP "192.168.0.0,192.168.0.1" -org <org-name> -space <space-name> -region <region> -orgRoles "BillingManager,Manager" -spaceRoles "Developer,Manager"
```



