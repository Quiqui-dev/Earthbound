package middleware

import "net/http"

// a func to wrap our handlers with custom functions
type Middleware func(http.Handler) http.Handler

// a helper to build out an order in which to apply multiple middleware functions
// to a final handler
type Chain struct {
	middlwares []Middleware
}

// Append a new middleware to the existing list
func (c *Chain) Use(m Middleware) {
	c.middlwares = append(c.middlwares, m)
}

func (c *Chain) Then(h http.Handler) http.Handler {
	for i := len(c.middlwares) - 1; i >= 0; i-- {
		h = c.middlwares[i](h)
	}

	return h
}
