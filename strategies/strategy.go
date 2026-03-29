package strategies

type Strategy[CTX Context] func(ctx CTX)

type Registry[CTX Context] interface {
	Register(action string, fn Strategy[CTX]) error
	Invoke(action string, ctx CTX) Result
}

type Router[CTX Context] interface {
	Register(namespace string, registry Registry[CTX]) error
	Invoke(ctx CTX) Result
}

type AppStrategy = Strategy[Context]

type AppRegistry = Registry[Context]
type AppRouter = Router[Context]
