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

package vulkan

import (
	"fmt"

	"github.com/google/gapid/core/data/pod"
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/gapis/atom"
	"github.com/google/gapid/gapis/atom/transform"
	"github.com/google/gapid/gapis/gfxapi"
	"github.com/google/gapid/gapis/replay"
	"github.com/google/gapid/gapis/replay/builder"
	"github.com/google/gapid/gapis/replay/value"
)

// findIssues is an atom transform that detects issues when replaying the
// stream of atoms. Any issues that are found are written to all the chans in
// the slice out. Once the last issue is sent (if any) all the chans in out are
// closed.
// NOTE: right now this transform is just used to close chans passed in requests.
type findIssues struct {
	out []chan<- replay.Issue
}

// reportTo adds the chan c to the list of issue listeners.
func (t *findIssues) reportTo(c chan<- replay.Issue) { t.out = append(t.out, c) }

func (t *findIssues) Transform(ctx log.Context, i atom.ID, a atom.Atom, out transform.Writer) {
	out.MutateAndWrite(ctx, i, a)
}

func (t *findIssues) Flush(ctx log.Context, out transform.Writer) {
	out.MutateAndWrite(ctx, atom.NoID, replay.Custom(func(ctx log.Context, s *gfxapi.State, b *builder.Builder) error {
		// Since the PostBack function is called before the replay target has actually arrived at the post command,
		// we need to actually write some data here. r.Uint32() is what actually waits for the replay target to have
		// posted the data in question. If we did not do this, we would shut-down the replay as soon as the second-to-last
		// Post had occurred, which may not be anywhere near the end of the stream.
		code := uint32(0xe11de11d)
		b.Push(value.U32(code))
		b.Post(b.Buffer(1), 4, func(r pod.Reader, err error) error {
			if err != nil {
				t.out = nil
				return err
			}
			if r.Uint32() != code {
				return fmt.Errorf("Flush did not get expected EOS code")
			}
			for _, c := range t.out {
				close(c)
			}
			t.out = nil
			return err
		})
		return nil
	}))
}
