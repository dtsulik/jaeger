// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adjuster

import (
	"github.com/jaegertracing/jaeger/model"
)

var spanTagsToMove = map[string]struct{}{
	"otel.library.name": {},
}

// SpanTagsToProcessAdjuster moves certain tags from span.tags to
// span.process.tags that should be there. This should ideally be
// fixed in upstream OTEL.
func SpanTagsToProcessAdjuster() Adjuster {
	return Func(func(trace *model.Trace) (*model.Trace, error) {
		for _, span := range trace.Spans {
			processTagsMap := make(map[string]struct{})
			for _, check := range span.Process.Tags {
				processTagsMap[check.Key] = struct{}{}
			}

			index := 0
			for _, tag := range span.Tags {
				if _, ok := spanTagsToMove[tag.Key]; ok {
					if _, exists := processTagsMap[tag.Key]; !exists {
						span.Process.Tags = append(span.Process.Tags, tag)
						continue
					}
				}
				span.Tags[index] = tag
				index++
			}
			span.Tags = span.Tags[:index]
		}
		return trace, nil
	})
}
