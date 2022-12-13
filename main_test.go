package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/joho/godotenv"

	"net/url"
)

func mock_http_get(get_url string) []byte {
	url_parts, _ := url.Parse(get_url)
	query, _ := url.ParseQuery(url_parts.RawQuery)

	currency_lookup := map[string]float64{
		"AUD": 1.5,
		"GBP": 2.5,
	}

	currency_value, _ := currency_lookup[query["to_currency"][0]]

	return []byte(fmt.Sprintf(`{
    "Realtime Currency Exchange Rate": {
        "1. From_Currency Code": "USD",
        "2. From_Currency Name": "United States Dollar",
        "3. To_Currency Code": "AUD",
        "4. To_Currency Name": "Australian Dollar",
        "5. Exchange Rate": "%f",
        "6. Last Refreshed": "2022-11-28 23:37:49",
        "7. Time Zone": "UTC",
        "8. Bid Price": "1.50328000",
        "9. Ask Price": "1.50328000"
    }
}`, currency_value))
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

type currencyKey struct{}
type exchangeKey struct{}
type apiKey struct{}

func thatTheAPIIs(ctx context.Context, api_value string) (context.Context, error) {
	return context.WithValue(ctx, apiKey{}, api_value), nil
}

func theConversionCurrencyIs(ctx context.Context, currency string) (context.Context, error) {
	return context.WithValue(ctx, currencyKey{}, currency), nil
}

func theConversionApiIsCalled(ctx context.Context) (context.Context, error) {
	currency, ok := ctx.Value(currencyKey{}).(string)
	if !ok {
		return ctx, errors.New("There is no currency")
	}

	exchange_rate := Get_exchange(mock_http_get, os.Getenv("API_KEY"), "USD", currency)

	return context.WithValue(ctx, exchangeKey{}, exchange_rate), nil
}

func theConversionGivenShouldBe(ctx context.Context, expected string) error {
	exchange_rate, ok := ctx.Value(exchangeKey{}).(float64)
	if !ok {
		return errors.New("There is no exchange rate result")
	}

	currency, ok := ctx.Value(currencyKey{}).(string)
	if !ok {
		return errors.New("There is no currency")
	}

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
		t.Fatal("Non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^that the API is ([a-z]+)$`, thatTheAPIIs)

	ctx.Step(`^the conversion currency is "([^"]*)"$`, theConversionCurrencyIs)
	ctx.Step(`^the conversion API is called$`, theConversionApiIsCalled)
	ctx.Step(`^the conversion given should be "([^"]*)"$`, theConversionGivenShouldBe)
}
