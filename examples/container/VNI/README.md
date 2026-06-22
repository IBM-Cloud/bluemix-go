# VNI Integration Test Tool

This tool allows you to test the VNI (Virtual Network Interface) GraphQL operations against a real IBM Cloud environment.

## Prerequisites

1. **IBM Cloud Account** with access to Kubernetes Service
2. **Bare Metal Cluster** with at least one worker node
3. **VPC VNI** created and available (not attached to any resource)
4. **API Key** with appropriate permissions
5. **Go 1.21+** installed

## Setup

### 1. Set Environment Variables

```bash
export IBMCLOUD_API_KEY="your-api-key-here"
export IBMCLOUD_ACCOUNT_ID="your-account-id-here"
```

### 2. Build the Test Tool

```bash
cd ../bluemix-go/examples/container/VNI
go build -o vni-test main.go
```

## Usage

### List VNI Attachments

List all VNI attachments for a cluster:

```bash
./vni-test -op=list -cluster=<cluster-id>
```

List VNI attachments for a specific worker:

```bash
./vni-test -op=list -worker=<worker-id>
```

With pagination:

```bash
./vni-test -op=list -cluster=<cluster-id> -first=5 -after=<cursor>
```

### Attach VNI to Bare Metal Worker

Attach to a specific worker:

```bash
./vni-test -op=attach \
  -worker=<worker-id> \
  -vni=<vni-id> \
  -vlan=100
```

Attach to any available worker in a cluster:

```bash
./vni-test -op=attach \
  -cluster=<cluster-id> \
  -vni=<vni-id> \
  -vlan=100
```

With auto-delete enabled:

```bash
./vni-test -op=attach \
  -worker=<worker-id> \
  -vni=<vni-id> \
  -vlan=100 \
  -auto-delete
```

### Detach VNI from Worker

```bash
./vni-test -op=detach \
  -worker=<worker-id> \
  -vni=<vni-id>
```

With auto-delete (deletes the VNI after detaching):

```bash
./vni-test -op=detach \
  -worker=<worker-id> \
  -vni=<vni-id> \
  -auto-delete
```

### Get Specific VNI Attachment

```bash
./vni-test -op=get \
  -cluster=<cluster-id> \
  -worker=<worker-id> \
  -vni=<vni-id>
```

## Complete Test Workflow

Here's a complete workflow to test all operations:

```bash
# 1. Set your credentials
export IBMCLOUD_API_KEY="your-api-key"
export IBMCLOUD_ACCOUNT_ID="your-account-id"

# 2. List existing attachments (should be empty initially)
./vni-test -op=list -cluster=<your-cluster-id>

# 3. Attach a VNI to a worker
./vni-test -op=attach \
  -worker=<your-worker-id> \
  -vni=<your-vni-id> \
  -vlan=100

# 4. List attachments again (should show the new attachment)
./vni-test -op=list -cluster=<your-cluster-id>

# 5. Get the specific attachment
./vni-test -op=get \
  -cluster=<your-cluster-id> \
  -worker=<your-worker-id> \
  -vni=<your-vni-id>

# 6. Detach the VNI
./vni-test -op=detach \
  -worker=<your-worker-id> \
  -vni=<your-vni-id>

# 7. Verify it's gone
./vni-test -op=list -cluster=<your-cluster-id>
```

## Getting Required IDs

### Get Cluster ID

```bash
ibmcloud ks clusters
```

### Get Worker ID

```bash
ibmcloud ks workers --cluster <cluster-id>
```

### Get VNI ID

```bash
ibmcloud is virtual-network-interfaces
```

Or create a new VNI:

```bash
ibmcloud is virtual-network-interface-create \
  --name test-vni \
  --subnet <subnet-id> \
  --vpc <vpc-id>
```

## Command-Line Flags

| Flag | Description | Required | Default |
|------|-------------|----------|---------|
| `-apikey` | IBM Cloud API Key | Yes (or env var) | `$IBMCLOUD_API_KEY` |
| `-account` | IBM Cloud Account ID | Yes (or env var) | `$IBMCLOUD_ACCOUNT_ID` |
| `-region` | IBM Cloud region | No | `us-south` |
| `-op` | Operation: attach, detach, list, get | Yes | `list` |
| `-cluster` | Cluster ID or name | Depends on operation | - |
| `-worker` | Worker ID | Depends on operation | - |
| `-vni` | VNI ID (e.g., r006-xxx) | For attach/detach/get | - |
| `-vlan` | VLAN ID (1-500) | For attach | `100` |
| `-auto-delete` | Auto-delete VNI on detach | No | `false` |
| `-first` | Number of items to fetch | For list | `10` |
| `-after` | Pagination cursor | For list | - |

## Expected Output

### Successful Attach

```
Attaching VNI r006-xxx to worker kube-xxx with VLAN ID 100...

✅ VNI attached successfully!
Worker ID: kube-xxx
VNI ID: r006-xxx
VNI Name: test-vni
VLAN ID: 100
Primary IP: 10.240.0.5
MAC Address: 02:00:00:00:00:01
```

### Successful List

```
Listing VNI attachments for cluster cluster-xxx...

Attachable Type: KubernetesCluster
Total Attachments: 2

VNI Attachments:
================

1. Worker: kube-worker-1
   VNI ID: r006-vni-1
   VNI Name: test-vni-1
   VLAN ID: 100
   Primary IP: 10.240.0.5
   MAC Address: 02:00:00:00:00:01

2. Worker: kube-worker-2
   VNI ID: r006-vni-2
   VNI Name: test-vni-2
   VLAN ID: 101
   Primary IP: 10.240.0.6
   MAC Address: 02:00:00:00:00:02
```

## Troubleshooting

### Error: "VNI not found"

- Verify the VNI ID is correct
- Ensure the VNI exists and is not already attached to another resource
- Check that the VNI is in the same VPC as the cluster

### Error: "Worker not found"

- Verify the worker ID is correct
- Ensure the worker is in a ready state
- Check that you have access to the cluster

### Error: "VLAN ID already in use"

- VLAN IDs must be unique per subnet
- Use a different VLAN ID or ensure VNIs are from the same subnet

### Error: "VNI attachments are only supported for bare metal workers"

- This feature only works with bare metal workers
- Virtual workers are not supported in this phase

## Debug Mode

Enable detailed HTTP request/response logging:

```bash
export TRACE_LEVEL=true
./vni-test -op=list -cluster=<cluster-id>
```

## Testing Checklist

- [ ] List attachments for empty cluster
- [ ] Attach VNI to specific worker
- [ ] List attachments shows new attachment
- [ ] Get specific attachment
- [ ] Attach VNI to cluster (any worker)
- [ ] Test VLAN ID reuse with VNIs from same subnet
- [ ] Test pagination with multiple attachments
- [ ] Detach VNI without auto-delete
- [ ] Detach VNI with auto-delete
- [ ] Verify error handling for invalid inputs

## Next Steps

After successful testing:

1. Document any issues or edge cases discovered
2. Update the implementation if needed
3. Create PR for bluemix-go repository
4. Proceed with terraform-provider-ibm implementation
