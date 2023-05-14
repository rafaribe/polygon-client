# Polygon Client

## Description

The Polygon Client is a Go application that retrieves information from the Polygon network using the Polygon RPC endpoint. It periodically fetches the latest block number and latest block hash and logs the information to the console.

# Prerequisites

To run this application, you need to have the following:

- Docker installed
- Polygon RPC endpoint URL

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


## Configuration

There is no configuration required for this application as the Polygon RPC endpoint is hardcoded. Simply run the executable and it will periodically fetch the latest block information and print it to the console.

## Improvements

In its current state, the Polygon Client is limited in that it only supports a single RPC endpoint that is hardcoded into the application. In order to make this application more flexible and reusable, it would be beneficial to add support for configuration files or command line arguments that would allow the user to specify the RPC endpoint they wish to use. This would enable the application to be used with a wider range of Polygon-compatible RPC endpoints.

## Usage

Once the Docker container is running, the application will periodically log the latest block number and hash to the console. To stop the application, use Ctrl + C.

# Support

For any questions or issues, please contact the [rafaribe](mailto:rafael.ntw@gmail.com) of the application.

# License

This project is licensed under the MIT License - see the LICENSE file for details.
