Feature: Update account information

    Add allergies to account

    Scenario: Add information to the account
        When I put this information
            | Allergies | ["Wheat"] |
        And send "POST" request to "/account"
        Then the response should be 200