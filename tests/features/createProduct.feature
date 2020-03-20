Feature: Create product

    Create a product and give it some information

    Scenario: Product with full information
        Given the product test001 doesn't exist
        When I request to create a product
        And  I put this information
            | attribute   | value                    |
            | productName | test product             |
            | ingredients | LIST Wheat,Egg,Sugar     |
            | nutrition   | {"Energy":500,"Fibre":0} |
        Then the response should include
            | attribute   | value                     |
            | productName | test product              |
            | ingredients | ["Wheat", "Egg", "Sugar"] |
            | nutrition   | {"Energy":500,"Fibre":0}  |

    Scenario: Product with some information
        Given the product test002 doesn't exist
        When I request to create a product
        And  I put this information
            | attribute   | value                |
            | productName | test product         |
            | ingredients | LIST Wheat,Egg,Sugar |
        Then the response should include
            | attribute   | value                     |
            | productName | test product              |
            | ingredients | ["Wheat", "Egg", "Sugar"] |
            | nutrition   |                           |