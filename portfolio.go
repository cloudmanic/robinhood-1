package robinhood

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Portfolio information

type position struct {
	BuyPrice string `json:"average_buy_price"`
	URL      string `json:"instrument"`
}

func (c *Client) rawPortfolio() ([]position, error) {
	parms := url.Values{}
	parms.Set("nonzero", "true")
	resp, err := c.paginatedGet(accountsURI + c.AccountID + "/" + positionsURI + "?" + parms.Encode())
	if err != nil {
		return nil, err
	}
	fmt.Printf("PaginatedGet got: %s", resp)
	var positions []position
	err = json.Unmarshal(resp, &positions)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

type detailedPosition struct {
  Symbol string `json:"symbol"`
	Name   string `json:"simple_name"`
}

type Position struct {
  Symbol string
  Name string
  BuyPrice float64
}

func (c *Client) Portfolio() ([]Position, error) {
  var positions []Position
	pos, err := c.rawPortfolio()
	if err != nil {
		return nil, err
	}
	for _, p := range pos {
		fmt.Printf("Going to get %s", p.URL)
		req, err := http.NewRequest("GET", p.URL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := doReq(req, c.Token)
		if err != nil {
			// TODO: return partial info from here.
			return nil, err
		}
		fmt.Printf("\nGot detailed position: %s\n", resp)
    var detail detailedPosition
    err = json.Unmarshal(resp, &detail)
    if err != nil {
			// TODO: return partial info from here.
			return nil, err
		}
    buyPrice, err := strconv.ParseFloat(p.BuyPrice, 64)
    if err != nil {
			// TODO: return partial info from here.
      return nil, err
    }
    positions = append(positions, Position{
      Symbol: detail.Symbol,
      Name: detail.Name,
      BuyPrice: buyPrice,
    })
	}
	return positions, nil
}
