package middleware

import "frame/context"

type HandleFun func(c *context.Conetxt)

type FilterBuilder func(next Filter) Filter

type Filter func(c *context.Conetxt)

var _ FilterBuilder = MetricsFilterBuilder

func MetricsFilterBuilder(next Filter) Filter {
	return func(c *context.Conetxt) {
		next(c)
	}
}
