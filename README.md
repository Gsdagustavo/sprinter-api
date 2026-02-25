# Sprinter API

This repository contains the backend service for the **Sprinter** project.
The frontend application is available at:
[https://github.com/VitorFranciscoDev/sprinter-app](https://github.com/VitorFranciscoDev/sprinter-app)

---

## Overview

The Sprinter API is responsible for:

* Implementing the business logic
* Managing database communication
* Handling authentication and authorization
* Processing and validating data

It is designed to operate in conjunction with the Sprinter frontend application.

---

## Development Configuration

The local development configuration is defined in the file:

```
dev-settings.toml
```

It is strongly recommended not to modify this file. Changes may disrupt API services or lead to unintended behavior.

---

## Running the Project Locally

To run the API in a local development environment, follow the steps below.

### 1. Create the Database Container

Create and start the database container using the configuration specified in `dev-settings.toml`.

Ensure that:

* The database service is running
* The connection parameters match those defined in the configuration file

### 2. Run the Application

Start the application using the following program arguments:

```
-action run
-settings dev-settings.toml
```

These arguments must be provided through your development environment or command-line interface.

---

## Recommended Execution Order

1. Start the database container
2. Verify that the database is accessible
3. Run the API with the required arguments
4. Start the frontend application

---

## Additional Notes

* Always use `dev-settings.toml` for local development.
* Avoid committing sensitive configuration changes.
* Ensure that the database container is fully initialized before starting the API.
