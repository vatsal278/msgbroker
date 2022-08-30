# Message Broker

[Message Broker](https://en.wikipedia.org/wiki/Message_broker#:~:text=A%20message%20broker%20(also%20known,messaging%20protocol%20of%20the%20receiver.)) an independent service that follows the [Observer System Design Pattern](https://refactoring.guru/design-patterns/observer) in virtue of using a publisher-subscriber methodology.

This acts as a webhook that notifies all subscribers via their own HTTP specification.

## Overview
A publisher can push data onto a topic which will be forwarded to all the subscribers subscribed to the topic. This is similar to following a topic on social media sites. Whenever there is some post for a topic, all users who are subscribed to the topic gets a notification. There can be multiple publishers publishing messages to a topic and multiple subscribers subscribed to a single topic.

This infra can be started as an independent service(server) and should be containerizable for independent scale-up/down.

## API Specification
It primarily has 3 `HTTP POST` endpoints:
1. `register` a publisher for a particular topic
2. `register` a subscriber to a particular topic
3. `publish` a message to a topic which in turn pushes this message to all the subscribers subscribed to the topic.

Additionally, there may be endpoints to `de-register` a publisher and a subscriber, which can be thought of as an enhancement but provision must be made to support this feature without additional change to existing endpoint logic(request/response).

### Subscriber
- A subscriber can be thought of as a server, which can have multiple endpoints.
- While registering a subscriber, the request to this infra layer should specify the endpoint details like the URL on the subscriber server, HTTP method configured for the endpoint on the subscriber server, etc.

### Puiblisher
- A publisher can be thought of as a HTTP client that can make HTTP calls to this infra layer.
- Only publishers registered to publish on a topic can push a message on the topic. This will prevent anonymous push to a topic from non-registered publishers.

## Precaution
- Avoid duplication of subscribers for a single topic that may arise when replaying the subscriber registration call. This will prevent a case where a subscriber receives same messages pushed on a single topic.

## Consideration
- There need not be any retry of message push if a subscriber is down or on failure. Thus, on publishing messages, the subscriber endpoints will be hit with the data, pushed by a registered publisher on a particular topic, and need not be checked for the response from the subscriber’s endpoint. This will be a fire-and-forget call.
- Messages once published need not be stored.
- The infra service must use the [base golang net/http](https://pkg.go.dev/net/http) package without any third party libraries.
- There need **NOT** be usage of any external data storage solutions like MySQL, MongoDB, Redis, etc. unless absolutely required.
- All publisher & subscriber records can be stored in-memory for faster retrieval. 