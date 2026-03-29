package realtime

import "context"

type SubscribeOpt func(*subscribeConfig)

type subscribeConfig struct {
	all          bool
	topics       []Topic
	lifecycleCtx context.Context
}

func WithTopic(topic Topic) SubscribeOpt {
	return func(cfg *subscribeConfig) {
		if cfg.all {
			return
		}
		cfg.topics = append(cfg.topics, topic)
	}
}

func WithTopics(topics ...Topic) SubscribeOpt {
	return func(cfg *subscribeConfig) {
		if cfg.all {
			return
		}
		cfg.topics = append(cfg.topics, topics...)
	}
}

func WithAllTopics() SubscribeOpt {
	return func(cfg *subscribeConfig) {
		cfg.all = true
		cfg.topics = nil
	}
}

func WithLifecycleContext(ctx context.Context) SubscribeOpt {
	return func(cfg *subscribeConfig) {
		cfg.lifecycleCtx = ctx
	}
}

func newSubscribeConfig(opts ...SubscribeOpt) subscribeConfig {
	cfg := subscribeConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
