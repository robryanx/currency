package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "io/ioutil"
    "net/http"
    "encoding/json"
    
    "github.com/joho/godotenv"
)

var to_currency *string

type CurrencyContainer struct {
    Conversion CurrencyConversion `json:"Realtime Currency Exchange Rate"`
}

type CurrencyConversion struct {
    FromCurrency string `json:"1. From_Currency Code"`
    ToCurrency string `json:"3. To_Currency Code"`
    ExchangeRate float64 `json:"5. Exchange Rate,string"`
}

func init() {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }

    to_currency = flag.String("to_currency", "AUD", "currency")
}

func http_get(url string) []byte {
    resp, err := http.Get(url)
    if err != nil {
       log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       log.Fatalln(err)
    }

    return body
}

func Get_exchange(http_get_func func(url string) []byte, api_key string, from_currency string, to_currency string) float64 {
    url := "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s&apikey=%s"

    body := http_get_func(fmt.Sprintf(url, "USD", to_currency, api_key))

    var container CurrencyContainer

    json.Unmarshal(body, &container)

    return container.Conversion.ExchangeRate
}

func main() {
    flag.Parse()

    exchange_rate := Get_exchange(http_get, os.Getenv("API_KEY"), "USD", *to_currency)

    Update(*to_currency, exchange_rate)
}
