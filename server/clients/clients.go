/*
 * Copyright 2020 The Yorkie Authors. All rights reserved.
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

package clients

import (
	"context"
	"errors"

	"github.com/yorkie-team/yorkie/api/types"
	"github.com/yorkie-team/yorkie/pkg/document/time"
	"github.com/yorkie-team/yorkie/server/backend"
	"github.com/yorkie-team/yorkie/server/backend/database"
)

var (
	// ErrInvalidClientKey is returned when the given Key is not valid ClientKey.
	ErrInvalidClientKey = errors.New("invalid client key")

	// ErrInvalidClientID is returned when the given Key is not valid ClientID.
	ErrInvalidClientID = errors.New("invalid client id")
)

// Activate activates the given client.
func Activate(
	ctx context.Context,
	be *backend.Backend,
	project *types.Project,
	clientKey string,
) (*database.ClientInfo, error) {
	return be.DB.ActivateClient(ctx, project.ID, clientKey)
}

// Deactivate deactivates the given client.
func Deactivate(
	ctx context.Context,
	be *backend.Backend,
	project *types.Project,
	actorID *time.ActorID,
) (*database.ClientInfo, error) {
	return be.DB.DeactivateClient(ctx, project.ID, types.IDFromActorID(actorID))
}

// FindClientInfo finds the client with the given id.
func FindClientInfo(
	ctx context.Context,
	be *backend.Backend,
	project *types.Project,
	clientID *time.ActorID,
) (*database.ClientInfo, error) {
	return be.DB.FindClientInfoByID(
		ctx,
		project.ID,
		types.IDFromActorID(clientID),
	)
}