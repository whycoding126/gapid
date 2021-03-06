// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package image

import (
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/gapis/database"
)

// Resolve returns the byte array holding the resized image for the
// ResizeResolver request.
// TODO: Can this be moved to the resolve package?
func (r *ResizeResolvable) Resolve(ctx log.Context) (interface{}, error) {
	data, err := database.Resolve(ctx, r.Data.ID())
	if err != nil {
		return nil, err
	}

	return r.Format.Resize(data.([]byte),
		int(r.SrcWidth), int(r.SrcHeight),
		int(r.DstWidth), int(r.DstHeight),
	)
}
