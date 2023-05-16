# Polygon Client

## Description

The Polygon Client is a Go application that retrieves information from the Polygon network using the Polygon RPC endpoint. It periodically fetches the latest block number and latest block hash and logs the information to the console. A makefile is provided to facilitate usage for a developer, it has `build` `test` and other developer-centric makefile targets.

# Prerequisites

To run this application, you need to have the following:

- Docker installed
- Polygon RPC endpoint URL (hardcoded for now)

## Installation

- Clone the repository to your local machine
- Navigate to the root directory of the project
- Run the following command to build the Docker image:

```sh
  docker build -t polygon-client .
```

- Run the following command to start the container:

```sh
 docker run --rm -it polygon-client
```

If you don't want to use docker:

```sh
make build
./bin/polygon-client
```

## Configuration

There is no configuration required for this application as the Polygon RPC endpoint is hardcoded. Simply run the executable and it will periodically fetch the latest block information and print it to the console.

## Improvements

In its current state, the Polygon Client is limited in that it only supports a single RPC endpoint that is hardcoded into the application. In order to make this application more flexible and reusable, it would be beneficial to add support for configuration files or command line arguments that would allow the user to specify the RPC endpoint they wish to use. This would enable the application to be used with a wider range of Polygon-compatible RPC endpoints.
Another improvement could be on Terraform, an alternative approach would be to deploy it to a Kubernetes Cluster, maybe also generate an Helm chart for this application and implement a semantic release CI workflow. There are some examples on my github on how to do the above so they can be omitted here.

## Usage

```bash
go run main.go
```

Once the Docker container is running, the application will periodically log the latest block number and hash to the console. To stop the application, use Ctrl + C.

# CI/CD

- **Test**: The "test" workflow is triggered by a push event and runs on the latest version of Ubuntu. It performs Go testing on the project and formats the test results using the "gotestfmt" tool. The original test log is saved as an artifact for later review.

- **Docker**: The "docker" workflow automates the building and pushing of Docker images. It sets up the necessary environment, authenticates with DockerHub and GitHub Container Registry, extracts metadata for the images, and performs multi-platform builds based on event triggers. The built images are then pushed to the respective registries.

# Terraform

This Terraform configuration sets up the infrastructure required to deploy the "polygon-client" application. It utilizes AWS ECS (Elastic Container Service) to host the application and AWS ECR (Elastic Container Registry) to store the Docker image.

The configuration includes the creation of an ECS cluster named "app" where the application will run. It also creates an AWS CloudWatch Log Group to capture logs from the application for monitoring and troubleshooting purposes.

The main component is the ECS service, which defines how the application should be run. It specifies the task definition, which includes details about the Docker image, logging configuration, and resource requirements such as CPU and memory.

Networking is configured to ensure the application's accessibility. The service is set up to run on AWS Fargate, which eliminates the need to manage the underlying infrastructure. The application is placed in private subnets and associated with security groups to control inbound and outbound traffic.

To expose the application to the internet, an AWS ALB (Application Load Balancer) is created. It distributes incoming traffic to the ECS service using a target group. The ALB is configured with HTTP and HTTPS listeners, redirecting HTTP traffic to HTTPS for secure communication.

Additionally, an ACM certificate is provisioned for enabling HTTPS communication with the ALB.

Overall, this Terraform configuration automates the setup of infrastructure components required to deploy the "polygon-client" application, ensuring scalability, high availability, and secure access.

# Support

For any questions or issues, please contact [rafaribe](mailto:rafael.ntw@gmail.com).

# License

This project is licensed under the MIT License - see the LICENSE file for details.
