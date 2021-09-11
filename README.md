[![Go Reference](https://pkg.go.dev/badge/astuart.co/go-robinhood/v2.svg)](https://pkg.go.dev/astuart.co/go-robinhood/v2)

# Robinhood the rich and feeding the poor, now automated

> Even though robinhood makes me poor

## Notice

The v2 API has changed to allow a context to be passed into various calls.

If you have used this library before, and use credential caching, you will need
to remove any credential cache and rebuild if you experience errors.

## General usage

```go
ctx = context.Background()
c, err := robinhood.Dial(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
// handle err

i, err := c.GetInstrumentForSymbol(ctx, "SPY")
// handle err

o, err := c.Order(ctx, i, robinhood.OrderOpts{
  Price: 100.0,
  Side: robinhood.Buy,
  Quantity: 1,
})
// handle err

time.Sleep(5*time.Second) //Let me think about it some more...

// Ah crap, I need to buy groceries.
err := o.Cancel(ctx)

if err != nil {
  // Oh well
}
```

See [cmd/example](cmd/example/example.go) for a runnable example.