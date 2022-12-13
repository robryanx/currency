Feature: Get currency
  I need to be able to get conversions

  Background:
    Given that the API is x

  Scenario: Get the AUD conversion
    Given the conversion currency is "AUD"
    When the conversion API is called
    Then the conversion given should be "1.5"

  Scenario: Get the GBP conversion
    Given the conversion currency is "GBP"
    When the conversion API is called
    Then the conversion given should be "2.5"