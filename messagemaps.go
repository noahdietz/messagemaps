// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package messagemaps

import (
	"flag"
	"fmt"
	"io"
	"os"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	Flags flag.FlagSet

	outFile       *string
	resourcesOnly *bool
)

func init() {
	outFile = Flags.String("out_file", "", "")
	resourcesOnly = Flags.Bool("resources_only", false, "")
}

func Analyze(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	var out io.Writer
	if f := *outFile; f != "" {
		file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		out = file
	} else {
		out = os.Stderr
	}
	w := func(f string, args ...any) {
		out.Write([]byte(fmt.Sprintf(f, args...) + "\n"))
	}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		messages := collectMessages(file.Messages)
		for _, r := range messages {
			for _, f := range r.Fields {
				d := f.Desc
				if !d.IsMap() || d.MapValue().Kind() != protoreflect.MessageKind {
					continue
				}
				w("%s: map with message value type %q", d.FullName(), d.MapValue().Message().FullName())
			}
		}
	}

	return nil
}

func collectMessages(messages []*protogen.Message) (resources []*protogen.Message) {
	for _, m := range messages {
		res := (proto.GetExtension(m.Desc.Options(), annotations.E_Resource)).(*annotations.ResourceDescriptor)
		if !*resourcesOnly || res.GetType() != "" {
			resources = append(resources, m)
		}
		resources = append(resources, collectMessages(m.Messages)...)
	}
	return
}
