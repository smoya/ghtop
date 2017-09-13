Feature: Top endpoint
    As a user I want to retrieve all the top github contributors

    Scenario: Default values
        Given a server
        And having these contributors "gthreepwood,bbernoulli"
        When I "GET" request to "/top?location='barcelona'"
        Then response is successful
        And response is JSON
        And response data should match JSON list:
        """
        [
           {
              "username":"gthreepwood",
              "vcs":"github",
              "profile_url":"https://github.com/gthreepwood",
              "avatar_url":"http://foo.bar"
           },
           {
              "username":"bbernoulli",
              "vcs":"github",
              "profile_url":"https://github.com/bbernoulli",
              "avatar_url":"http://foo.bar"
           }
        ]
        """

    Scenario: Limiting results
        Given a server
        And having these contributors "gthreepwood,bbernoulli"
        When I "GET" request to "/top?location='barcelona'&limit=1"
        Then response is successful
        And response is JSON
        And response data should match JSON list:
        """
        [
           {
              "username":"gthreepwood",
              "vcs":"github",
              "profile_url":"https://github.com/gthreepwood",
              "avatar_url":"http://foo.bar"
           }
        ]
        """

    Scenario: Erroring because missing location
        Given a server
        And having these contributors "gthreepwood,bbernoulli"
        When I "GET" request to "/top"
        Then response is error 400

