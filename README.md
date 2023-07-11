
## Load Balancer

The load balancer is responsible for distributing incoming requests among the available services. It uses a round-robin algorithm to evenly distribute the load. The load balancer periodically checks the health of the services and adjusts the routing accordingly.

### Configuration

The load balancer configuration is stored in `config.json`, which contains the URLs of the services to load balance.

### Building and Running

To build and run the load balancer, follow these steps:

1. Navigate to the `load-balancer` directory.
2. Build the Docker image using the provided Dockerfile: `docker build -t load-balancer .`
3. Start the load balancer container: `docker run -p 8080:8080 load-balancer`

The load balancer will listen on port 8080 and distribute incoming requests to the services.

## Services

The project includes five services, namely `service1`, `service2`, and other 3 services. Each service handles specific functionality and communicates with the load balancer.

### Building and Running

To build and run a service, follow these steps:

1. Navigate to the respective service directory (`service1`, `service2`, or `service3`, ...).
2. Build the Docker image using the provided Dockerfile: `docker build -t service1 .`
3. Start the service container: `docker run -p <port>:8080 service1`

Replace `<port>` with the desired port number for the service.

## Usage

Once the load balancer and services are running, you can send requests to the load balancer's address (e.g., `http://localhost:8080`). The load balancer will distribute the requests to the available services based on the load balancing algorithm.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

