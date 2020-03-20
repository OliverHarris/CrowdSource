    Scenario: Get the users account information
        Given The user has set some information
        Then the response should include
            | Allergies              | ["Wheat"]                                                                                           |
            | "RecommendedNutrition" | { "Carbohydrate": 260,"Energy": 2000,"Fat": 70,"Protein": 50,"Salt": 6,"Saturates": 20,"Sugar": 90} |