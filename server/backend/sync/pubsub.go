/*
 * Copyright 2021 The Yorkie Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sync

import (
	"sync"
	gotime "time"

	"github.com/rs/xid"

	"github.com/yorkie-team/yorkie/api/types"
	"github.com/yorkie-team/yorkie/pkg/document/time"
)

const (
	// publishTimeout is the timeout for publishing an event.
	publishTimeout = 100 * gotime.Millisecond
)

// Subscription represents a subscription of a subscriber to documents.
type Subscription struct {
	id         string
	subscriber *time.ActorID
	mu         sync.Mutex
	closed     bool
	events     chan DocEvent
}

// NewSubscription creates a new instance of Subscription.
func NewSubscription(subscriber *time.ActorID) *Subscription {
	return &Subscription{
		id:         xid.New().String(),
		subscriber: subscriber,
		events:     make(chan DocEvent, 1),
		closed:     false,
	}
}

// ID returns the id of this subscription.
func (s *Subscription) ID() string {
	return s.id
}

// DocEvent represents events that occur related to the document.
type DocEvent struct {
	Type           types.DocEventType
	Publisher      *time.ActorID
	DocumentRefKey types.DocRefKey
	Body           types.DocEventBody
}

// Events returns the DocEvent channel of this subscription.
func (s *Subscription) Events() chan DocEvent {
	return s.events
}

// Subscriber returns the subscriber of this subscription.
func (s *Subscription) Subscriber() *time.ActorID {
	return s.subscriber
}

// Close closes all resources of this Subscription.
func (s *Subscription) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		s.closed = true
		close(s.events)
	}
}

// Publish publishes the given event to the subscriber.
func (s *Subscription) Publish(event DocEvent) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return false
	}

	// NOTE(hackerwins): When a subscription is being closed by a subscriber,
	// the subscriber may not receive messages.
	select {
	case s.Events() <- event:
		return true
	case <-gotime.After(publishTimeout):
		return false
	}
}
