Feature: Get currency
  I need to be able to get conversions

  Scenario: Get the AUD conversion
    Given the conversion currency is "AUD"
    Then the conversion given should be "1.5"