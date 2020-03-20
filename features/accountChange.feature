Scenario: Add information to the account
    When I put this information
    | Allergies | ["Wheat"] |
    Then the response should be 200 