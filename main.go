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

func main() {
    flag.Parse()

    url := "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s&apikey=%s"

    resp, err := http.Get(fmt.Sprintf(url, "USD", *to_currency, os.Getenv("API_KEY")))
    if err != nil {
       log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       log.Fatalln(err)
    }

    var container CurrencyContainer

    json.Unmarshal(body, &container)

    Update(container.Conversion.ToCurrency, container.Conversion.ExchangeRate)
}
