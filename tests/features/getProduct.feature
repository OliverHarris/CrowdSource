Feature: Get product information

    Requesting a product using the barcode

    Scenario: Request a product
        Given the product "test002" does exist
        When  I send "GET" request to "/product/test001"
        Then the response should include
            """
            "name": "test product"
            """
        And the response should include
            """
            "ingredients": [
            "Wheat",
            "Egg",
            "Sugar"
            ]
            """
        And the response should include
            """
            "nutrition": {
            "Energy": 500,
            "Fibre": 0
            }
            """