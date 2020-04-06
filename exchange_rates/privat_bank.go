package exchange_rates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
)

var monthConst = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

type ExchangeRate struct {
	BaseCurrency   string  `json:"baseCurrency"`
	Currency       string  `json:"currency"`
	SaleRateNB     float64 `json:"saleRateNB"`
	PurchaseRateNB float64 `json:"purchaseRateNB"`
	SaleRate       float64 `json:"saleRate,omitempty"`
	PurchaseRate   float64 `json:"purchaseRate,omitempty"`
}

type Rates struct {
	Date            string         `json:"date"`
	Bank            string         `json:"bank"`
	BaseCurrency    int            `json:"baseCurrency"`
	BaseCurrencyLit string         `json:"baseCurrencyLit"`
	ExchangeRate    []ExchangeRate `json:"exchangeRate"`
}

func GetRates() (rates string, err error) {
	log.Info("GetRates")
	client := http.Client{}
	var req *http.Request
	year, month, date := time.Now().Date()

	if req, err = http.NewRequest("GET",
		"https://api.privatbank.ua/p24api/exchange_rates?json", nil); err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("date", strconv.Itoa(date)+"."+strconv.Itoa(monthConst[month.String()])+"."+strconv.Itoa(year))
	req.URL.RawQuery = q.Encode()

	var httpRsp *http.Response

	if httpRsp, err = client.Do(req); err != nil {
		return "", err
	}
	var bytesDataFromHttp []byte

	if bytesDataFromHttp, err = ioutil.ReadAll(httpRsp.Body); err != nil {
		return "", err
	}
	var ratesData Rates

	if err = json.Unmarshal(bytesDataFromHttp, &ratesData); err != nil {
		return "", err
	}
	rates = ratesData.String()
	return rates, nil
}

func (r *Rates) String() string {
	log.Info("GetRates String()")

	res := "Date: " + r.Date +
		"\nBank: " + r.Bank

	for _, rate := range r.ExchangeRate {

		switch rate.Currency {
		case "USD":
			res += "\nUSD" +
				"\n								Sale Rate:" + fmt.Sprintf("%.2f", rate.SaleRate) + " UAN" +
				"\n								Purchase Rate:" + fmt.Sprintf("%.2f", rate.PurchaseRate) + " UAN"

		case "EUR":
			res += "\nEUR" +
				"\n								Sale Rate:" + fmt.Sprintf("%.2f", rate.SaleRate) + " UAN" +
				"\n								Purchase Rate:" + fmt.Sprintf("%.2f", rate.PurchaseRate) + " UAN"

		case "RUB":
			res += "\nRUB" +
				"\n								Sale Rate:" + fmt.Sprintf("%.2f", rate.SaleRate) + " UAN" +
				"\n								Purchase Rate:" + fmt.Sprintf("%.2f", rate.PurchaseRate) + " UAN"
		}
	}
	return res
}
