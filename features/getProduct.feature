Feature: Get product information

    Requesting a product using the barcode

    Scenario: Request a product
        Given the product test001 exists
        When I request test001
        Then my response should include
            | attribute   | value                     |
            | productName | test product              |
            | ingredients | ["Wheat", "Egg", "Sugar"] |
            | nutrition   | {"Energy":500,"Fibre":0}  |