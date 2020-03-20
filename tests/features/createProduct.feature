Feature: Create product

    Create a product and give it some information

    Scenario: Product with full information
        Given the product "test001" doesn't exist
        When I put this information
            """
            {
                "ID": "test001",
                "productName": {
                    "name": "test product"
                },
                "ingredients": {
                    "ingredients": [
                        "Wheat",
                        "Egg",
                        "Sugar"
                    ]
                },
                "nutrition": {
                    "nutrition": {
                        "Energy": [
                            500
                        ],
                        "Fibre": [
                            0
                        ]
                    }
                }
            }
            """
        And I send "POST" request to "/product/test001"
        Then the response should be 200

    Scenario: Product with some information
        Given the product "test002" doesn't exist
        When I put this information
            """
            {
                "ID": "test002",
                "productName": {
                    "name": "test product"
                },
                "ingredients": {
                    "ingredients": [
                        "Wheat",
                        "Egg",
                        "Sugar"
                    ]
                }
            }
            """
        And I send "POST" request to "/product/test002"
        Then the response should be 200