# Golang Web App with Active Directory Authentication

A little proof of concept showing how to authenticate a user against Active Directory, serve HTML templates and handling sessions with a custom in-memory storage using `mux` from the [Gorilla Toolkit](https://www.gorillatoolkit.org/pkg/mux).

App flow:

- Login
  - [x] Check username and password against AD
  - [x] If it's valid create a session
  - [x] Store session in db
  - [x] Create an insecure cookie with session id (in the future use HMAC)
  - [x] Check user session id from cookie against session stored in db to verify if it's valid
- Restricted resource behind login
  - [x] If session is invalid redirect to login otherwise show the resource
- Logout
  - [ ] Deletes session in db
  - [x] Expires cookie
