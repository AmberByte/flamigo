package strategies

type Strategy[CTX Context] func(ctx CTX)

type Registry[CTX Context] interface {
	Register(topic string, fn Strategy[CTX]) error
	Use(ctx CTX) StrategyResult
}

type AppStrategy = Strategy[Context]

type AppRegistry = Registry[Context]
