/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package compose

import (
	"context"

	"github.com/docker/compose-cli/api/compose"

	"github.com/compose-spec/compose-go/types"
	"golang.org/x/sync/errgroup"
)

func (s *composeService) Start(ctx context.Context, project *types.Project, consumer compose.LogConsumer) error {
	var group *errgroup.Group
	if consumer != nil {
		eg, err := s.attach(ctx, project, consumer)
		if err != nil {
			return err
		}
		group = eg
	}

	err := InDependencyOrder(ctx, project, func(c context.Context, service types.ServiceConfig) error {
		return s.startService(ctx, project, service)
	})
	if err != nil {
		return err
	}
	if group != nil {
		return group.Wait()
	}
	return nil
}
