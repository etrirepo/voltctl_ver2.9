# bossctl, a CLI for VOLTHA

This repository contains `bossctl`, a CLI tool for managing and operating
VOLTHA components.

It functions similarly to the `docker` CLI or kubernetes `kubectl` CLI, in that
it's a simple standalone control application which can perform various
functions and has flexible and customizable output formats either as a table or
`JSON`.

## Build / Install

To install the `bossctl` command, downloads are available for multiple
platforms and architectures from the [github releases
page](https://github.com/opencord/bossctl/releases), or you can compile your
own copy by installing Go 1.13.x, checking out the code and running `make
build`.

## Shell Completion

`bossctl` supports shell completion for the `bash` shell. To enable shell
Completion you can use the following command on *most* \*nix based system.

```shell
source <(bossctl completion bash)
```

If you are running an older bash 3.x shell (default on MacOS), then you can try
the following command:

```shell
source /dev/stdin <<<"$(bossctl completion bash)"
```

If you which to make `bash` shell completion automatic when you login to your
account you can append the output of `bossctl completion bash` to your
`$HOME/.bashrc`:

```shell
bossctl completion base >> $HOME/.bashrc
```

## Configuration file

`bossctl` stores it's configuration file in `~/.volt/config`. An example of the
configuration file can be found in the `bossctl.config` file in this repo.

## Usage and Commands

```shell
$ ./bossctl -h
Usage:
  bossctl [OPTIONS] <command>

Global Options:
  -c, --config=FILE                     Location of client config file [$VOLTCONFIG]
  -s, --server=SERVER:PORT              IP/Host and port of VOLTHA
  -k, --kafka=SERVER:PORT               IP/Host and port of Kafka
  -e, --kvstore=SERVER:PORT             IP/Host and port of KV store (etcd) [$KVSTORE]
  -d, --debug                           Enable debug mode
  -t, --timeout=DURATION                API call timeout duration
      --tls                             Use TLS
      --tlscacert=CA_CERT_FILE          Trust certs signed only by this CA
      --tlscert=CERT_FILE               Path to TLS vertificate file
      --tlskey=KEY_FILE                 Path to TLS key file
      --tlsverify                       Use TLS and verify the remote
  -8, --k8sconfig=FILE                  Location of Kubernetes config file [$KUBECONFIG]
      --kvstoretimeout=DURATION         timeout for calls to KV store [$KVSTORE_TIMEOUT]
  -o, --command-options=FILE            Location of command options default configuration file [$VOLTCTL_COMMAND_OPTIONS]
  -m, --maxcallrecvmsgsize=SIZE         Max GRPC Client request size limit in bytes (eg: 4MB)

Help Options:
  -h, --help                            Show this help message

Available commands:
  adapter        adapter commands
  completion     generate shell compleition
  component      component instance commands
  config         generate bossctl configuration
  device         device commands
  devicegroup    device group commands
  event          event commands
  log            log config commands
  logicaldevice  logical device commands
  message        message commands
  version        display version
```

Help specific to each command can be found by running `volctl <command> -h`.

### Changing the command output format

Each command has a default output table format. This can be overridden from the
command line using the `bossctl --format=...` option. The specification of the
format is roughly equivalent to the `docker` or `kubectl` command. If the
prefix `table` is specified a table with headers will be displayed, else each
line will be output as specified.

The output of a command may also be written as `JSON` or `YAML` by using the
`--outputas` or `-o` command line option. Valid values for this options are
`table`, `json`, or `yaml`.

### Overriding Default Command Format and Order

The default format and ordering of commands can be overridden (specified) by
the command line options, but they can also be set via a configuration file so
that the overrides don't have to be specified on each invocation.

By default the file `~/.volt/command_options` is loaded, but the file used can
also be specified by the environment variable `VOLTCTL_COMMAND_OPTIONS` or via
the command line arguments.

A sample of this file is include in the repository as
`bossctl_command_options.config`.

### Examples

```shell
$ bossctl adapter list
ID                   VENDOR            VERSION       SINCELASTCOMMUNICATION
openolt              VOLTHA OpenOLT    2.3.0-dev     NEVER
brcm_openomci_onu    VOLTHA OpenONU    2.3.0-dev     32m10s
```

```shell
$ bossctl adapter list --outputas json
[{"Id":"openolt","Vendor":"VOLTHA OpenOLT","Version":"2.3.0-dev","LogLevel":"","LastCommunication":"NEVER","SinceLastCommunication":"NEVER"},{"Id":"brcm_openomci_onu","Vendor":"VOLTHA OpenONU","Version":"2.3.0-dev","LogLevel":"DEBUG","LastCommunication":"2020-04-04T20:48:59Z","SinceLastCommunication":"1s"}]
```

After piping through `python -m json.tool`:

```json
[
    {
        "Id": "openolt",
        "LastCommunication": "NEVER",
        "LogLevel": "",
        "SinceLastCommunication": "NEVER",
        "Vendor": "VOLTHA OpenOLT",
        "Version": "2.3.0-dev"
    },
    {
        "Id": "brcm_openomci_onu",
        "LastCommunication": "2020-04-04T20:46:45Z",
        "LogLevel": "DEBUG",
        "SinceLastCommunication": "1m57s",
        "Vendor": "VOLTHA OpenONU",
        "Version": "2.3.0-dev"
    }
]
```


```shell
$ bossctl device list
ID                                      TYPE                 ROOT     PARENTID                                SERIALNUMBER    ADMINSTATE    OPERSTATUS    CONNECTSTATUS    REASON
1398f977-3630-43d2-8d3b-1ae395a95162    openolt              true     540fc38e-cf35-4d14-8b01-7760acecefaa    BBSIM_OLT_0     ENABLED       ACTIVE        REACHABLE
5bacc996-b922-41fc-8ddc-a92d29729955    brcm_openomci_onu    false    1398f977-3630-43d2-8d3b-1ae395a95162    BBSM00000001    ENABLED       ACTIVE        REACHABLE        omci-flows-pushed
```

```shell
$ bossctl device list --format 'table{{.Id}}\t{{.SerialNumber}}\t{{.ConnectStatus}}'
ID                                      SERIALNUMBER    CONNECTSTATUS
1398f977-3630-43d2-8d3b-1ae395a95162    BBSIM_OLT_0     REACHABLE
5bacc996-b922-41fc-8ddc-a92d29729955    BBSM00000001    REACHABLE
```

```shell
$ bossctl device list --outputas json
[{"id":"d2960b6e-f963-4acb-83c9-492b3211cf6b","type":"brcm_openomci_onu","root":false,"parentid":"e2c1d2cd-c260-4285-8632-7b205aed660a","parentportno":536870912,"vendor":"OpenONU","model":"","hardwareversion":"","firmwareversion":"","serialnumber":"BBSM00000001","vendorid":"BBSM","adapter":"brcm_openomci_onu","vlan":0,"macaddress":"","address":"unknown","extraargs":"","proxyaddress":{"deviceId":"e2c1d2cd-c260-4285-8632-7b205aed660a","devicetype":"openolt","channelid":0,"channelgroup":0,"onuid":1,"onusessionid":0},"adminstate":"ENABLED","operstatus":"DISCOVERED","reason":"stopping-openomci","connectstatus":"UNREACHABLE","ports":[{"portno":536870912,"label":"PON port","type":"PON_ONU","adminstate":"ENABLED","operstatus":"ACTIVE","deviceid":"","peers":[{"deviceid":"e2c1d2cd-c260-4285-8632-7b205aed660a","portno":536870912}]},{"portno":16,"label":"uni-16","type":"ETHERNET_UNI","adminstate":"ENABLED","operstatus":"UNKNOWN","deviceid":"","peers":[]},{"portno":17,"label":"uni-17","type":"ETHERNET_UNI","adminstate":"ENABLED","operstatus":"DISCOVERED","deviceid":"","peers":[]},{"portno":18,"label":"uni-18","type":"ETHERNET_UNI","adminstate":"ENABLED","operstatus":"DISCOVERED","deviceid":"","peers":[]},{"portno":19,"label":"uni-19","type":"ETHERNET_UNI","adminstate":"ENABLED","operstatus":"DISCOVERED","deviceid":"","peers":[]}],"flows":[{"id":"8c4fd2d0f768700a","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":1000,"cookie":"~3fd5629a","inport":"16","vlanid":"0","setvlanid":"900","output":"536870912","writemetadata":"0x0384004000100000","meter":"2","tunnelid":"16"},{"id":"5d0b3499cd2bf4ac","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":1000,"cookie":"~4df91e40","inport":"536870912","vlanid":"900","metadata":"0x0000000000000010","setvlanid":"0","output":"16","writemetadata":"0x0000004000000000","meter":"2"},{"id":"21a5ad60293e6c60","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":10000,"cookie":"~ba31a4f2","inport":"16","ethtype":"0x0800","ipproto":"17","udpsrc":"68","dstsrc":"67","setvlanid":"900","pushvlanid":"0x8100","output":"536870912","writemetadata":"0x0000004000000000","meter":"2","tunnelid":"16"}]},{"id":"e2c1d2cd-c260-4285-8632-7b205aed660a","type":"openolt","root":true,"parentid":"28d1128f-7d9a-48fa-b60b-d96e1491d92a","parentportno":0,"vendor":"BBSim","model":"asfvolt16","hardwareversion":"","firmwareversion":"","serialnumber":"BBSIM_OLT_0","vendorid":"","adapter":"openolt","vlan":0,"macaddress":"0a:0a:0a:0a:0a:00","address":"bbsim.voltha.svc:50060","extraargs":"","adminstate":"DISABLED","operstatus":"UNKNOWN","reason":"","connectstatus":"REACHABLE","ports":[{"portno":1048576,"label":"nni-1048576","type":"ETHERNET_NNI","adminstate":"ENABLED","operstatus":"ACTIVE","deviceid":"","peers":[]},{"portno":536870912,"label":"pon-536870912","type":"PON_OLT","adminstate":"ENABLED","operstatus":"DISCOVERED","deviceid":"","peers":[{"deviceid":"d2960b6e-f963-4acb-83c9-492b3211cf6b","portno":536870912}]}],"flows":[{"id":"e1746c5320441c57","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":10000,"cookie":"~f81586a7","inport":"1048576","ethtype":"0x0800","ipproto":"17","udpsrc":"67","dstsrc":"68","output":"CONTROLLER"},{"id":"12f8e0237d36dcab","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":10000,"cookie":"~ce6c3527","inport":"1048576","ethtype":"0x88cc","output":"CONTROLLER"},{"id":"35f0a5d7315c8b8a","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":1000,"cookie":"~986cca9a","inport":"536870912","vlanid":"900","setvlanid":"900","pushvlanid":"0x8100","output":"1048576","writemetadata":"0x0000004000000000","meter":"2","tunnelid":"16"},{"id":"755a065fb691c418","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":1000,"cookie":"~531d5ec9","inport":"1048576","vlanid":"900","metadata":"0x0000000000000384","popvlan":"yes","output":"536870912","writemetadata":"0x0384004000000010","meter":"2","tunnelid":"16"},{"id":"e11f009524a53eb2","tableid":0,"durationsec":0,"durationnsec":0,"idletimeout":0,"hardtimeout":0,"packetcount":0,"bytecount":0,"priority":10000,"cookie":"~ba31a4f2","inport":"536870912","ethtype":"0x0800","vlanid":"900","ipproto":"17","udpsrc":"68","dstsrc":"67","output":"CONTROLLER","writemetadata":"0x0000004000000000","meter":"2","tunnelid":"16"}]}]
```

### Configure logs in the VOLTHA components

Check the current log levels:

```shell
$ bossctl log level list
COMPONENTNAME       PACKAGENAME    LEVEL
read-write-core     default        DEBUG
adapter-open-olt    default        DEBUG
global              default        WARN
open-flow-agent     default        DEBUG
```

Configure the global log level:

```shell
$ bossctl log level set INFO
COMPONENTNAME    PACKAGENAME    STATUS     ERROR
global           default        Success
```

Configure the default log level of a component:

```shell
$ bossctl log level set ERROR read-write-core
COMPONENTNAME      PACKAGENAME    STATUS     ERROR
read-write-core    default        Success
```

Set the log level of a specific package within a component:

```shell
$ bossctl log level set WARN adapter-open-olt#github.com/opencord/voltha-lib-go/kafka
COMPONENTNAME       PACKAGENAME                                 STATUS     ERROR
adapter-open-olt    github.com/opencord/voltha-lib-go/kafka    Success
```

Get list of valid log packages for a component:

```shell
$ bossctl log package list adapter-open-olt
COMPONENTNAME       PACKAGENAME
adapter-open-olt    default
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/adapters/common
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/config
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/db
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/db/kvstore
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/flows
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/kafka
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/ponresourcemanager
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/probe
adapter-open-olt    github.com/opencord/voltha-lib-go/v3/pkg/techprofile
adapter-open-olt    github.com/opencord/voltha-openolt-adapter/internal/pkg/config
adapter-open-olt    github.com/opencord/voltha-openolt-adapter/internal/pkg/core
adapter-open-olt    github.com/opencord/voltha-openolt-adapter/internal/pkg/resourcemanager
adapter-open-olt    main
```

If you do not see packages in your output it is possible you are not retrieving the data from the correct path in ``ETCD``.
The data is stored according to the helm chart defined value for `kv_store_data_prefix`.

If you deployed according to the README in the [voltha-helm-charts README]() the value for `kv_store_data_prefix` is
`service/voltha/{{ .Release.Name }}_{{ .Values.global.stack_name }}`, e.g `service/voltha/voltha1_voltha1`.

If you are unsure about the `kv_store_data_prefix` you can check through one of the pods, e.g. `rw-core`.
Make sure to use the correct namespace for `-n` by checking it with `kubectl get pods -A` and then issue:
```shell
$ kubectl describe -n voltha1 $(kubectl -n voltha1 get pods -l app=rw-core -o name) | grep KV
KV_STORE_DATAPATH_PREFIX:  service/voltha/voltha1_voltha1
```

You can set the path for any `bossctl` command by exporting the `KV_STORE_DATAPATH_PREFIX` variable,
e.g. `KV_STORE_DATAPATH_PREFIX=service/voltha/voltha1_voltha1`.

A full example command is:
```shell
KV_STORE_DATAPATH_PREFIX=service/voltha/voltha1_voltha1 bossctl log package list adapter-open-olt
```

### Configuring the message size

When you run VOLTHA with a high number of OLTs/ONUs is possible that the gRPC response exceeds the default 4MB in size.
This will be addressed on the server side in a future release, but you have the option to exceed the allowed response size
using the `-m, --maxcallrecvmsgsize=SIZE` option.

The common error message when this happens is:
```shell
ERROR: RESOURCEEXHAUSTED: grpc: received message larger than max (4241141 vs. 4194304)
```

If that happens, retry the command passing the -m option, eg:
```shell
$ bossctl device list -m 8M
```
