package realtime

type SubscribeOpt func(*subscribeConfig)

type subscribeConfig struct {
	all    bool
	topics []string
}

func WithTopic(topic string) SubscribeOpt {
	return func(cfg *subscribeConfig) {
		if cfg.all {
			return
		}
		cfg.topics = append(cfg.topics, topic)
	}
}

func WithTopics(topics ...string) SubscribeOpt {
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

func newSubscribeConfig(opts ...SubscribeOpt) subscribeConfig {
	cfg := subscribeConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
