# Project 'Tigerhall Kittens' Web app based on Go lang - Testing Guide

This document provides a comprehensive guide for testing the various HTTP handlers implemented in our Go project. The handlers facilitate operations related to users, tigers, and sightings, including creation, authentication, and data retrieval. The testing process outlined below has been conducted using Postman.

## Pre-requisites

- Ensure the Go application is running locally and listening on the expected port (default is `:8080`).
- Install Postman on your local machine to execute the HTTP requests.

## Testing Endpoints

### 1. Create User (`/users/create`)

- **Method:** `POST`
- **Body:** JSON payload containing user details (e.g., name, email, password).
- **Purpose:** Registers a new user in the system.


Example json playload data as Input :-

{
  "username": "nwUse44",
  "email": "nwuer@eempwee.com",
  "password": "secueassword2123"
}

Expected: Status Code 201 and a JSON response containing an object with user details, including attributes such as id, username, email, and created_at and excluding password.

............................

### 2. User Login (`/users/login`)

- **Method:** `POST`
- **Body:** JSON payload with user login credentials (e.g., email, password).
- **Purpose:** Authenticates a user and returns a session token or equivalent credentials.

Example json playload data as Input :-

{
  "username": "newUser",
  "password": "securepassword123"
}

Expected: Status Code 200 & and a JSON Response with a message "Logged in successfully"

.....................

### 3. Create Tiger (`/tigers/create`)

- **Method:** `POST`
- **Body:** JSON payload with tiger details (e.g., name, species).
- **Purpose:** Adds a new tiger record to the database.

Example json playload data as Input :-

{
  "name": "Tiger7",
  "date_of_birth": "2013-05-20T12:00:00Z",
  "last_seen_timestamp": "2023-01-11T12:00:00Z",
  "last_seen_lat": 21.0,
  "last_seen_lon": 49.0
}

Expected: Status Code 201 and a JSON response containing an object with detailed information about a tiger, including attributes such as id, name, date_of_birth, last_seen_timestamp, last_seen_lat, and last_seen_lon

.....................


### 4. List All Tigers (`/tigers/list`)

- **Method:** `GET`
- **Purpose:** Retrieves a list of all tigers in the database.

Ex url :=  http://localhost:8080/tigers/list?page=1&pageSize=3

Expected: Status Code 201  & and a JSON array containing objects representing tigers.

.....................

### 5. Create Sighting (`/sightings/create`)

- **Method:** `POST`
- **Body:** Form-data or JSON payload with sighting details, including `tigerID`, location coordinates (`lat`, `lon`), timestamp, and an image file.
- **Purpose:** Records a new sighting of a tiger along with an image.
 
 Url : http://localhost:8080/sightings/create

 with form-data checked it and with key Sightings - and a value like 

{
  "tiger_id": 1,
  "user_id": 2 ,
  "lat": 23.5551,
  "lon": 55.2708,
  "timestamp": "2024-02-11T12:00:00Z"
}

 and then key as image upload a file.
 Should get a JSON response with created sighting information.


Testable Combination :

Scenario: 1

For example, if you pass above json data and it's saved another user with below json info for the same tiger at that point of time code will check was there any record of sightings for the same tiger is there it will pull and calculate the distance as you see coordinates are same then you will end with a custom error message.

{
  "tiger_id": 1,
  "user_id": 3 ,
  "lat": 23.5551,
  "lon": 55.2708,
  "timestamp": "2024-02-11T12:00:00Z"
}

Error Response:

{"code":"TOO_CLOSE_TO_PREVIOUS_SIGHTING","message":"New sighting is too close to the last sighting. Sightings must be at least 5 kilometers apart."}

Scenario: 2

Lets say another user have sent the below json data for the same tiger but with different coordinates at that point of time, code will check and see are there any previous sightings for the tiger by any users, if there then we have notification system provided using a message queue with channels , this will send out notification email message for the previous sighted users and sighting will be saved with a JSON response back to the user about the current saved sighting.

{
  "tiger_id": 1,
  "user_id": 7 ,
  "lat": 23.5551,
  "lon": 55.2708,
  "timestamp": "2024-02-11T12:00:00Z"
}


Log Messages :

Sending email to user 3 about tiger 1
2024/02/11 13:27:20 Sending email to User ID: 3, for Tiger ID: 1;

Sending email to user 2 about tiger 1
2024/02/11 13:27:20 Sending email to User ID: 2, for Tiger ID: 1;

..........................

### 6. List Sightings (`/sightings/list`)

- **Method:** `GET`
- **Parameters:** `tigerID` (required), `page`, `pageSize` for pagination.
- **Purpose:** Fetches a paginated list of sightings for a specified tiger.


Ex : http://localhost:8080/sightings/list?tigerID=4&page=1&pageSize=10

Expected: Status Code 200  & and a JSON array containing objects representing Sightings.


................

## Testing Instructions

1. **Start the Application:** Ensure your Go server is running.
2. **Open Postman:** Launch Postman and create a new request for each handler according to the details above.
3. **Configure the Request:** Select the appropriate HTTP method, enter the request URL (e.g., `http://localhost:8080/users/create` for user creation), and, if required, set up the request body.
4. **Send the Request:** Click the "Send" button in Postman to execute the request.
5. **Review the Response:** Examine the status code and response body in Postman to verify the expected outcome.

## Troubleshooting

- Ensure the server is running and accessible at the specified port.
- Validate the request payload and parameters for correctness.
- Review server logs for any errors or warnings that might indicate issues with request processing.

## Unit Tests

This project includes comprehensive unit tests covering both models and handlers.


### Test Structure

- **Models Tests:** Located in the `models` directory, these tests cover our data structures and database interactions, ensuring that our models behave as expected under various conditions.

- **Handlers Tests:** Found in the `handlers` directory, these tests validate the logic within our HTTP handlers, confirming that they correctly process requests and respond as intended.

NOTE: Was not able to complete Unitests for Sightings hanlders ontime.

To run the unit tests for models and handlers, follow these steps:

- ** Go to Models director and you can run " go test " 
- ** other way individually you can run going to each test file.

