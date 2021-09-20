[![Go Reference](https://pkg.go.dev/badge/github.com/tstromberg/roho.svg)](https://pkg.go.dev/github.com/tstromberg/roho)

# RoHo

An idiomatic Go client for Robinhood. Based on https://github.com/andrewstuart/go-robinhood

## General usage

```go
c, err := roho.New(nil)
i, err := c.Lookup("SPY")

o, err := cli.Buy(i, roho.OrderOpts{
  Price: 100.0,
  Quantity: 1,
})

// Oh no! Don't buy that!
o.Cancel()
```
