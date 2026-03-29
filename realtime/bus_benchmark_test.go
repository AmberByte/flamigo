package realtime

import (
	"fmt"
	"testing"
)

type multiTopicBenchmarkEvent struct {
	topics []Topic
}

func (e multiTopicBenchmarkEvent) Topics() []Topic {
	return e.topics
}

func BenchmarkBusPublishWildcard(b *testing.B) {
	tests := []struct {
		subscribers int
		patterns    int
	}{
		{subscribers: 1, patterns: 1},
		{subscribers: 10, patterns: 1},
		{subscribers: 100, patterns: 1},
		{subscribers: 10, patterns: 5},
		{subscribers: 100, patterns: 5},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("subs=%d/patterns=%d", tt.subscribers, tt.patterns), func(b *testing.B) {
			bus := NewBus[benchmarkEvent]()
			for i := 0; i < tt.subscribers; i++ {
				bus.Subscribe(func(ctx Context, event benchmarkEvent) {
					_ = event.payload
				}, wildcardTopicOpts(tt.patterns)...)
			}

			event := benchmarkEvent{
				topic:   "region/eu-west/server-1",
				payload: "test payload",
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bus.Publish(event)
			}
		})
	}
}

func BenchmarkBusPublishMultiTopicEvent(b *testing.B) {
	tests := []struct {
		subscribers int
		topics      int
		eventTopics int
	}{
		{subscribers: 10, topics: 2, eventTopics: 2},
		{subscribers: 100, topics: 2, eventTopics: 2},
		{subscribers: 100, topics: 5, eventTopics: 5},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("subs=%d/topics=%d/event-topics=%d", tt.subscribers, tt.topics, tt.eventTopics), func(b *testing.B) {
			bus := NewBus[multiTopicBenchmarkEvent]()
			for i := 0; i < tt.subscribers; i++ {
				bus.Subscribe(func(ctx Context, event multiTopicBenchmarkEvent) {
					_ = len(event.topics)
				}, benchmarkTopicOpts(tt.topics)...)
			}

			event := multiTopicBenchmarkEvent{
				topics: benchmarkEventTopics(tt.eventTopics),
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bus.Publish(event)
			}
		})
	}
}

func BenchmarkSubscriptionMutation(b *testing.B) {
	b.Run("add-remove-topic", func(b *testing.B) {
		bus := NewBus[benchmarkEvent]()
		subscription := bus.Subscribe(func(ctx Context, event benchmarkEvent) {
			_ = event.payload
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			topic := ParseTopic(fmt.Sprintf("dynamic/%d", i))
			subscription.AddTopic(topic)
			subscription.RemoveTopic(topic)
		}
	})

	b.Run("subscribe-all-after-topics", func(b *testing.B) {
		bus := NewBus[benchmarkEvent]()
		subscription := bus.Subscribe(func(ctx Context, event benchmarkEvent) {
			_ = event.payload
		}, WithTopics(ParseTopic("server/1"), ParseTopic("server/2"), ParseTopic("server/3")))

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			subscription.SubscribeAll()
		}
	})
}

func wildcardTopicOpts(patterns int) []SubscribeOpt {
	opts := make([]SubscribeOpt, 0, patterns)
	for i := 0; i < patterns; i++ {
		opts = append(opts, WithTopic(ParseTopic(fmt.Sprintf("region/*/server-%d", i+1))))
	}
	return opts
}

func benchmarkEventTopics(count int) []Topic {
	topics := make([]Topic, 0, count)
	for i := 0; i < count; i++ {
		topics = append(topics, NewTopic(fmt.Sprintf("topic-%d", i)))
	}
	return topics
}

func benchmarkTopicOpts(topics int) []SubscribeOpt {
	opts := make([]SubscribeOpt, 0, topics)
	for j := 0; j < topics; j++ {
		opts = append(opts, WithTopic(ParseTopic(fmt.Sprintf("topic-%d", j))))
	}
	return opts
}

type benchmarkEvent struct {
	topic    string
	payload  string
	receiver bool
}

func (e benchmarkEvent) Topics() []Topic {
	return []Topic{NewTopic(e.topic)}
}

func BenchmarkBusPublish(b *testing.B) {
	tests := []struct {
		subscribers int
		topics      int
	}{
		{subscribers: 1, topics: 1},
		{subscribers: 10, topics: 1},
		{subscribers: 100, topics: 1},
		{subscribers: 10, topics: 10},
		{subscribers: 100, topics: 10},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("subs=%d/topics=%d", tt.subscribers, tt.topics), func(b *testing.B) {
			bus := NewBus[benchmarkEvent]()

			// Setup subscribers
			for i := 0; i < tt.subscribers; i++ {
				bus.Subscribe(func(ctx Context, event benchmarkEvent) {
					// Simulate some work
					_ = event.payload
				}, benchmarkTopicOpts(tt.topics)...)
			}

			event := benchmarkEvent{
				topic:    "topic-0",
				payload:  "test payload",
				receiver: true,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bus.Publish(event)
			}
		})
	}
}

func BenchmarkBusPublishSync(b *testing.B) {
	tests := []struct {
		subscribers int
		topics      int
	}{
		{subscribers: 1, topics: 1},
		{subscribers: 10, topics: 1},
		{subscribers: 100, topics: 1},
		{subscribers: 10, topics: 10},
		{subscribers: 100, topics: 10},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("subs=%d/topics=%d", tt.subscribers, tt.topics), func(b *testing.B) {
			bus := NewBus[benchmarkEvent]()

			// Setup subscribers
			for i := 0; i < tt.subscribers; i++ {
				bus.Subscribe(func(ctx Context, event benchmarkEvent) {
					// Simulate some work
					_ = event.payload
				}, benchmarkTopicOpts(tt.topics)...)
			}

			event := benchmarkEvent{
				topic:    "topic-0",
				payload:  "test payload",
				receiver: true,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bus.PublishSync(event)
			}
		})
	}
}

func BenchmarkBusSubscribe(b *testing.B) {
	bus := NewBus[benchmarkEvent]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sub := bus.Subscribe(func(ctx Context, event benchmarkEvent) {
			// Simulate some work
			_ = event.payload
		}, WithAllTopics())
		sub.Cancel()
	}
}
