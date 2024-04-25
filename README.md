# terraform-cloud-metrics

This repository contains a simple project which will connect to the terraform-cloud API and return details of the number of projects/workspaces we have.



## Build & Launch

Build the application via the provided [Dockerfile](Dockerfile)

```sh
docker build -t terraform-cloud-metrics:latest .

docker run -p 8090:8090 -e TFE_TOKEN=xxx terraform-cloud-metrics
```



## Sample Output

Access the metrics on http://127.0.0.0.1:8090/metrics and you should see:

```
..
terraform_cloud_workspace_count{project="aws-iam"} 123
terraform_cloud_workspace_count{project="beholder"} 1
terraform_cloud_workspace_count{project="container-registries"} 1
terraform_cloud_workspace_count{project="default"} 0
terraform_cloud_workspace_count{project="doltsoevsky"} 1
terraform_cloud_workspace_count{project="duckpond"} 3
terraform_cloud_workspace_count{project="game5"} 23
terraform_cloud_workspace_count{project="generic"} 9
terraform_cloud_workspace_count{project="lighthouse"} 4
terraform_cloud_workspace_count{project="marine"} 28
terraform_cloud_workspace_count{project="sandbox"} 0
terraform_cloud_workspace_count{project="techservices"} 3
```
