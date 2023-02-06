# Calculator

A simple system to compute the mean, median, and mode of a sequence of numbers.

## Getting Started

To use this system, clone the repository and run the code on your local machine. The system is written in Go, so having Go installed is a requirement.

## Prerequisites

Before you start, make sure you have Go installed on your machine. You can download it from the official website at https://golang.org/. Installing Docker is optional, but if you prefer to run the system another way, it must be installed as well.

* GO 1.19.1
* Docker 20.10.12

## Running the System

### Keyboard Input
To run the system, follow these steps:

1. Clone the repository:

```shell=
$ git clone https://github.com/dalaoqi/calculator.git
```

2. Go to the repository directory:
```shell=
$ cd calculator
```
3. In the `bin` directory, you will find four binary files. To execute them, follow this order: run `bin/server` in one command window, then run `bin/client1`, `bin/client2`, and `bin/client3` in three separate command windows.

* You can regenerate binary files in the `/bin` directory by running `scripts/run_binary.sh`.

4. Once the clients are ready, you will see the message, 
```
You can now input integers and press [ENTER]. The clients will display the mean, median, and mode of the input values.
```
You can then input a sequence of numbers using your keyboard, and the results of computing the mean, median, and mode will be shown.

#### Example

![](https://i.imgur.com/4nZEYhh.png)


### Run with Docker, HTTP Request

This system is capable of supporting multiple concurrent processes running inside Docker. To use it, you will first need to build the Docker image. The relevant scripts can be found in the scripts directory.

1. build the image
Use the following command to build the Docker image:
```shell=
$ ./scripts/build_docker.sh
```
If the build is successful, you can verify the image by running docker images and looking for an image named `calculator`.

2. run the image
To run the Docker image, use the following command:
```shell=
$ ./scripts/run_docker.sh
```

To verify that the image is running, use `docker ps` and look for an image named `calculator`.

3. Checking the logs
Use the following command to view the logs of the container:
```shell=
$ docker logs -f calculator
```
This will display the logs and show that all the binaries are running under `supervisor`. 

4. Sending a Request

To send a request, open a command window and enter the following:
```shell=
$ curl --location --request GET 'localhost:8081/' \
--header 'Content-Type: text/plain' \
--data-raw '1 5 5 10 15 2 3'
```

This will return the results, which will be displayed in the logs of the container and in the command window.

#### Example
The logs of calculator are displayed on the left, and the results of the request sent are displayed on the right.
![](https://i.imgur.com/VVrBgsb.png)
