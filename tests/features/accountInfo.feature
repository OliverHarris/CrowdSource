Feature: Get the account information

    Get allergies and RecommendedNutrition from the user account

    Scenario: Get the users account information
        Given I send "GET" request to "/account"
        Then the response should be 200