/*
 * Copyright 2019-present Ciena Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package commands

import (
	"context"
  "fmt"
	"github.com/golang/protobuf/ptypes/empty"
	flags "github.com/jessevdk/go-flags"
	"github.com/opencord/bossctl/pkg/format"
	"github.com/opencord/bossctl/pkg/model"
	"github.com/opencord/voltha-protos/v5/go/voltha"
)

const (
	DEFAULT_OUTPUT_FORMAT = "table{{ .Id }}\t{{ .Vendor }}\t{{ .Type }}\t{{ .Endpoint }}\t{{ .Version }}\t{{ .CurrentReplica }}\t{{ .TotalReplicas }}\t{{ gosince .LastCommunication}}"
  ETCD_OUTPUT_FORMAT =`
  ETCD_GET_VALUE_VERSION1 : {{.ver1}}
  ETCD_GET_VALUE_VERSION2 : {{.ver2}}
  `
)

type AdapterList struct {
	ListOutputOptions
}
type etcdList struct {
	ListOutputOptions
  Args struct {
		Ids []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type AdapterOpts struct {
	List AdapterList `command:"list"`
  EtcdList etcdList `command:"etcdList"`
}

var adapterOpts = AdapterOpts{}

func RegisterAdapterCommands(parent *flags.Parser) {
	if _, err := parent.AddCommand("adapter", "adapter commands", "Commands to query and manipulate VOLTHA adapters", &adapterOpts); err != nil {
		Error.Fatalf("Unexpected error while attempting to register adapter commands : %s", err)
	}
}

func (options *AdapterList) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	adapters, err := client.ListAdapters(ctx, &empty.Empty{})
	if err != nil {
		return err
	}
	data := make([]model.AdapterInstance, len(adapters.Items))
	for i, item := range adapters.Items {
		data[i].PopulateFrom(item)
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("adapter-list", "format", DEFAULT_OUTPUT_FORMAT)
	}
	if options.Quiet {
		outputFormat = "{{.Id}}"
	}
	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("adapter-list", "order", "")
	}

	// TODO: lastCommunication ends up formatted as `seconds:1589415656 nanos:775740000`
	//   need to think through where to do presentation formatting.

	result := CommandResult{
		Format:    format.Format(outputFormat),
		Filter:    options.Filter,
		OrderBy:   orderBy,
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      data,
	}
	GenerateOutput(&result)

	return nil
}
func (options *etcdList) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	var lastErr error
	for _, i := range options.Args.Ids {
		ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
		defer cancel()

		id := voltha.ID{Id: string(i)}

		returnEtcd, err := client.GetEtcdList(ctx, &id)
		if err != nil {
			Error.Printf("Error while enabling '%s': %s\n", i, err)
			lastErr = err
			continue
		}
    fmt.Println(returnEtcd)
//    outputFormat := GetCommandOptionWithDefault("device-ports", "format", ETCD_OUTPUT_FORMAT)
//    orderBy := GetCommandOptionWithDefault("device-ports", "order", "")
//    result := CommandResult{
//		  Format : format.Format(outputFormat),
//		  OrderBy : orderBy,
//		  Data : returnEtcd,
//	  }
//	GenerateOutput(&result)
	}

	if lastErr != nil {
		return NoReportErr
	}
	return nil
}
