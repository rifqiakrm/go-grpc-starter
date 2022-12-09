package pubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// Subscriber is an interface that defines the methods that a pubsub receiver must implement.
type Subscriber interface {
	SubscriptionName() string
	ProcessMessage(context.Context, *pubsub.Message)
}

// Publisher is an interface that defines the methods that a pubsub publisher must implement.
type Publisher interface {
	Send(ctx context.Context, topicName string, data interface{}, attributes interface{}) error
}

// PubSub is a pubsub engine.
type PubSub struct {
	*pubsub.Client
}

// NewPubSub creates an instance of PubSub.
// It needs three parameters.
// The first parameter is project ID where the topic resides.
// The second parameter is credential file. Usually, it is a Service Account.
func NewPubSub(projectID string, credentialFile *string) *PubSub {
	// subscribe to pubsub topic
	var client *pubsub.Client
	var err error

	if credentialFile != nil {
		client, err = pubsub.NewClient(context.Background(), projectID, option.WithCredentialsFile(*credentialFile))
	} else {
		client, err = pubsub.NewClient(context.Background(), projectID)
	}
	if err != nil {
		panic(err)
	}

	return &PubSub{
		Client: client,
	}
}

// StartSubscriptions starts pub sub engine to receive subs
func (ps *PubSub) StartSubscriptions(subscribers ...Subscriber) error {
	for idx := range subscribers {
		go func(snh Subscriber) {
			if err := ps.Client.Subscription(snh.SubscriptionName()).Receive(context.Background(),
				snh.ProcessMessage); err != nil {
				log.Println("GooglePubSub-StartSubscriptions: Error subscribe " + snh.SubscriptionName())
				log.Println(err)
				panic(err)
			}
		}(subscribers[idx])
	}

	return nil
}
