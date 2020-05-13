package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"os"
)

func getMarketData(params ...interface{}) {
	var (
		err      error
		rapidkey string
		result   struct {
			MarketSummaryResponse struct {
				Result []struct {
					FullExchangeName     string `json:"fullExchangeName"`
					ExchangeTimezoneName string `json:"exchangeTimezoneName"`
					Symbol               string `json:"symbol"`
					RegularMarketChange  struct {
						Raw float64 `json:"raw"`
						Fmt string  `json:"fmt"`
					} `json:"regularMarketChange"`
					GmtOffSetMilliseconds int    `json:"gmtOffSetMilliseconds"`
					ExchangeDataDelayedBy int    `json:"exchangeDataDelayedBy"`
					Language              string `json:"language"`
					RegularMarketTime     struct {
						Raw int    `json:"raw"`
						Fmt string `json:"fmt"`
					} `json:"regularMarketTime"`
					ExchangeTimezoneShortName  string `json:"exchangeTimezoneShortName"`
					RegularMarketChangePercent struct {
						Raw float64 `json:"raw"`
						Fmt string  `json:"fmt"`
					} `json:"regularMarketChangePercent"`
					QuoteType          string `json:"quoteType"`
					MarketState        string `json:"marketState"`
					RegularMarketPrice struct {
						Raw float64 `json:"raw"`
						Fmt string  `json:"fmt"`
					} `json:"regularMarketPrice"`
					Market                     string `json:"market"`
					PriceHint                  int    `json:"priceHint,omitempty"`
					Tradeable                  bool   `json:"tradeable"`
					SourceInterval             int    `json:"sourceInterval"`
					Exchange                   string `json:"exchange"`
					Region                     string `json:"region"`
					ShortName                  string `json:"shortName,omitempty"`
					Triggerable                bool   `json:"triggerable"`
					RegularMarketPreviousClose struct {
						Raw float64 `json:"raw"`
						Fmt string  `json:"fmt"`
					} `json:"regularMarketPreviousClose"`
					HeadSymbolAsString string `json:"headSymbolAsString,omitempty"`
					Currency           string `json:"currency,omitempty"`
					LongName           string `json:"longName,omitempty"`
					QuoteSourceName    string `json:"quoteSourceName,omitempty"`
				} `json:"result"`
				Error interface{} `json:"error"`
			} `json:"marketSummaryResponse"`
		}
	)
	logger.Debug().Msg("Getting Market Data")
	mktdata := params[0].(*marketData)
	rapidkey = os.Getenv("RAPIDAPI_KEY")
	if len(rapidkey) == 0 {
		logger.Error().Msg("Missing environment variable RAPIDAPI_KEY")
		return
	}
	url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/market/get-summary?region=US&lang=en"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "apidojo-yahoo-finance-v1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", rapidkey)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		logger.Error().Err(err).Msg("error decoding JSON from rapidApi")
	} else if result.MarketSummaryResponse.Error != nil {
		logger.Error().Msg(spew.Sdump(result.MarketSummaryResponse.Error))
	} else {
		for _, r := range result.MarketSummaryResponse.Result {
			if r.Symbol == "^IXIC" {
				mktdata.lines[0] = "NASDAQ"
				mktdata.lines[1] = fmt.Sprintf("%s %s", r.RegularMarketPrice.Fmt, r.RegularMarketChangePercent.Fmt)
			}
		}
	}
}
