# Sear Reservation API

## Introduction

This repository provides sample endpoints for creating event and performing seat reservations.
Each user is allowed to perform at most 1 seat.

## Features

### API Functionality

Repository provides following functionalities:

- Create Event
- Query Event
- Create Reservation
- Cancel Reservation
- Confirm if reservation is done

### Avoiding Double Booking

It prevents users from performing double booking for the same seat.
Specific scenarios where double booking could occur and methods designed to prevent it possible to find within tests.

### Environment

System designed considering scalability in mind allowing for adding more features.

## Tech Stack

- Language: Go
- Deployment: Docker
- Client: No UI, link for postman can be found in [here](https://www.postman.com/titicorp/workspace/plugo-seat-reservation-task/collection/27702330-af68028f-76d0-42fd-ba0e-c6d422508c0e?action=share&creator=27702330)

## Getting Started

1. Clone the repository

```bash
#!/bin/bash
git clone https://github.com/Kalomidin/plugo-seat-reservation-go.git
```

2. Setup the docker

```bash
#!/bin/bash
./docker-start.sh
```

3. Create repo where u want reside the build:

```bash
#!/bin/bash
mkdir bin
```

4. Build the repository

```bash
#!/bin/bash
go build -o bin ./...
```

5. Run the builded image

```bash
#!/bin/bash
./bin/seat-reservation-api 
```


6. Open the postman and start using the endpoints.
