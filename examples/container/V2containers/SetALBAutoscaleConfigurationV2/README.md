# SetALBAutoscaleConfiguration example

This example shows how to set the ALB autoscale config.

Example: 

```
go run main.go -clusterNameOrID mycluster -albID public-crck9aaedd0p8vjmqa0asg-alb1 -minReplicas 2 -maxReplicas 4 -cpuAverageUtilization 600
```