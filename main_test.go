package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"log"
	"os"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
)

func mock_http_get(url string) []byte {
	return []byte(`{
    "Realtime Currency Exchange Rate": {
        "1. From_Currency Code": "USD",
        "2. From_Currency Name": "United States Dollar",
        "3. To_Currency Code": "AUD",
        "4. To_Currency Name": "Australian Dollar",
        "5. Exchange Rate": "1.5",
        "6. Last Refreshed": "2022-11-28 23:37:49",
        "7. Time Zone": "UTC",
        "8. Bid Price": "1.50328000",
        "9. Ask Price": "1.50328000"
    }
}`)
}

func init() {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func TestGet_exchange(t *testing.T) {
    exchange_rate := Get_exchange(mock_http_get, os.Getenv("API_KEY"), "USD", "AUD")

    if exchange_rate != 1.5 {
    	t.Errorf("USD/AUD exchange = %f; want 1.5", exchange_rate)
    }
}

// godogsCtxKey is the key used to store the available godogs in the context.Context.
type godogsCtxKey struct{}

func theConversionCurrencyIs(ctx context.Context, currency string) (context.Context, error) {
  return context.WithValue(ctx, godogsCtxKey{}, currency), nil
}

func theConversionGivenShouldBe(ctx context.Context, expected string) error {
  currency, ok := ctx.Value(godogsCtxKey{}).(string)
  if !ok {
    return errors.New("There is no currency")
  }

  exchange_rate := Get_exchange(mock_http_get, os.Getenv("API_KEY"), "USD", currency)

  expected_f, _ := strconv.ParseFloat(expected, 64)

  if expected_f != exchange_rate {
    return fmt.Errorf("expected %s currency conversion to be %.2f, but it is %.2f", currency, expected_f, exchange_rate)
  }

  return nil
}

func TestFeatures(t *testing.T) {
  suite := godog.TestSuite{
    ScenarioInitializer: InitializeScenario,
    Options: &godog.Options{
      Format:   "pretty",
      Paths:    []string{"features"},
      TestingT: t, // Testing instance that will run subtests.
    },
  }

  if suite.Run() != 0 {
    t.Fatal("non-zero status returned, failed to run feature tests")
  }
}

func InitializeScenario(sc *godog.ScenarioContext) {
  sc.Step(`^the conversion currency is "([^"]*)"$`, theConversionCurrencyIs)
  sc.Step(`^the conversion given should be "([^"]*)"$`, theConversionGivenShouldBe)
}
