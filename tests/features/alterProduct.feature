Feature: Changing an existing product

    Change some product information for an existing product

    Scenario: Update product name
        Given the product test002 does exist
        When I request to create a product
        And  I put this information
            | attribute   | value                |
            | productName | test product changed |
        Then the response should include
            | attribute   | value                     |
            | productName | test product changed      |
            | ingredients | ["Wheat", "Egg", "Sugar"] |
            | nutrition   |                           |

    Scenario: Update product name and ingredients
        Given the product test002 does exist
        When I request to create a product
        And  I put this information
            | attribute   | value                      |
            | productName | test product               |
            | ingredients | LIST Wheat,Egg,Sugar,Water |
        Then the response should include
            | attribute   | value                             |
            | productName | test product                      |
            | ingredients | ["Wheat", "Egg", "Sugar","Water"] |
            | nutrition   |                                   |

    Scenario: Update product nutrition
        Given the product test002 does exist
        When I request to create a product
        And  I put this information
            | attribute | value                    |
            | nutrition | {"Energy":500,"Fibre":0} |
        Then the response should include
            | attribute   | value                             |
            | productName | test product                      |
            | ingredients | ["Wheat", "Egg", "Sugar","Water"] |
            | nutrition   | {"Energy":500,"Fibre":0}          |