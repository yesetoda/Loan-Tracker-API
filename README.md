
# Loan Tracker API

## Overview

The Loan Tracker API is a RESTful API built with Golang using the Gin framework. It allows users to manage loans, apply for new loans, and view loan statuses. Admin functionalities are also included for monitoring system activity and managing loan records.

## Features

- **User Management**
  - Register a new user
  - Email verification
  - Password reset

- **Loan Management**
  - Apply for a loan
  - View loan status
  - Admin functionalities for loan management

- **Admin Functionalities**
  - View system logs
  - Manage loan records

- **Security**
  - JWT-based authentication
  - Password hashing with bcrypt

## Technologies

- **Golang**: The programming language used.
- **Gin**: Web framework for building the API.
- **MongoDB**: NoSQL database for storing data.
- **Gingonic**: Middleware for enhanced functionality.
- **JWT-go**: For JSON Web Token handling.
- **bcrypt**: For secure password hashing.

## Endpoints

### User Endpoints

- **POST /users/register**
  - Register a new user.
  - Request Body: `{ "email": "string", "password": "string" }`

- **POST /users/verify-email**
  - Verify a user's email address.
  - Request Body: `{ "token": "string" }`

- **POST /users/reset-password**
  - Request a password reset.
  - Request Body: `{ "email": "string" }`

### Loan Endpoints

- **POST /loans/apply**
  - Apply for a new loan.
  - Request Body: `{ "amount": "number", "term": "number" }`

- **GET /loans/status**
  - View the status of an existing loan.
  - Request Params: `{ "loan_id": "string" }`

### Admin Endpoints

- **GET /admin/logs**
  - View system logs.
  
- **PUT /admin/loans/{loan_id}**
  - Manage (update) loan records.
  - Request Body: `{ "status": "string" }`

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/loan-tracker-api.git
   ```

2. Navigate to the project directory:
   ```bash
   cd loan-tracker-api
   ```

3. Install the necessary dependencies:
   ```bash
   go mod tidy
   ```

4. Configure environment variables as needed. Create a `config.yaml` file based on `config_sample.yaml`.

5. Run the application:
   ```bash
   go run main.go
   ```


## Documentation

API documentation can be accessed using Postman. Import the provided Postman collection file located in the [postman](https://documenter.getpostman.com/view/37276877/2sAXjGduTp) 



## Contact

For any questions or issues, please contact [yeneineh seiba](mailto:yeneinehseiba@gmail.com).
