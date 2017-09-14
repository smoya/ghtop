# ghtop
[![Build Status](https://travis-ci.org/smoya/ghtop.svg?branch=master)](https://travis-ci.org/smoya/ghtop)

A simple Golang http service that lists the top Github contributors given a location.

Features included:

* Top contributors by a given location.
* Cached results with a given TTL.
* Number of results customizable (Max 150).
* Sort by the amount of **repositories**, **followers** or by the **date they joined**. 

## Installation

### From docker:

There is a Docker image called `smoya/ghtop`. In order to run the service, just do:

```bash
docker run --publish 8080:8080 smoya/ghtop:latest -gh-token=<github token>
```

### From source

Ghtop requires Go 1.9 or later.
```bash
go get -u github.com/smoya/ghtop
```

## Usage

Run the server:

```bash
ghtop -gh-token=<github token> -ttl=<ttl in seconds>
```

The `GET /top` endpoint should be now mounted.

Curl [http://localhost:8080/top?location=barcelona](http://localhost:8080/top?location=barcelona) in order to see the top dev contributors in Barcelona area.
Replace `barcelona` with whatever location you want to look at. 

```bash
curl 'http://localhost:8080/top?location=barcelona'
```

## Authentication

If authentication is required, the application is ready for use Basic HTTP authentication [rfc2617](https://tools.ietf.org/html/rfc2617).
In further versions [JWT](https://jwt.io) could be implemented.

In order to prompt the authentication dialog on the `/top` endpoint, you must specify the user and password when running the service.

### Arguments reference:

| name           | type   | description                                                                             | required | default |
|----------------|--------|-----------------------------------------------------------------------------------------|----------|---------|
| -port          | int    | Server's listening port                                                                 | no       | 8080    |
| -env           | string | Sets the environment. Just for logging.                                                 | no       | prod    |
| -gh-token      | string | Github personal access token.  Create yours from https://github.com/settings/tokens/new | yes      |         |
| -ttl           | int    | The ttl in seconds for the repository cache.                                            | no       | 300     |
| -auth-user     | string | The username for basic Http Authentication. No Auth if empty                            | no       |         |
| -auth-password | string | The password for basic Http Authentication.                                             | no       |         |

## Tests

Ghtop library tests are split in two:

* Unit tests
* End To End tests (features are located in the [e2e pkg](pkg/e2e/features))

Run `make tests` and `make e2e` respectively, or `make check` in order to run all of them.

## Design considerations

* The code allows us to easy implement a new repository for any other vcs system rather than github.
* The limit has no fixed values but the max is 150 results (avoiding possible performance issues). I considered that adding the possibility to the user to 
limit the amount of results in each request was better for UX.
* The Github token is set at application level rather than by the user in each request. 
This is a design decision since I considered this service as a simple API endpoint for a possible nice frontend application.
However this could be changed handling a given Github token in each `GET /top` request.  

## TODO

* Repository cache can be improved by checking if the specified limit is lower or equal than the cached one and use that cache.
* We could retrieve more data from Github like user full name, email, organizations, etc. 
* Add a multiple repository constructor in order to allow using several repositories (from different vcs) at the same time.
* Add a proper authentication mechanism ([JWT](https://jwt.io) for example). This could be done through an in memory database + a simple http middleware.
* Adding more tests.
