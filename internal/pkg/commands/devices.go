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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	flags "github.com/jessevdk/go-flags"
	"github.com/opencord/bossctl/pkg/format"
	"github.com/opencord/voltha-protos/v5/go/common"
	"github.com/opencord/voltha-protos/v5/go/extension"
	"github.com/opencord/voltha-protos/v5/go/voltha"
	"github.com/opencord/voltha-protos/v5/go/bossopenolt"
)

const (
	DEFAULT_DEVICE_FORMAT         = "table{{ .Id }}\t{{.Type}}\t{{.Root}}\t{{.ParentId}}\t{{.SerialNumber}}\t{{.AdminState}}\t{{.OperStatus}}\t{{.ConnectStatus}}\t{{.Reason}}"
	DEFAULT_DEVICE_PORTS_FORMAT   = "table{{.PortNo}}\t{{.Label}}\t{{.Type}}\t{{.AdminState}}\t{{.OperStatus}}\t{{.DeviceId}}\t{{.Peers}}"
	DEFAULT_DEVICE_INSPECT_FORMAT = `ID: {{.Id}}
  TYPE:          {{.Type}}
  ROOT:          {{.Root}}
  PARENTID:      {{.ParentId}}
  SERIALNUMBER:  {{.SerialNumber}}
  VLAN:          {{.Vlan}}
  ADMINSTATE:    {{.AdminState}}
  OPERSTATUS:    {{.OperStatus}}
  CONNECTSTATUS: {{.ConnectStatus}}`
	DEFAULT_DEVICE_PM_CONFIG_GET_FORMAT         = "table{{.DefaultFreq}}\t{{.Grouped}}\t{{.FreqOverride}}"
	DEFAULT_DEVICE_PM_CONFIG_METRIC_LIST_FORMAT = "table{{.Name}}\t{{.Type}}\t{{.Enabled}}\t{{.SampleFreq}}"
	DEFAULT_DEVICE_PM_CONFIG_GROUP_LIST_FORMAT  = "table{{.GroupName}}\t{{.Enabled}}\t{{.GroupFreq}}"
	DEFAULT_DEVICE_VALUE_GET_FORMAT             = "table{{.Name}}\t{{.Result}}"
	DEFAULT_DEVICE_IMAGE_LIST_GET_FORMAT        = "table{{.Name}}\t{{.Url}}\t{{.Crc}}\t{{.DownloadState}}\t{{.ImageVersion}}\t{{.LocalDir}}\t{{.ImageState}}\t{{.FileSize}}"
	ONU_IMAGE_LIST_FORMAT                       = "table{{.Version}}\t{{.IsCommited}}\t{{.IsActive}}\t{{.IsValid}}\t{{.ProductCode}}\t{{.Hash}}"
	ONU_IMAGE_STATUS_FORMAT                     = "table{{.DeviceId}}\t{{.ImageState.Version}}\t{{.ImageState.DownloadState}}\t{{.ImageState.Reason}}\t{{.ImageState.ImageState}}\t"
	DEFAULT_DEVICE_GET_PORT_STATUS_FORMAT       = `
  TXBYTES:		{{.TxBytes}}
  TXPACKETS:		{{.TxPackets}}
  TXERRPACKETS:		{{.TxErrorPackets}}
  TXBCASTPACKETS:	{{.TxBcastPackets}}
  TXUCASTPACKETS:	{{.TxUcastPackets}}
  TXMCASTPACKETS:	{{.TxMcastPackets}}
  RXBYTES:		{{.RxBytes}}
  RXPACKETS:		{{.RxPackets}}
  RXERRPACKETS:		{{.RxErrorPackets}}
  RXBCASTPACKETS:	{{.RxBcastPackets}}
  RXUCASTPACKETS:	{{.RxUcastPackets}}
  RXMCASTPACKETS:	{{.RxMcastPackets}}`
	DEFAULT_DEVICE_GET_UNI_STATUS_FORMAT = `
  ADMIN_STATE:          {{.AdmState}}
  OPERATIONAL_STATE:    {{.OperState}}
  CONFIG_IND:           {{.ConfigInd}}`
	DEFAULT_ONU_PON_OPTICAL_INFO_STATUS_FORMAT = `
  POWER_FEED_VOLTAGE__VOLTS:      {{.PowerFeedVoltage}}
  RECEIVED_OPTICAL_POWER__dBm:    {{.ReceivedOpticalPower}}
  MEAN_OPTICAL_LAUNCH_POWER__dBm: {{.MeanOpticalLaunchPower}}
  LASER_BIAS_CURRENT__mA:         {{.LaserBiasCurrent}}
  TEMPERATURE__Celsius:           {{.Temperature}}`
	DEFAULT_RX_POWER_STATUS_FORMAT = `
	INTF_ID: {{.IntfId}}
	ONU_ID: {{.OnuId}}
	STATUS: {{.Status}}
	FAIL_REASON: {{.FailReason}}
	RX_POWER : {{.RxPower}}`
	DEFAULT_ETHERNET_FRAME_EXTENDED_PM_COUNTERS_FORMAT = `Upstream_Drop_Events:	        {{.UDropEvents}}
Upstream_Octets:	        {{.UOctets}}
UFrames:	                {{.UFrames}}
UBroadcastFrames:	        {{.UBroadcastFrames}}
UMulticastFrames:	        {{.UMulticastFrames}}
UCrcErroredFrames:	        {{.UCrcErroredFrames}}
UUndersizeFrames:	        {{.UUndersizeFrames}}
UOversizeFrames:	        {{.UOversizeFrames}}
UFrames_64Octets:	        {{.UFrames_64Octets}}
UFrames_65To_127Octets:	        {{.UFrames_65To_127Octets}}
UFrames_128To_255Octets:	{{.UFrames_128To_255Octets}}
UFrames_256To_511Octets:	{{.UFrames_256To_511Octets}}
UFrames_512To_1023Octets:	{{.UFrames_512To_1023Octets}}
UFrames_1024To_1518Octets:	{{.UFrames_1024To_1518Octets}}
DDropEvents:	                {{.DDropEvents}}
DOctets:	                {{.DOctets}}
DFrames:	                {{.DFrames}}
DBroadcastFrames:	        {{.DBroadcastFrames}}
DMulticastFrames:	        {{.DMulticastFrames}}
DCrcErroredFrames:	        {{.DCrcErroredFrames}}
DUndersizeFrames:	        {{.DUndersizeFrames}}
DOversizeFrames:	        {{.DOversizeFrames}}
DFrames_64Octets:	        {{.DFrames_64Octets}}
DFrames_65To_127Octets:	        {{.DFrames_65To_127Octets}}
DFrames_128To_255Octets:	{{.DFrames_128To_255Octets}}
DFrames_256To_511Octets:	{{.DFrames_256To_511Octets}}
DFrames_512To_1023Octets:	{{.DFrames_512To_1023Octets}}
DFrames_1024To_1518Octets:	{{.DFrames_1024To_1518Octets}}
PmFormat:	                {{.PmFormat}}`
	DEFAULT_DEVICE_BOSS_VLAN_FORMAT = "table{{.DeviceId}}\t{{.STAG_mode}}\t{{.STAG_fields}}"
	DEFAULT_DEVICE_BOSS_GET_CONNECTION_FORMAT = "table{{.DeviceId}}\t{{.IP}}\t{{.MAC}}"
	DEFAULT_DEVICE_BOSS_GET_DEVINFO_FORMAT = "table{{.DeviceId}}\t{{.Fpga_Type}}\t{{.Fpga_ver}}\t{{.Fpga_date}}\t{{.Sw_ver}}\t{{.Sw_date}}"
	DEFAULT_DEVICE_BOSS_GET_PMDTXDIS_FORMAT = "table{{.PortNo}}\t{{.Status}}"
	DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT = "{{.Result}}"
	DEFAULT_DEVICE_BOSS_GET_PMDDEVICESTATUS_FORMAT = "table{{.PortNo}}\t{{.Loss}}\t{{.Module}}\t{{.Fault}}\t{{.Link}}"
	DEFAULT_DEVICE_BOSS_GETDEVICEPORT_FORMAT = "table{{.PortNo}}\t{{.State}}"
	DEFAULT_DEVICE_BOSS_GETMTUSIZE_FORMAT = "table{{.Mtu}}"
	DEFAULT_DEVICE_BOSS_GETMODE_FORMAT = "table{{.DeviceId}}\t{{.Mode}}"
	DEFAULT_DEVICE_BOSS_GETAGINGTime_FORMAT = "table{{.DeviceId}}\t{{.AgingTime}}"
	DEFAULT_DEVICE_BOSS_GETDEVICEMACINFO_FORMAT = "table{{.DeviceId}}\t{{.Mtu}}\t{{.VlanMode}}\t{{.AgingMode}}\t{{.AgingTime}}"
	DEFAULT_DEVICE_BOSS_SETSDNTABLE_FORMAT = "table{{.HashKey}}"
	DEFAULT_DEVICE_BOSS_GETSDNTABLE_FORMAT = "table{{.DeviceId}}\t{{.Address}}\t{{.PortId}}\t{{.Vlan}}"
	DEFAULT_DEVICE_BOSS_GETVALUE_FORMAT = "table{{.DeviceId}}\t{{.Value}}"
	DEFAULT_DEVICE_BOSS_ADDONU_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Result}}\t{{.Rate}}\t{{.VendorId}}\t{{.Vssn}}"
	DEFAULT_DEVICE_BOSS_GETSLATABLE_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Tcont}}\t{{.Type}}\t{{.Si}}\t{{.Abmin}}\t{{.Absur}}\t{{.Fec}}\t{{.Distance}}"
	DEFAULT_DEVICE_BOSS_GETONUVSSN_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Vssn}}"
	DEFAULT_DEVICE_BOSS_GETONUDISTANCE_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Distance}}"
	DEFAULT_DEVICE_BOSS_GETBURSTDELIMITER_FORMAT = "table{{.DeviceId}}\t{{.Length}}\t{{.Delimiter}}"
	DEFAULT_DEVICE_BOSS_GETBURSTPREAMBLE_FORMAT = "table{{.DeviceId}}\t{{.Length}}\t{{.Preamble}}\t{{.Repeat}}"
	DEFAULT_DEVICE_BOSS_GETBURSTVERSION_FORMAT = "table{{.DeviceId}}\t{{.Version}}\t{{.Index}}\t{{.Pontag}}"
	DEFAULT_DEVICE_BOSS_GETBURSTPROFILE_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Version}}\t{{.Index}}\t{{.DelimiterLength}}\t{{.Delimiter}}\t{{.PreambleLength}}\t{{.Preamble}}\t{{.Repeat}}\t{{.Pontag}}"
	DEFAULT_DEVICE_BOSS_GETREGISTERSTATUS_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Status}}"
	DEFAULT_DEVICE_BOSS_GETONUINFO_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Rate}}\t{{.VendorId}}\t{{.Vssn}}\t{{.Distance}}\t{{.Status}}"
	DEFAULT_DEVICE_BOSS_GETOMCISTATUS_FORMAT = "table{{.DeviceId}}\t{{.Status}}"
	DEFAULT_DEVICE_BOSS_GETOMCIDATA_FORMAT = "table{{.DeviceId}}\t{{.Control}}\t{{.Data}}"
	DEFAULT_DEVICE_BOSS_GETTOD_FORMAT = "table{{.DeviceId}}\t{{.Mode}}\t{{.Time}}"
	DEFAULT_DEVICE_BOSS_GETPMCONTROL_FORMAT = `Device ID : {{.DeviceId}}
Action : {{.Action}}
Onu Mode : {{.OnuMode}}
Transinit : {{.Transinit}}
TxInit : {{.Txinit}}`
	DEFAULT_DEVICE_BOSS_GETPMTABLE_FORMAT = `Device ID : {{.DeviceId}}
Onu ID : {{.OnuId}}
Mode : {{.Mode}}
Sleep : {{.Sleep}}
Aware : {{.Aware}}
Rxoff : {{.Rxoff}}
Hold : {{.Hold}}
Action : {{.Action}}
Status : {{.Status}}`
	DEFAULT_DEVICE_BOSS_GETSLICEBW_FORMAT = "table{{.DeviceId}}\t{{.Bw}}"
	DEFAULT_DEVICE_BOSS_SLAV2_FORMAT = "table{{.DeviceId}}\t{{.OnuId}}\t{{.Tcont}}\t{{.AllocId}}\t{{.Slice}}\t{{.Bw}}\t{{.Dba}}\t{{.Type}}\t{{.Fixed}}\t{{.Assur}}\t{{.Nogur}}\t{{.Max}}\t{{.Reach}}"

)
type DevOltConnect struct{
        ListOutputOptions
        Args struct{
                Id      DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
        }`positional-args:"yes"`
}
type GetDevicePmdStatus struct {
        ListOutputOptions
        Args struct {
                PortType string   `positional-arg-name:"PORT_TYPE" required:"yes"`
                Id       DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
                PortNo   int   `positional-arg-name:"PORT_NO" required:"yes"`
        } `positional-args:"yes"`
}

type GetTod struct{
	DeviceId string
	Mode 	 string
	Time 	 string
}
type AddOnu struct{
	DeviceId string
	OnuId	 string
	Result	 string
	Rate	 string
	VendorId string
	Vssn	 string
}
type GetValue struct {
	DeviceId string
	Value	 string
}
type GetSdnTable struct{
	DeviceId string
	Address  string
	PortId   string
	Vlan     string
}
type SetSdnTable struct {
	HashKey string
}
type GetDeviceMacInfo struct{
	DeviceId string
	Mtu string
	VlanMode string
	AgingMode string
	AgingTime string
}
type GetDeviceVlan struct{
	DeviceId string
	STAG_mode string
	STAG_fields string
}
type GetOltConnect struct{
	DeviceId string
	IP	string
	MAC	string
}
type GetDevInfo struct{
	DeviceId  string
	Fpga_Type string
	Fpga_ver  string
	Fpga_date string
	Sw_ver	  string
	Sw_date   string
}
type GetResponse struct{
	Result string
}
type GetPmdTxDis struct{
	PortNo string
	Status string
}
type GetPmdDeviceStatus struct{
	PortNo string
	Loss   string
	Module string
	Fault  string
	Link   string
}
type GetDevicePort struct{
	PortNo string
	State  string
}
type GetMtuSize struct{
	Mtu string
}
type GetMode struct{
	DeviceId string
	Mode string
}
type GetAgingTime struct{
	DeviceId string
	AgingTime string
}
type DeviceList struct {
	ListOutputOptions
}

type DeviceCreate struct {
	DeviceType  string `short:"t" required:"true" long:"devicetype" description:"Device type"`
	MACAddress  string `short:"m" long:"macaddress" default:"" description:"MAC Address"`
	IPAddress   string `short:"i" long:"ipaddress" default:"" description:"IP Address"`
	HostAndPort string `short:"H" long:"hostandport" default:"" description:"Host and port"`
}
type DeviceGetVlan struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type DeviceOltInfo struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type GetMtuSizeRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetPmdTxDis struct{
	ListOutputOptions
	Args struct {
		PortType string `positional-arg-name:"PORT_TYPE" required:"yes"`
		Mode string `positional-arg-name:"MODE" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortNo int `positional-arg-name:"PORT_NO" required:"yes"`
	}`positional-args:"yes"`
}
type SetDevicePort struct{
	ListOutputOptions
	Args struct {
		Mode string `positional-arg-name:"MODE" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortNo int `positional-arg-name:"PORT_NO" required:"yes"`
	}`positional-args:"yes"`
}
type GetDevicePortRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortNo int `positional-arg-name:"PORT_NO" required:"yes"`
	}`positional-args:"yes"`
}

type GetPmdTxDisRequest struct {
	ListOutputOptions
	Args struct {
		PortType string `positional-arg-name:"PORT_TYPE" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortNo int `positional-arg-name:"PORT_NO" required:"yes"`
	}`positional-args:"yes"`

}

type PortResetRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortNo int `positional-arg-name:"PORT_NO" required:"yes"`
	}`positional-args:"yes"`
}
type SetMtuSize struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		MtuSize int `positional-arg-name:"PORT_NO" required:"yes"`
	}`positional-args:"yes"`
}
type SetVlan struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Action int `positional-arg-name:"PORT_TYPE" required:"yes"`
		Vid int `positional-arg-name:"MODE" required:"yes"`
		Pbit int `positional-arg-name:"PORT_NO" required:"yes"`
	}`positional-args:"yes"`
}
type SetLutMode struct{
	ListOutputOptions
	Args struct {
		Direction string `positional-arg-name:"DIRECTION" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Mode string `positional-arg-name:"MODE" required:"yes"`
	}`positional-args:"yes"`
}
type GetLutModeRequest struct{
	ListOutputOptions
	Args struct {
		Direction string `positional-arg-name:"DIRECTION" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetAgingMode struct{
	ListOutputOptions
	Args struct {
		Value string `positional-arg-name:"VALUE" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type GetAgingMode struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetAgingTime struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Value int `positional-arg-name:"VALUE" required:"yes"`
	}`positional-args:"yes"`
}

type GetAgingTimeRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type GetDeviceMacInfoRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetSdnTableRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortId string `positional-arg-name:"PortId" required:"yes"`
		Vid int `positional-arg-name:"VID" required:"yes"`
		Pbit int `positional-arg-name:"pbit" required:"yes"`
	}`positional-args:"yes"`
}
type GetSdnTableRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Address string `positional-arg-name:"Address" required:"yes"`
	}`positional-args:"yes"`
}

type SetLength struct{
	ListOutputOptions
	Args struct {
		Operation string `positional-arg-name:"Operation" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Value int `positional-arg-name:"Address" required:"yes"`
	}`positional-args:"yes"`
}
type GetLengthRequest struct{
	ListOutputOptions
	Args struct {
		Operation string `positional-arg-name:"Operation" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetQuietZone struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Value int `positional-arg-name:"VALUE" required:"yes"`
	}`positional-args:"yes"`
}
type GetQuietZone struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetFecMode struct{
	ListOutputOptions
	Args struct {
		Direction string `positional-arg-name:"Direction" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Mode int `positional-arg-name:"Mode" required:"yes"`
	}`positional-args:"yes"`
}
type GetFecMode struct{
	ListOutputOptions
	Args struct {
		Direction string `positional-arg-name:"DIRECTION" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type AddOnuRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}
type DeleteOnu struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}
type AddOnuSla struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
		Tcont int `positional-arg-name:"Tcont" required:"yes"`
		Type int `positional-arg-name:"Type" required:"yes"`
		Si int `positional-arg-name:"Si" required:"yes"`
		Abmin int `positional-arg-name:"Abmin" required:"yes"`
		Absur int `positional-arg-name:"Absur" required:"yes"`
	}`positional-args:"yes"`
}
type ClearOnuSla struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
		Tcont int `positional-arg-name:"Tcont" required:"yes"`
	}`positional-args:"yes"`
}
type GetSlaTable struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetOnuAllocid struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
		AllocId string `positional-arg-name:"AllocId" required:"yes"`
	}`positional-args:"yes"`
}
type DeleteOnuAllocid struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
		AllocId string `positional-arg-name:"AllocId" required:"yes"`
	}`positional-args:"yes"`
}
type SetVssn struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
		Vssn string `positional-arg-name:"Vssn" required:"yes"`
	}`positional-args:"yes"`
}
type GetOnuVssn struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}
type GetOnuDistance struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}
type SetBurstDelimiter struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Length int `positional-arg-name:"Length" required:"yes"`
		Delimiter string `positional-arg-name:"Delimiter" required:"yes"`
	}`positional-args:"yes"`
}
type GetBurstDelimiter struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetBurstPreamble struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Repeat string `positional-arg-name:"Repeat" required:"yes"`

	}`positional-args:"yes"`
}

type GetBurstPreamble struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetBurstVersion struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Version string `positional-arg-name:"Repeat" required:"yes"`
		Index int `positional-arg-name:"Repeat" required:"yes"`
		Pontag string `positional-arg-name:"Repeat" required:"yes"`
	}`positional-args:"yes"`
}

type GetBurstVersion struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetBurstProfile struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}

type GetBurstProfile struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}

type GetRegisterStatus struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}

type GetOnuInfo struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}

type SetDsOmciOnu struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`

}
type GetOmciStatus struct{
	ListOutputOptions
	Args struct {
		Direction string `positional-arg-name:"Direction" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetDsOmciData struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Control int `positional-arg-name:"Control" required:"yes"`
		Data string `positional-arg-name:"Data" required:"yes"`
	}`positional-args:"yes"`
}

type GetUsOmciData struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetTod struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Mode int `positional-arg-name:"Mod" required:"yes"`
		Time int `positional-arg-name:"Time" required:"yes"`
	}`positional-args:"yes"`
}

type GetTodRequest struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetDataMode struct{
	ListOutputOptions
	Args struct {
		Direction string `positional-arg-name:"Direction" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Mode int `positional-arg-name:"Direction" required:"yes"`
	}`positional-args:"yes"`
}

type GetDataMode struct{
	ListOutputOptions
	Args struct {
		Direction string `positional-arg-name:"Direction" required:"yes"`
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetFecDecMode struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Value int `positional-arg-name:"Value" required:"yes"`
	}`positional-args:"yes"`
}

type GetFecDecMode struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetDelimiter struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Value string `positional-arg-name:"Value" required:"yes"`
	}`positional-args:"yes"`
}

type GetDelimiter struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetErrorPermit struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Value int `positional-arg-name:"Value" required:"yes"`
	}`positional-args:"yes"`
}

type GetErrorPermit struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}

type SetPmControl struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"ONU_ID" required:"yes"`
		Mode int `positional-arg-name:"MODE" required:"yes"`
		PowerTime int `positional-arg-name:"POWER_TIME" required:"yes"`
		AwareTime int `positional-arg-name:"AWARE_TIME" required:"yes"`
	}`positional-args:"yes"`
}

type GetPmControl struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}
type GetPmTable struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`

}
type SetSAOn struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}
type SetSAOff struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"OnuId" required:"yes"`
	}`positional-args:"yes"`
}
type CreateDeviceHandler struct{
	ListOutputOptions
	Args struct {
		Ids []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetSliceBw struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Slice int `positional-arg-name:"SLICEe" required:"yes"`
		Bw int `positional-arg-name:"BW" required:"yes"`
	}`positional-args:"yes"`

}
type GetSliceBw struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}
type SetSlaV2 struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		OnuId int `positional-arg-name:"ONUID" required:"yes"`
		Tcont int `positional-arg-name:"TCONT" required:"yes"`
		Slice int `positional-arg-name:"SLICE" required:"yes"`
		CoDba int `positional-arg-name:"CODBA" required:"yes"`
		Type int `positional-arg-name:"TYPE" required:"yes"`
		Rf int `positional-arg-name:"RF" required:"yes"`
		Ra int `positional-arg-name:"RA" required:"yes"`
		Rn int `positional-arg-name:"RN" required:"yes"`
	}`positional-args:"yes"`
}
type GetSlaV2 struct{
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	}`positional-args:"yes"`
}


type DeviceId string

type MetricName string
type GroupName string
type PortNum uint32
type ValueFlag string

type DeviceDelete struct {
	Force bool `long:"force" description:"Delete device forcefully"`
	Args  struct {
		Ids []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceEnable struct {
	Args struct {
		Ids []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceDisable struct {
	Args struct {
		Ids []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceReboot struct {
	Args struct {
		Ids []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceFlowList struct {
	ListOutputOptions
	FlowIdOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceFlowGroupList struct {
	ListOutputOptions
	GroupListOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}
type DevicePortList struct {
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceInspect struct {
	OutputOptionsJson
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePortEnable struct {
	Args struct {
		Id     DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortId PortNum  `positional-arg-name:"PORT_NUMBER" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePortDisable struct {
	Args struct {
		Id     DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortId PortNum  `positional-arg-name:"PORT_NUMBER" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigsGet struct {
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigMetricList struct {
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigGroupList struct {
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigGroupMetricList struct {
	ListOutputOptions
	Args struct {
		Id    DeviceId  `positional-arg-name:"DEVICE_ID" required:"yes"`
		Group GroupName `positional-arg-name:"GROUP_NAME" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigFrequencySet struct {
	OutputOptions
	Args struct {
		Id       DeviceId      `positional-arg-name:"DEVICE_ID" required:"yes"`
		Interval time.Duration `positional-arg-name:"INTERVAL" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigMetricEnable struct {
	Args struct {
		Id      DeviceId     `positional-arg-name:"DEVICE_ID" required:"yes"`
		Metrics []MetricName `positional-arg-name:"METRIC_NAME" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigMetricDisable struct {
	Args struct {
		Id      DeviceId     `positional-arg-name:"DEVICE_ID" required:"yes"`
		Metrics []MetricName `positional-arg-name:"METRIC_NAME" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigGroupEnable struct {
	Args struct {
		Id    DeviceId  `positional-arg-name:"DEVICE_ID" required:"yes"`
		Group GroupName `positional-arg-name:"GROUP_NAME" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigGroupDisable struct {
	Args struct {
		Id    DeviceId  `positional-arg-name:"DEVICE_ID" required:"yes"`
		Group GroupName `positional-arg-name:"GROUP_NAME" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigGroupFrequencySet struct {
	OutputOptions
	Args struct {
		Id       DeviceId      `positional-arg-name:"DEVICE_ID" required:"yes"`
		Group    GroupName     `positional-arg-name:"GROUP_NAME" required:"yes"`
		Interval time.Duration `positional-arg-name:"INTERVAL" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceGetExtValue struct {
	ListOutputOptions
	Args struct {
		Id        DeviceId  `positional-arg-name:"DEVICE_ID" required:"yes"`
		Valueflag ValueFlag `positional-arg-name:"VALUE_FLAG" required:"yes"`
	} `positional-args:"yes"`
}

type DevicePmConfigSetMaxSkew struct {
	Args struct {
		Id      DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		MaxSkew uint32   `positional-arg-name:"MAX_SKEW" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceOnuListImages struct {
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceOnuDownloadImage struct {
	Args struct {
		Id           DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Name         string   `positional-arg-name:"IMAGE_NAME" required:"yes"`
		Url          string   `positional-arg-name:"IMAGE_URL" required:"yes"`
		ImageVersion string   `positional-arg-name:"IMAGE_VERSION" required:"yes"`
		Crc          uint32   `positional-arg-name:"IMAGE_CRC" required:"yes"`
		LocalDir     string   `positional-arg-name:"IMAGE_LOCAL_DIRECTORY"`
	} `positional-args:"yes"`
}

type DeviceOnuActivateImageUpdate struct {
	Args struct {
		Id           DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		Name         string   `positional-arg-name:"IMAGE_NAME" required:"yes"`
		ImageVersion string   `positional-arg-name:"IMAGE_VERSION" required:"yes"`
		SaveConfig   bool     `positional-arg-name:"SAVE_EXISTING_CONFIG"`
		LocalDir     string   `positional-arg-name:"IMAGE_LOCAL_DIRECTORY"`
	} `positional-args:"yes"`
}

type OnuDownloadImage struct {
	ListOutputOptions
	Args struct {
		ImageVersion      string     `positional-arg-name:"IMAGE_VERSION" required:"yes"`
		Url               string     `positional-arg-name:"IMAGE_URL" required:"yes"`
		Vendor            string     `positional-arg-name:"IMAGE_VENDOR"`
		ActivateOnSuccess bool       `positional-arg-name:"IMAGE_ACTIVATE_ON_SUCCESS"`
		CommitOnSuccess   bool       `positional-arg-name:"IMAGE_COMMIT_ON_SUCCESS"`
		Crc               uint32     `positional-arg-name:"IMAGE_CRC"`
		IDs               []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type OnuActivateImage struct {
	ListOutputOptions
	Args struct {
		ImageVersion    string     `positional-arg-name:"IMAGE_VERSION" required:"yes"`
		CommitOnSuccess bool       `positional-arg-name:"IMAGE_COMMIT_ON_SUCCESS"`
		IDs             []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type OnuAbortUpgradeImage struct {
	ListOutputOptions
	Args struct {
		ImageVersion string     `positional-arg-name:"IMAGE_VERSION" required:"yes"`
		IDs          []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type OnuCommitImage struct {
	ListOutputOptions
	Args struct {
		ImageVersion string     `positional-arg-name:"IMAGE_VERSION" required:"yes"`
		IDs          []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type OnuImageStatus struct {
	ListOutputOptions
	Args struct {
		ImageVersion string     `positional-arg-name:"IMAGE_VERSION" required:"yes"`
		IDs          []DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type OnuListImages struct {
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceGetPortStats struct {
	ListOutputOptions
	Args struct {
		Id       DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortNo   uint32   `positional-arg-name:"PORT_NO" required:"yes"`
		PortType string   `positional-arg-name:"PORT_TYPE" required:"yes"`
	} `positional-args:"yes"`
}
type UniStatus struct {
	ListOutputOptions
	Args struct {
		Id       DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		UniIndex uint32   `positional-arg-name:"UNI_INDEX" required:"yes"`
	} `positional-args:"yes"`
}
type OnuPonOpticalInfo struct {
	ListOutputOptions
	Args struct {
		Id DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
	} `positional-args:"yes"`
}

type GetOnuStats struct {
	ListOutputOptions
	Args struct {
		OltId  DeviceId `positional-arg-name:"OLT_DEVICE_ID" required:"yes"`
		IntfId uint32   `positional-arg-name:"PON_INTF_ID" required:"yes"`
		OnuId  uint32   `positional-arg-name:"ONU_ID" required:"yes"`
	} `positional-args:"yes"`
}

type GetOnuEthernetFrameExtendedPmCounters struct {
	ListOutputOptions
	Reset bool `long:"reset" description:"Reset the counters"`
	Args  struct {
		Id       DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		UniIndex *uint32  `positional-arg-name:"UNI_INDEX"`
	} `positional-args:"yes"`
}

type RxPower struct {
	ListOutputOptions
	Args struct {
		Id     DeviceId `positional-arg-name:"DEVICE_ID" required:"yes"`
		PortNo uint32   `positional-arg-name:"PORT_NO" required:"yes"`
		OnuNo  uint32   `positional-arg-name:"ONU_NO" required:"yes"`
	} `positional-args:"yes"`
}

type DeviceOpts struct {
	List    DeviceList          `command:"list"`
	Create  DeviceCreate        `command:"create"`
	Delete  DeviceDelete        `command:"delete"`
	Enable  DeviceEnable        `command:"enable"`
	Disable DeviceDisable       `command:"disable"`
	Flows   DeviceFlowList      `command:"flows"`
	Groups  DeviceFlowGroupList `command:"groups"`
	Port    struct {
		List    DevicePortList    `command:"list"`
		Enable  DevicePortEnable  `command:"enable"`
		Disable DevicePortDisable `command:"disable"`
                Set SetDevicePort `command:"set"`
                Get GetDevicePortRequest `command:"get"`
                Reset PortResetRequest `command:"reset"`
	} `command:"port"`
	Inspect DeviceInspect `command:"inspect"`
	Reboot  DeviceReboot  `command:"reboot"`
	Value   struct {
		Get DeviceGetExtValue `command:"get"`
	} `command:"value"`
	PmConfig struct {
		Get     DevicePmConfigsGet `command:"get"`
		MaxSkew struct {
			Set DevicePmConfigSetMaxSkew `command:"set"`
		} `command:"maxskew"`
		Frequency struct {
			Set DevicePmConfigFrequencySet `command:"set"`
		} `command:"frequency"`
		Metric struct {
			List    DevicePmConfigMetricList    `command:"list"`
			Enable  DevicePmConfigMetricEnable  `command:"enable"`
			Disable DevicePmConfigMetricDisable `command:"disable"`
		} `command:"metric"`
		Group struct {
			List    DevicePmConfigGroupList         `command:"list"`
			Enable  DevicePmConfigGroupEnable       `command:"enable"`
			Disable DevicePmConfigGroupDisable      `command:"disable"`
			Set     DevicePmConfigGroupFrequencySet `command:"set"`
		} `command:"group"`
		GroupMetric struct {
			List DevicePmConfigGroupMetricList `command:"list"`
		} `command:"groupmetric"`
	} `command:"pmconfig"`
	Image struct {
		Get      DeviceOnuListImages          `command:"list"`
		Download DeviceOnuDownloadImage       `command:"download"`
		Activate DeviceOnuActivateImageUpdate `command:"activate"`
	} `command:"image"`
	DownloadImage struct {
		Download     OnuDownloadImage     `command:"download"`
		Activate     OnuActivateImage     `command:"activate"`
		Commit       OnuCommitImage       `command:"commit"`
		AbortUpgrade OnuAbortUpgradeImage `command:"abort"`
		Status       OnuImageStatus       `command:"status"`
		List         OnuListImages        `command:"list" `
	} `command:"onuimage"`
	GetExtVal struct {
		Stats                   DeviceGetPortStats                    `command:"portstats"`
		UniStatus               UniStatus                             `command:"unistatus"`
		OpticalInfo             OnuPonOpticalInfo                     `command:"onu_pon_optical_info"`
		OnuStats                GetOnuStats                           `command:"onu_stats"`
		EthernetFrameExtendedPm GetOnuEthernetFrameExtendedPmCounters `command:"ethernet_frame_extended_pm"`
		RxPower                 RxPower                               `command:"rxpower"`
	} `command:"getextval"`
	BossCommand struct{
		Vlan  DeviceGetVlan `command:"vlan"`
		Olt_connect DevOltConnect `command:"olt_connect"`
		Olt_info  DeviceOltInfo `command:"olt_info"`
		MtuSize GetMtuSizeRequest `command:"mtu"`
		Get_Lut GetLutModeRequest `command:"lookup"`
		Get_Aging_Mode GetAgingMode `command:"aging_mode"`
		Get_Aging_Time GetAgingTimeRequest `command:"aging_time"`
		Get_Mac_Info GetDeviceMacInfoRequest `command:"mac_info"`
		Get_Sdn_Table GetSdnTableRequest `command:"sdn_table"`
		Get_Length GetLengthRequest `command:"length"`
		Get_Quiet_Zone GetQuietZone `command:"quiet_zone"`
		Get_Fec_Mode GetFecMode `command:"fec_mode"`
		BossCommandGetSlaTable struct{
			Get_Sla_Table GetSlaTable `command:"sla_table"`
		}`command:"onu"`
		Get_Burst_Delimiter GetBurstDelimiter `command:"burst_delimeter"`
		Get_Burst_Preamble GetBurstPreamble `command:"burst_preamble"`
		Get_Burst_Version GetBurstVersion `command:"burst_version"`
		Get_Burst_Profile GetBurstProfile `command:"burst_profile"`
		Get_Register_Status GetRegisterStatus `command:"register_status"`
		Get_Onu_Info GetOnuInfo `command:"onu_info"`
		BossCommandGetOmci struct {
			Get_Omci_Status GetOmciStatus `command:"status"`
		}`command:"omci"`
		Get_Tod GetTodRequest `command:"tod"`
		Get_Data_mod GetDataMode `command:"data_mode"`
		Get_FecDec_Mode GetFecDecMode `command:"fec_dec"`
		Get_Delimit GetDelimiter `command:"delimeter"`
		Get_Error_permit GetErrorPermit `command:"error_permit"`
		BossCommandUsOmci struct{
			Get_UsOmci_Data GetUsOmciData `command:"data"`
		}`command:"usomci"`
		Get_Slice_Bw GetSliceBw `command:"slice_bw"`
		Get_Sla2_Table GetSlaV2 `command:"sla2_table"`
	} `command:"get"`
	BossCommandPmd struct{
		Set SetPmdTxDis `command:"set"`
		Get GetPmdTxDisRequest `command:"get"`
		Get_Status GetDevicePmdStatus `command:"get_status"`
	}`command:"pmd"`
	BossCommandSet struct{
		SetMtu SetMtuSize `command:"mtu"`
		Set_vlan SetVlan `command:"vlan"`
		Set_Lut SetLutMode `command:"lookup"`
		Set_Aging SetAgingMode `command:"aging_mode"`
		Set_Aging_Time SetAgingTime `command:"aging_time"`
		Set_Sdn_Table SetSdnTableRequest `command:"sdn_table"`
		Length SetLength `command:"length"`
		Set_Quiet_Zone SetQuietZone `command:"quiet_zone"`
		Set_Fec_Mode SetFecMode `command:"fec_mode"`
		Set_Brust_Delimiter SetBurstDelimiter `command:"burst_delimeter"`
		Set_Bruest_Preamble SetBurstPreamble `command:"burst_preamble"`
		Set_Bruest_Version SetBurstVersion `command:"burst_version"`
		Set_Bruest_Profile SetBurstProfile `command:"burst_profile"`
		BossCommandDownOmci struct{
			Set_DsOmci_Onu SetDsOmciOnu `command:"onu"`
			Set_DsOmci_Data SetDsOmciData `command:"data"`
		}`command:"downomci"`
		Set_Tod SetTod `command:"tod"`
		Set_Data_Mode SetDataMode `command:"data_mode"`
		Set_FecDEC SetFecDecMode `command:"fec_dec"`
		Set_Delimiter SetDelimiter `command:"delimeter"`
		Set_Error_Permit SetErrorPermit `command:"error_permit"`
		Set_Slice_Bw SetSliceBw `command:"slice_bw"`
		Set_Sla2 SetSlaV2 `command:"sla2"`
	}`command:"set"`
	BossCommandAdd struct{
		Add_Onu AddOnuRequest `command:"onu"`
	}`command:"add"`
	BossCommandDelete struct{
		Delete_Onu DeleteOnu `command:"onu"`
	}`command:"deleted"`
	BossCommandOnu struct{
		BossCommandOnuAdd struct{
			Add_Onu_Sla AddOnuSla `command:"sla"`
		}`command:"add"`
		BossCommandOnuClear struct{
			Clear_Onu_Sla ClearOnuSla `command:"sla"`
		}`command:"clear"`
		BossCommandOnuSet struct {
			Set_Onu_Alloc_id SetOnuAllocid `command:"allocid"`
			Set_Vssn SetVssn `command:"vssn"`
			Set_Pm_Control SetPmControl `command:"pm_control"`
		}`command:"set"`
		BossCommandOnuDelete struct {
			Delete_Onu_Alloc_id DeleteOnuAllocid `command:"allocid"`
		}`command:"delete"`
		BossCommandOnuGet struct{
			Get_Onu_Vssn GetOnuVssn `command:"vssn"`
			Get_Onu_Distance GetOnuDistance `command:"distance"`
			Get_Pm_Table GetPmTable `command:"pm_table"`
			Get_Pm_Control GetPmControl `command:"pm_control"`
		}`command:"get"`
		BossCommandOnuEnable struct{
			Set_Sa_On SetSAOn `command:"sa"`
		}`command:"enable"`
		BossCommandOnuDisable struct{
			Set_Sa_Off SetSAOff `command:"sa"`
		}`command:"disable"`

	}`command:"onu"`
	BossCommandDeviceHadler CreateDeviceHandler `command:"createDeviceHandler"`
}

var deviceOpts = DeviceOpts{}

func RegisterDeviceCommands(parser *flags.Parser) {
	if _, err := parser.AddCommand("device", "device commands", "Commands to query and manipulate VOLTHA devices", &deviceOpts); err != nil {
		Error.Fatalf("Unexpected error while attempting to register device commands : %s", err)
	}
}

func (i *MetricName) Complete(match string) []flags.Completion {
	conn, err := NewConnection()
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	var deviceId string
found:
	for i := len(os.Args) - 1; i >= 0; i -= 1 {
		switch os.Args[i] {
		case "enable":
			fallthrough
		case "disable":
			if len(os.Args) > i+1 {
				deviceId = os.Args[i+1]
			} else {
				return nil
			}
			break found
		default:
		}
	}

	if len(deviceId) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(deviceId)}

	pmconfigs, err := client.ListDevicePmConfigs(ctx, &id)

	if err != nil {
		return nil
	}

	list := make([]flags.Completion, 0)
	for _, metrics := range pmconfigs.Metrics {
		if strings.HasPrefix(metrics.Name, match) {
			list = append(list, flags.Completion{Item: metrics.Name})
		}
	}

	return list
}

func (i *GroupName) Complete(match string) []flags.Completion {
	conn, err := NewConnection()
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	var deviceId string
found:
	for i := len(os.Args) - 1; i >= 0; i -= 1 {
		switch os.Args[i] {
		case "list":
			fallthrough
		case "enable":
			fallthrough
		case "disable":
			if len(os.Args) > i+1 {
				deviceId = os.Args[i+1]
			} else {
				return nil
			}
			break found
		default:
		}
	}

	if len(deviceId) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(deviceId)}

	pmconfigs, err := client.ListDevicePmConfigs(ctx, &id)

	if err != nil {
		return nil
	}

	list := make([]flags.Completion, 0)
	for _, group := range pmconfigs.Groups {
		if strings.HasPrefix(group.GroupName, match) {
			list = append(list, flags.Completion{Item: group.GroupName})
		}
	}
	return list
}

func (i *PortNum) Complete(match string) []flags.Completion {
	conn, err := NewConnection()
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	/*
	 * The command line args when completing for PortNum will be a DeviceId
	 * followed by one or more PortNums. So walk the argument list from the
	 * end and find the first argument that is enable/disable as those are
	 * the subcommands that come before the positional arguments. It would
	 * be nice if this package gave us the list of optional arguments
	 * already parsed.
	 */
	var deviceId string
found:
	for i := len(os.Args) - 1; i >= 0; i -= 1 {
		switch os.Args[i] {
		case "enable":
			fallthrough
		case "disable":
			if len(os.Args) > i+1 {
				deviceId = os.Args[i+1]
			} else {
				return nil
			}
			break found
		default:
		}
	}

	if len(deviceId) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(deviceId)}

	ports, err := client.ListDevicePorts(ctx, &id)
	if err != nil {
		return nil
	}

	list := make([]flags.Completion, 0)
	for _, item := range ports.Items {
		pn := strconv.FormatUint(uint64(item.PortNo), 10)
		if strings.HasPrefix(pn, match) {
			list = append(list, flags.Completion{Item: pn})
		}
	}

	return list
}

func (i *DeviceId) Complete(match string) []flags.Completion {
	conn, err := NewConnection()
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	devices, err := client.ListDevices(ctx, &empty.Empty{})
	if err != nil {
		return nil
	}

	list := make([]flags.Completion, 0)
	for _, item := range devices.Items {
		if strings.HasPrefix(item.Id, match) {
			list = append(list, flags.Completion{Item: item.Id})
		}
	}

	return list
}

func (options *DeviceList) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	devices, err := client.ListDevices(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-list", "format", DEFAULT_DEVICE_FORMAT)
	}
	if options.Quiet {
		outputFormat = "{{.Id}}"
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-list", "order", "")
	}

	// Make sure json output prints an empty list, not "null"
	if devices.Items == nil {
		devices.Items = make([]*voltha.Device, 0)
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		Filter:    options.Filter,
		OrderBy:   orderBy,
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      devices.Items,
	}

	GenerateOutput(&result)
	return nil
}

func (options *DeviceCreate) Execute(args []string) error {

	device := voltha.Device{}
	if options.HostAndPort != "" {
		device.Address = &voltha.Device_HostAndPort{HostAndPort: options.HostAndPort}
	} else if options.IPAddress != "" {
		device.Address = &voltha.Device_Ipv4Address{Ipv4Address: options.IPAddress}
	}
	if options.MACAddress != "" {
		device.MacAddress = strings.ToLower(options.MACAddress)
	}
	if options.DeviceType != "" {
		device.Type = options.DeviceType
	}

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	createdDevice, err := client.CreateDevice(ctx, &device)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", createdDevice.Id)

	return nil
}

func (options *DeviceDelete) Execute(args []string) error {

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
		if options.Force {
			_, err = client.ForceDeleteDevice(ctx, &id)
		} else {
			_, err = client.DeleteDevice(ctx, &id)
		}

		if err != nil {
			Error.Printf("Error while deleting '%s': %s\n", i, err)
			lastErr = err
			continue
		}
		fmt.Printf("%s\n", i)
	}

	if lastErr != nil {
		return NoReportErr
	}
	return nil
}

func (options *DeviceEnable) Execute(args []string) error {
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

		_, err := client.EnableDevice(ctx, &id)
		if err != nil {
			Error.Printf("Error while enabling '%s': %s\n", i, err)
			lastErr = err
			continue
		}
		fmt.Printf("%s\n", i)
	}

	if lastErr != nil {
		return NoReportErr
	}
	return nil
}

func (options *DeviceDisable) Execute(args []string) error {
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

		_, err := client.DisableDevice(ctx, &id)
		if err != nil {
			Error.Printf("Error while disabling '%s': %s\n", i, err)
			lastErr = err
			continue
		}
		fmt.Printf("%s\n", i)
	}

	if lastErr != nil {
		return NoReportErr
	}
	return nil
}

func (options *DeviceReboot) Execute(args []string) error {
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

		_, err := client.RebootDevice(ctx, &id)
		if err != nil {
			Error.Printf("Error while rebooting '%s': %s\n", i, err)
			lastErr = err
			continue
		}
		fmt.Printf("%s\n", i)
	}

	if lastErr != nil {
		return NoReportErr
	}
	return nil
}

func (options *DevicePortList) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	ports, err := client.ListDevicePorts(ctx, &id)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_PORTS_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		Filter:    options.Filter,
		OrderBy:   orderBy,
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      ports.Items,
	}

	GenerateOutput(&result)
	return nil
}

func (options *DeviceFlowList) Execute(args []string) error {
	fl := &FlowList{}
	fl.ListOutputOptions = options.ListOutputOptions
	fl.FlowIdOptions = options.FlowIdOptions
	fl.Args.Id = string(options.Args.Id)
	fl.Method = "device-flows"
	return fl.Execute(args)
}

func (options *DeviceFlowGroupList) Execute(args []string) error {
	grp := &GroupList{}
	grp.ListOutputOptions = options.ListOutputOptions
	grp.GroupListOptions = options.GroupListOptions
	grp.Args.Id = string(options.Args.Id)
	grp.Method = "device-groups"
	return grp.Execute(args)
}

func (options *DeviceInspect) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("only a single argument 'DEVICE_ID' can be provided")
	}

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	device, err := client.GetDevice(ctx, &id)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-inspect", "format", DEFAULT_DEVICE_INSPECT_FORMAT)
	}
	if options.Quiet {
		outputFormat = "{{.Id}}"
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      device,
	}
	GenerateOutput(&result)
	return nil
}

/*Device  Port Enable */
func (options *DevicePortEnable) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	port := voltha.Port{DeviceId: string(options.Args.Id), PortNo: uint32(options.Args.PortId)}

	_, err = client.EnablePort(ctx, &port)
	if err != nil {
		Error.Printf("Error enabling port number %v on device Id %s,err=%s\n", options.Args.PortId, options.Args.Id, ErrorToString(err))
		return err
	}

	return nil
}

/*Device  Port Disable */
func (options *DevicePortDisable) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	port := voltha.Port{DeviceId: string(options.Args.Id), PortNo: uint32(options.Args.PortId)}

	_, err = client.DisablePort(ctx, &port)
	if err != nil {
		Error.Printf("Error enabling port number %v on device Id %s,err=%s\n", options.Args.PortId, options.Args.Id, ErrorToString(err))
		return err
	}

	return nil
}

func (options *DevicePmConfigSetMaxSkew) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	pmConfigs.MaxSkew = options.Args.MaxSkew

	_, err = client.UpdateDevicePmConfigs(ctx, pmConfigs)
	if err != nil {
		return err
	}

	return nil
}

func (options *DevicePmConfigsGet) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-pm-configs", "format", DEFAULT_DEVICE_PM_CONFIG_GET_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-pm-configs", "order", "")
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		Filter:    options.Filter,
		OrderBy:   orderBy,
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      pmConfigs,
	}

	GenerateOutput(&result)
	return nil

}

func (options *DevicePmConfigMetricList) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	if !pmConfigs.Grouped {
		for _, metric := range pmConfigs.Metrics {
			if metric.SampleFreq == 0 {
				metric.SampleFreq = pmConfigs.DefaultFreq
			}
		}
		outputFormat := CharReplacer.Replace(options.Format)
		if outputFormat == "" {
			outputFormat = GetCommandOptionWithDefault("device-pm-configs", "format", DEFAULT_DEVICE_PM_CONFIG_METRIC_LIST_FORMAT)
		}

		orderBy := options.OrderBy
		if orderBy == "" {
			orderBy = GetCommandOptionWithDefault("device-pm-configs", "order", "")
		}

		result := CommandResult{
			Format:    format.Format(outputFormat),
			Filter:    options.Filter,
			OrderBy:   orderBy,
			OutputAs:  toOutputType(options.OutputAs),
			NameLimit: options.NameLimit,
			Data:      pmConfigs.Metrics,
		}

		GenerateOutput(&result)
		return nil
	} else {
		return fmt.Errorf("Device '%s' does not have Non Grouped Metrics", options.Args.Id)
	}
}

func (options *DevicePmConfigMetricEnable) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	if !pmConfigs.Grouped {
		metrics := make(map[string]struct{})
		for _, metric := range pmConfigs.Metrics {
			metrics[metric.Name] = struct{}{}
		}

		for _, metric := range pmConfigs.Metrics {
			for _, mName := range options.Args.Metrics {
				if _, exist := metrics[string(mName)]; !exist {
					return fmt.Errorf("Metric Name '%s' does not exist", mName)
				}

				if string(mName) == metric.Name && !metric.Enabled {
					metric.Enabled = true
					_, err := client.UpdateDevicePmConfigs(ctx, pmConfigs)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		return fmt.Errorf("Device '%s' does not have Non Grouped Metrics", options.Args.Id)
	}
	return nil
}

func (options *DevicePmConfigMetricDisable) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	if !pmConfigs.Grouped {
		metrics := make(map[string]struct{})
		for _, metric := range pmConfigs.Metrics {
			metrics[metric.Name] = struct{}{}
		}

		for _, metric := range pmConfigs.Metrics {
			for _, mName := range options.Args.Metrics {
				if _, have := metrics[string(mName)]; !have {
					return fmt.Errorf("Metric Name '%s' does not exist", mName)
				}
				if string(mName) == metric.Name && metric.Enabled {
					metric.Enabled = false
					_, err := client.UpdateDevicePmConfigs(ctx, pmConfigs)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("Metric '%s' cannot be disabled", string(mName))
				}
			}
		}
	} else {
		return fmt.Errorf("Device '%s' does not have Non Grouped Metrics", options.Args.Id)
	}
	return nil
}

func (options *DevicePmConfigGroupEnable) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	if pmConfigs.Grouped {
		groups := make(map[string]struct{})
		for _, group := range pmConfigs.Groups {
			groups[group.GroupName] = struct{}{}
		}
		for _, group := range pmConfigs.Groups {
			if _, have := groups[string(options.Args.Group)]; !have {
				return fmt.Errorf("Group Name '%s' does not exist", options.Args.Group)
			}
			if string(options.Args.Group) == group.GroupName && !group.Enabled {
				group.Enabled = true
				_, err := client.UpdateDevicePmConfigs(ctx, pmConfigs)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("Device '%s' does not have Group Metrics", options.Args.Id)
	}
	return nil
}

func (options *DevicePmConfigGroupDisable) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	if pmConfigs.Grouped {
		groups := make(map[string]struct{})
		for _, group := range pmConfigs.Groups {
			groups[group.GroupName] = struct{}{}
		}

		for _, group := range pmConfigs.Groups {
			if _, have := groups[string(options.Args.Group)]; !have {
				return fmt.Errorf("Group Name '%s' does not exist", options.Args.Group)
			}

			if string(options.Args.Group) == group.GroupName && group.Enabled {
				group.Enabled = false
				_, err := client.UpdateDevicePmConfigs(ctx, pmConfigs)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("Device '%s' does not have Group Metrics", options.Args.Id)
	}
	return nil
}

func (options *DevicePmConfigGroupFrequencySet) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	if pmConfigs.Grouped {
		groups := make(map[string]struct{})
		for _, group := range pmConfigs.Groups {
			groups[group.GroupName] = struct{}{}
		}

		for _, group := range pmConfigs.Groups {
			if _, have := groups[string(options.Args.Group)]; !have {
				return fmt.Errorf("group name '%s' does not exist", options.Args.Group)
			}

			if string(options.Args.Group) == group.GroupName {
				if !group.Enabled {
					return fmt.Errorf("group '%s' is not enabled", options.Args.Group)
				}
				group.GroupFreq = uint32(options.Args.Interval.Seconds())
				_, err = client.UpdateDevicePmConfigs(ctx, pmConfigs)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("device '%s' does not have group metrics", options.Args.Id)
	}
	return nil
}

func (options *DevicePmConfigGroupList) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	if pmConfigs.Grouped {
		for _, group := range pmConfigs.Groups {
			if group.GroupFreq == 0 {
				group.GroupFreq = pmConfigs.DefaultFreq
			}
		}
		outputFormat := CharReplacer.Replace(options.Format)
		if outputFormat == "" {
			outputFormat = GetCommandOptionWithDefault("device-pm-configs", "format", DEFAULT_DEVICE_PM_CONFIG_GROUP_LIST_FORMAT)
		}

		orderBy := options.OrderBy
		if orderBy == "" {
			orderBy = GetCommandOptionWithDefault("device-pm-configs", "order", "")
		}

		result := CommandResult{
			Format:    format.Format(outputFormat),
			Filter:    options.Filter,
			OrderBy:   orderBy,
			OutputAs:  toOutputType(options.OutputAs),
			NameLimit: options.NameLimit,
			Data:      pmConfigs.Groups,
		}

		GenerateOutput(&result)
	} else {
		return fmt.Errorf("Device '%s' does not have Group Metrics", string(options.Args.Id))
	}
	return nil
}

func (options *DevicePmConfigGroupMetricList) Execute(args []string) error {

	var metrics []*voltha.PmConfig
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	for _, groups := range pmConfigs.Groups {

		if string(options.Args.Group) == groups.GroupName {
			for _, metric := range groups.Metrics {
				if metric.SampleFreq == 0 && groups.GroupFreq == 0 {
					metric.SampleFreq = pmConfigs.DefaultFreq
				} else {
					metric.SampleFreq = groups.GroupFreq
				}
			}
			metrics = groups.Metrics
		}
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-pm-configs", "format", DEFAULT_DEVICE_PM_CONFIG_METRIC_LIST_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-pm-configs", "order", "")
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		Filter:    options.Filter,
		OrderBy:   orderBy,
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      metrics,
	}

	GenerateOutput(&result)
	return nil

}

func (options *DevicePmConfigFrequencySet) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := voltha.ID{Id: string(options.Args.Id)}

	pmConfigs, err := client.ListDevicePmConfigs(ctx, &id)
	if err != nil {
		return err
	}

	pmConfigs.DefaultFreq = uint32(options.Args.Interval.Seconds())

	_, err = client.UpdateDevicePmConfigs(ctx, pmConfigs)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-pm-configs", "format", DEFAULT_DEVICE_PM_CONFIG_GET_FORMAT)
	}
	if options.Quiet {
		outputFormat = "{{.Id}}"
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      pmConfigs,
	}

	GenerateOutput(&result)
	return nil

}

func (options *OnuDownloadImage) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	var devIDList []*common.ID
	for _, i := range options.Args.IDs {

		devIDList = append(devIDList, &common.ID{Id: string(i)})
	}

	downloadImage := voltha.DeviceImageDownloadRequest{
		DeviceId: devIDList,
		Image: &voltha.Image{
			Url:     options.Args.Url,
			Crc32:   options.Args.Crc,
			Vendor:  options.Args.Vendor,
			Version: options.Args.ImageVersion,
		},
		ActivateOnSuccess: options.Args.ActivateOnSuccess,
		CommitOnSuccess:   options.Args.CommitOnSuccess,
	}

	deviceImageResp, err := client.DownloadImageToDevice(ctx, &downloadImage)
	if err != nil {
		return err
	}

	outputFormat := GetCommandOptionWithDefault("onu-image-download", "format", ONU_IMAGE_STATUS_FORMAT)
	// Make sure json output prints an empty list, not "null"
	if deviceImageResp.DeviceImageStates == nil {
		deviceImageResp.DeviceImageStates = make([]*voltha.DeviceImageState, 0)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      deviceImageResp.DeviceImageStates,
	}
	GenerateOutput(&result)
	return nil

}

func (options *OnuActivateImage) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	var devIDList []*common.ID
	for _, i := range options.Args.IDs {

		devIDList = append(devIDList, &common.ID{Id: string(i)})
	}

	downloadImage := voltha.DeviceImageRequest{
		DeviceId:        devIDList,
		Version:         options.Args.ImageVersion,
		CommitOnSuccess: options.Args.CommitOnSuccess,
	}

	deviceImageResp, err := client.ActivateImage(ctx, &downloadImage)
	if err != nil {
		return err
	}

	outputFormat := GetCommandOptionWithDefault("onu-image-activate", "format", ONU_IMAGE_STATUS_FORMAT)
	// Make sure json output prints an empty list, not "null"
	if deviceImageResp.DeviceImageStates == nil {
		deviceImageResp.DeviceImageStates = make([]*voltha.DeviceImageState, 0)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      deviceImageResp.DeviceImageStates,
	}
	GenerateOutput(&result)

	return nil

}

func (options *OnuAbortUpgradeImage) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	var devIDList []*common.ID
	for _, i := range options.Args.IDs {

		devIDList = append(devIDList, &common.ID{Id: string(i)})
	}

	downloadImage := voltha.DeviceImageRequest{
		DeviceId: devIDList,
		Version:  options.Args.ImageVersion,
	}

	deviceImageResp, err := client.AbortImageUpgradeToDevice(ctx, &downloadImage)
	if err != nil {
		return err
	}

	outputFormat := GetCommandOptionWithDefault("onu-image-abort", "format", ONU_IMAGE_STATUS_FORMAT)
	// Make sure json output prints an empty list, not "null"
	if deviceImageResp.DeviceImageStates == nil {
		deviceImageResp.DeviceImageStates = make([]*voltha.DeviceImageState, 0)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      deviceImageResp.DeviceImageStates,
	}
	GenerateOutput(&result)

	return nil

}

func (options *OnuCommitImage) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	var devIDList []*common.ID
	for _, i := range options.Args.IDs {

		devIDList = append(devIDList, &common.ID{Id: string(i)})
	}
	downloadImage := voltha.DeviceImageRequest{
		DeviceId: devIDList,
		Version:  options.Args.ImageVersion,
	}

	deviceImageResp, err := client.CommitImage(ctx, &downloadImage)
	if err != nil {
		return err
	}

	outputFormat := GetCommandOptionWithDefault("onu-image-commit", "format", ONU_IMAGE_STATUS_FORMAT)
	// Make sure json output prints an empty list, not "null"
	if deviceImageResp.DeviceImageStates == nil {
		deviceImageResp.DeviceImageStates = make([]*voltha.DeviceImageState, 0)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      deviceImageResp.DeviceImageStates,
	}
	GenerateOutput(&result)

	return nil

}

func (options *OnuListImages) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := common.ID{Id: string(options.Args.Id)}

	onuImages, err := client.GetOnuImages(ctx, &id)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("onu-image-list", "format", ONU_IMAGE_LIST_FORMAT)
	}

	if options.Quiet {
		outputFormat = "{{.Id}}"
	}

	//TODO orderby

	// Make sure json output prints an empty list, not "null"
	if onuImages.Items == nil {
		onuImages.Items = make([]*voltha.OnuImage, 0)
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      onuImages.Items,
	}

	GenerateOutput(&result)
	return nil

}

func (options *OnuImageStatus) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	var devIDList []*common.ID
	for _, i := range options.Args.IDs {

		devIDList = append(devIDList, &common.ID{Id: string(i)})
	}

	imageStatusReq := voltha.DeviceImageRequest{
		DeviceId: devIDList,
		Version:  options.Args.ImageVersion,
	}
	imageStatus, err := client.GetImageStatus(ctx, &imageStatusReq)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-image-list", "format", ONU_IMAGE_STATUS_FORMAT)
	}

	if options.Quiet {
		outputFormat = "{{.Id}}"
	}

	//TODO orderby

	// Make sure json output prints an empty list, not "null"
	if imageStatus.DeviceImageStates == nil {
		imageStatus.DeviceImageStates = make([]*voltha.DeviceImageState, 0)
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      imageStatus.DeviceImageStates,
	}

	GenerateOutput(&result)
	return nil

}

func (options *DeviceOnuListImages) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := common.ID{Id: string(options.Args.Id)}

	imageDownloads, err := client.ListImageDownloads(ctx, &id)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-image-list", "format", DEFAULT_DEVICE_IMAGE_LIST_GET_FORMAT)
	}

	if options.Quiet {
		outputFormat = "{{.Id}}"
	}

	//TODO orderby

	// Make sure json output prints an empty list, not "null"
	if imageDownloads.Items == nil {
		imageDownloads.Items = make([]*voltha.ImageDownload, 0)
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      imageDownloads.Items,
	}

	GenerateOutput(&result)
	return nil

}

func (options *DeviceOnuDownloadImage) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	downloadImage := voltha.ImageDownload{
		Id:       string(options.Args.Id),
		Name:     options.Args.Name,
		Url:      options.Args.Url,
		Crc:      options.Args.Crc,
		LocalDir: options.Args.LocalDir,
	}

	_, err = client.DownloadImage(ctx, &downloadImage)
	if err != nil {
		return err
	}

	return nil

}

func (options *DeviceOnuActivateImageUpdate) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	downloadImage := voltha.ImageDownload{
		Id:           string(options.Args.Id),
		Name:         options.Args.Name,
		ImageVersion: options.Args.ImageVersion,
		SaveConfig:   options.Args.SaveConfig,
		LocalDir:     options.Args.LocalDir,
	}

	_, err = client.ActivateImageUpdate(ctx, &downloadImage)
	if err != nil {
		return err
	}

	return nil

}

type ReturnValueRow struct {
	Name   string      `json:"name"`
	Result interface{} `json:"result"`
}

func (options *DeviceGetPortStats) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := extension.NewExtensionClient(conn)
	var portType extension.GetOltPortCounters_PortType

	if options.Args.PortType == "pon" {
		portType = extension.GetOltPortCounters_Port_PON_OLT
	} else if options.Args.PortType == "nni" {

		portType = extension.GetOltPortCounters_Port_ETHERNET_NNI
	} else {
		return fmt.Errorf("expected interface type pon/nni, provided %s", options.Args.PortType)
	}

	singleGetValReq := extension.SingleGetValueRequest{
		TargetId: string(options.Args.Id),
		Request: &extension.GetValueRequest{
			Request: &extension.GetValueRequest_OltPortInfo{
				OltPortInfo: &extension.GetOltPortCounters{
					PortNo:   options.Args.PortNo,
					PortType: portType,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	rv, err := client.GetExtValue(ctx, &singleGetValReq)
	if err != nil {
		Error.Printf("Error getting value on device Id %s,err=%s\n", options.Args.Id, ErrorToString(err))
		return err
	}

	if rv.Response.Status != extension.GetValueResponse_OK {
		return fmt.Errorf("failed to get port stats %v", rv.Response.ErrReason.String())
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-get-port-status", "format", DEFAULT_DEVICE_GET_PORT_STATUS_FORMAT)
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      rv.GetResponse().GetPortCoutners(),
	}
	GenerateOutput(&result)
	return nil
}

func (options *GetOnuStats) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := extension.NewExtensionClient(conn)

	singleGetValReq := extension.SingleGetValueRequest{
		TargetId: string(options.Args.OltId),
		Request: &extension.GetValueRequest{
			Request: &extension.GetValueRequest_OnuPonInfo{
				OnuPonInfo: &extension.GetOnuCountersRequest{
					IntfId: options.Args.IntfId,
					OnuId:  options.Args.OnuId,
				},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	rv, err := client.GetExtValue(ctx, &singleGetValReq)
	if err != nil {
		Error.Printf("Error getting value on device Id %s,err=%s\n", options.Args.OltId, ErrorToString(err))
		return err
	}

	if rv.Response.Status != extension.GetValueResponse_OK {
		return fmt.Errorf("failed to get onu stats %v", rv.Response.ErrReason.String())
	}
	outputFormat := CharReplacer.Replace(options.Format)
	data, formatStr := buildOnuStatsOutputFormat(rv.GetResponse().GetOnuPonCounters())
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-get-onu-status", "format", formatStr)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      data,
	}
	GenerateOutput(&result)
	return nil
}

func (options *GetOnuEthernetFrameExtendedPmCounters) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := extension.NewExtensionClient(conn)
	var singleGetValReq extension.SingleGetValueRequest

	if options.Args.UniIndex != nil {
		singleGetValReq = extension.SingleGetValueRequest{
			TargetId: string(options.Args.Id),
			Request: &extension.GetValueRequest{
				Request: &extension.GetValueRequest_OnuInfo{
					OnuInfo: &extension.GetOmciEthernetFrameExtendedPmRequest{
						OnuDeviceId: string(options.Args.Id),
						Reset_:      options.Reset,
						IsUniIndex: &extension.GetOmciEthernetFrameExtendedPmRequest_UniIndex{
							UniIndex: *options.Args.UniIndex,
						},
					},
				},
			},
		}
	} else {
		singleGetValReq = extension.SingleGetValueRequest{
			TargetId: string(options.Args.Id),
			Request: &extension.GetValueRequest{
				Request: &extension.GetValueRequest_OnuInfo{
					OnuInfo: &extension.GetOmciEthernetFrameExtendedPmRequest{
						OnuDeviceId: string(options.Args.Id),
						Reset_:      options.Reset,
					},
				},
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	rv, err := client.GetExtValue(ctx, &singleGetValReq)
	if err != nil {
		Error.Printf("Error getting value on device Id %s,err=%s\n", options.Args.Id, ErrorToString(err))
		return err
	}

	if rv.Response.Status != extension.GetValueResponse_OK {
		return fmt.Errorf("failed to get ethernet frame extended pm counters %v", rv.Response.ErrReason.String())
	}
	outputFormat := CharReplacer.Replace(options.Format)
	data := buildOnuEthernetFrameExtendedPmOutputFormat(rv.GetResponse().GetOnuCounters())
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-get-onu-status", "format", DEFAULT_ETHERNET_FRAME_EXTENDED_PM_COUNTERS_FORMAT)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      data,
	}
	GenerateOutput(&result)
	return nil
}

func (options *UniStatus) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := extension.NewExtensionClient(conn)

	singleGetValReq := extension.SingleGetValueRequest{
		TargetId: string(options.Args.Id),
		Request: &extension.GetValueRequest{
			Request: &extension.GetValueRequest_UniInfo{
				UniInfo: &extension.GetOnuUniInfoRequest{
					UniIndex: options.Args.UniIndex,
				},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	rv, err := client.GetExtValue(ctx, &singleGetValReq)
	if err != nil {
		Error.Printf("Error getting value on device Id %s,err=%s\n", options.Args.Id, ErrorToString(err))
		return err
	}
	if rv.Response.Status != extension.GetValueResponse_OK {
		return fmt.Errorf("failed to get uni status %v", rv.Response.ErrReason.String())
	}
	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-get-uni-status", "format", DEFAULT_DEVICE_GET_UNI_STATUS_FORMAT)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      rv.GetResponse().GetUniInfo(),
	}
	GenerateOutput(&result)
	return nil
}

func (options *OnuPonOpticalInfo) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := extension.NewExtensionClient(conn)

	singleGetValReq := extension.SingleGetValueRequest{
		TargetId: string(options.Args.Id),
		Request: &extension.GetValueRequest{
			Request: &extension.GetValueRequest_OnuOpticalInfo{
				OnuOpticalInfo: &extension.GetOnuPonOpticalInfo{},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	rv, err := client.GetExtValue(ctx, &singleGetValReq)
	if err != nil {
		Error.Printf("Error getting value on device Id %s,err=%s\n", options.Args.Id, ErrorToString(err))
		return err
	}
	if rv.Response.Status != extension.GetValueResponse_OK {
		return fmt.Errorf("failed to get onu pon optical info %v", rv.Response.ErrReason.String())
	}
	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-get-onu-pon-optical-info", "format", DEFAULT_ONU_PON_OPTICAL_INFO_STATUS_FORMAT)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      rv.GetResponse().GetOnuOpticalInfo(),
	}
	GenerateOutput(&result)
	return nil
}

func (options *RxPower) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := extension.NewExtensionClient(conn)

	singleGetValReq := extension.SingleGetValueRequest{
		TargetId: string(options.Args.Id),
		Request: &extension.GetValueRequest{
			Request: &extension.GetValueRequest_RxPower{
				RxPower: &extension.GetRxPowerRequest{
					IntfId: options.Args.PortNo,
					OnuId:  options.Args.OnuNo,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	rv, err := client.GetExtValue(ctx, &singleGetValReq)
	if err != nil {
		Error.Printf("Error getting value on device Id %s,err=%s\n", options.Args.Id, ErrorToString(err))
		return err
	}
	if rv.Response.Status != extension.GetValueResponse_OK {
		return fmt.Errorf("failed to get rx power %v", rv.Response.ErrReason.String())
	}
	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-get-rx-power", "format", DEFAULT_RX_POWER_STATUS_FORMAT)
	}
	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      rv.GetResponse().GetRxPower(),
	}
	GenerateOutput(&result)
	return nil
}

/*Device  get Onu Distance */
func (options *DeviceGetExtValue) Execute(args []string) error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := voltha.NewVolthaServiceClient(conn)

	valueflag, okay := extension.ValueType_Type_value[string(options.Args.Valueflag)]
	if !okay {
		Error.Printf("Unknown valueflag %s\n", options.Args.Valueflag)
	}

	val := extension.ValueSpecifier{Id: string(options.Args.Id), Value: extension.ValueType_Type(valueflag)}

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	rv, err := client.GetExtValue(ctx, &val)
	if err != nil {
		Error.Printf("Error getting value on device Id %s,err=%s\n", options.Args.Id, ErrorToString(err))
		return err
	}

	var rows []ReturnValueRow
	for name, num := range extension.ValueType_Type_value {
		if num == 0 {
			// EMPTY is not a real value
			continue
		}
		if (rv.Error & uint32(num)) != 0 {
			row := ReturnValueRow{Name: name, Result: "Error"}
			rows = append(rows, row)
		}
		if (rv.Unsupported & uint32(num)) != 0 {
			row := ReturnValueRow{Name: name, Result: "Unsupported"}
			rows = append(rows, row)
		}
		if (rv.Set & uint32(num)) != 0 {
			switch name {
			case "DISTANCE":
				row := ReturnValueRow{Name: name, Result: rv.Distance}
				rows = append(rows, row)
			default:
				row := ReturnValueRow{Name: name, Result: "Unimplemented-in-bossctl"}
				rows = append(rows, row)
			}
		}
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-value-get", "format", DEFAULT_DEVICE_VALUE_GET_FORMAT)
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      rows,
	}
	GenerateOutput(&result)
	return nil
}
func(options *DevOltConnect) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetOltConnect(ctx, &id)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GET_CONNECTION_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	tmp := GetOltConnect{
		DeviceId : Bossresp.DeviceId,
		IP : Bossresp.Ip,
		MAC: Bossresp.Mac,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *DeviceGetVlan) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetVlan(ctx, &id)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_VLAN_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var mode string
	if Bossresp.VlanMode ==0{
		mode="Pass"
	}else{
		mode="Add/Remove"
	}
	
	tmp := GetDeviceVlan{
		DeviceId : Bossresp.DeviceId,
		STAG_mode : mode,
		STAG_fields : Bossresp.Fields,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *DeviceOltInfo) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetOltDeviceInfo(ctx, &id)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GET_DEVINFO_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	tmp := GetDevInfo{
		DeviceId : Bossresp.DeviceId,
		Fpga_Type : Bossresp.FpgaType,
		Fpga_ver : Bossresp.FpgaVer,
		Fpga_date : Bossresp.Fpga_Date,
		Sw_ver : Bossresp.SwVer,
		Sw_date : Bossresp.SwDate,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *SetPmdTxDis) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var mode int32
	var portType int32
	if options.Args.PortType == "nni"{
		portType = 0
	}else if options.Args.PortType == "pon"{
		portType = 1
	} else {
		fmt.Println("port Type Error")
		fmt.Println(options.Args.PortType)
		return nil
	}

	if options.Args.Mode == "enable"{
		mode = 0
	}else if options.Args.Mode =="disable"{
		mode = 1
	}else {
		fmt.Println("mode Error")
		return nil
	}
	param := bossopenolt.SetPmdTxdis{ PortType:int32(portType), Mode:int32(mode), PortNo:int32(options.Args.PortNo)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetpmdtxdisParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetPmdTxDis(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetPmdTxDisRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var portType int32
	if options.Args.PortType == "nni"{
		portType = 0
	}else if options.Args.PortType == "pon"{
		portType =1
	}else {
		fmt.Println("port Type Error " )
		fmt.Println(options.Args.PortType)
		return nil
	}
	param := bossopenolt.GetPmdsKind{ PortType:int32(portType),  PortNo:int32(options.Args.PortNo)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetpmdskindParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetPmdTxdis(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GET_PMDTXDIS_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var port int = int(Bossresp.PortNo)

	tmp := GetPmdTxDis{
		PortNo : strconv.Itoa(port),
		Status : Bossresp.Status,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetDevicePmdStatus) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var portType int32
	if options.Args.PortType == "nni"{
		portType = 0
	}else if options.Args.PortType == "pon"{
		portType =1
	}else {
		fmt.Println("port Type Error " )
		fmt.Println(options.Args.PortType)
		return nil
	}
	param := bossopenolt.GetPmdsKind{ PortType:int32(portType),  PortNo:int32(options.Args.PortNo)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetpmdskindParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetDevicePmdStatus(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GET_PMDDEVICESTATUS_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var port int = int(Bossresp.PortNo)

	tmp := GetPmdDeviceStatus{
		PortNo : strconv.Itoa(port),
		Loss : Bossresp.Loss,
		Module : Bossresp.Module,
		Fault : Bossresp.Fault,
		Link : Bossresp.Link,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetDevicePort) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var mode int32
	if options.Args.Mode=="enable"{
		mode=1
	}else if options.Args.Mode =="disable"{
		mode=0
	}else{
		fmt.Println("mode Error .. input Data : "+ options.Args.Mode )
	}
	param := bossopenolt.SetPort{Mode: mode, PortNo:int32(options.Args.PortNo)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetportAram{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetDevicePort(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetDevicePortRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	param := bossopenolt.SetPortKind{PortNo:int32(options.Args.PortNo)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetportkindParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetDevicePort(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETDEVICEPORT_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	tmp := GetDevicePort{
		PortNo : strconv.Itoa(int(Bossresp.PortNo)),
		State : Bossresp.State,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *PortResetRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	param := bossopenolt.SetPortKind{ PortNo:int32(options.Args.PortNo)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetportkindParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.PortReset(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *SetMtuSize) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	param := bossopenolt.SetMtuSize{ MtuSize:int32(options.Args.MtuSize )}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetmtusizeParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetMtuSize(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetMtuSizeRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	id := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetMtuSize(ctx, &id)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETMTUSIZE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	tmp := GetMtuSize{
		Mtu : strconv.Itoa(int(Bossresp.Mtu)),
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *SetVlan) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	param := bossopenolt.SetVlan{ Action:int32(options.Args.Action), Vid:int32(options.Args.Vid), Pbit:int32(options.Args.Pbit)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetvlanParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetVlan(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *SetLutMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var direct int32
	var mode int32
	if options.Args.Direction=="downstream"{
		direct = 0
	}else if options.Args.Direction=="upstream"{
		direct=1
	}else{
		fmt.Println("select downstream/upstream. input data :" + options.Args.Direction)
		return nil
	}
	if options.Args.Mode=="normal"{
		mode = 0
	}else if options.Args.Mode =="bypass"{
		mode=1
	}else{
		fmt.Println("select normal/bypass. input data :" + options.Args.Mode)
		return nil 
	}
	param := bossopenolt.SetDirectionMode{Direction:direct,Mode: mode}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetdirectiommodeParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetLutMode(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetLutModeRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var direct int32

	if options.Args.Direction=="downstream"{
		direct = 0
	}else if options.Args.Direction=="upstream"{
		direct=1
	}else{
		fmt.Println("select downstream/upstream. input data :" + options.Args.Direction)
		return nil
	}
	param := bossopenolt.GetDirectionValue{Direction:direct}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetdirectionvalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetLutMode(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETMODE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Mode ==0{
		response="normal"
	}else{
		response="bypass"
	}
	tmp := GetMode{
		DeviceId: Bossresp.DeviceId,
		Mode : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetAgingMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var mode int32
	if options.Args.Value=="disable"{
		mode=0
	}else if options.Args.Value=="enable"{
		mode=1
	}else {
		fmt.Println("select disable/enable. input :"+options.Args.Value)
		return nil
	}
	param := bossopenolt.IntegerValue{ Value : mode}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_IntegervalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetAgingMode(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetAgingMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetAgingMode(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETMODE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Mode ==0{
		response="disable"
	}else{
		response="enable"
	}
	tmp := GetMode{
		DeviceId: Bossresp.DeviceId,
		Mode : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetAgingTime) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	param := bossopenolt.IntegerValue{ Value : int32(options.Args.Value)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_IntegervalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetAgingTime(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetAgingTimeRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetAgingTime(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETAGINGTime_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	tmp := GetAgingTime{
		DeviceId: Bossresp.DeviceId,
		AgingTime : strconv.Itoa(int(Bossresp.AgingTime))+"sec",
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetDeviceMacInfoRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetDeviceMacInfo(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETDEVICEMACINFO_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var mode string
	if Bossresp.VlanMode ==0{
		mode="Pass"
	}else{
		mode="Add/Remove"
	}
	var agingMode string
	if Bossresp.AgingMode == 0{
		agingMode = "disable"
	}else {
		agingMode = "enable"
	}

	tmp := GetDeviceMacInfo{
		DeviceId: Bossresp.DeviceId,
		Mtu : strconv.Itoa(int(Bossresp.Mtu)),
		VlanMode : mode,
		AgingMode :agingMode,
		AgingTime: strconv.Itoa(int(Bossresp.AgingTime))+"sec",
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetSdnTableRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var port int64
	port, _ = strconv.ParseInt(options.Args.PortId, 0, 16)
//	port = strconv.ParseInt("0x1001", 0, 16)
	param := bossopenolt.SetSdnTable{PortId : int32(port), Vid : int32(options.Args.Vid), Pbit: int32(options.Args.Pbit)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetsdntableParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetSdnTable(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_SETSDNTABLE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	tmp := SetSdnTable{
		HashKey : strconv.Itoa(int(Bossresp.HashKey)),
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetSdnTableRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var address int64
	address ,_ = strconv.ParseInt(options.Args.Address, 0, 16)
	param := bossopenolt.GetSdnTable{Address : int32(address)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetsdntableParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetSdnTable(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETSDNTABLE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	tmp := GetSdnTable{
		DeviceId : Bossresp.DeviceId,
		Address : strconv.Itoa(int(Bossresp.Address)),
		PortId : strconv.Itoa(int(Bossresp.PortId)),
		Vlan : Bossresp.Vlan,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetLength) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	var operation int
	if options.Args.Operation == "enable"{
		operation = 0
	}else if options.Args.Operation == "disable" {
		operation =1
	}else if options.Args.Operation == "guard_time"{
		operation=2
	}else{
		fmt.Println("Validate Error")
	}
	defer cancel()
	param := bossopenolt.SetLength{Operation : int32(operation), Value : int32(options.Args.Value)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetlengthParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetLength(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var resultResp string
	if Bossresp.Result == 0{
		resultResp = "success"
	}else{
		resultResp = "fail"
	}

	tmp := GetResponse{
		Result : resultResp,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetLengthRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	var operation int
	if options.Args.Operation == "enable"{
		operation = 0
	}else if options.Args.Operation == "disable" {
		operation =1
	}else if options.Args.Operation == "guard_time"{
		operation=2
	}else{
		fmt.Println("Validate Error")
		return nil
	}
	defer cancel()
	param := bossopenolt.GetLength{Operation : int32(operation)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetlengthParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetLength(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETVALUE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
/*	var resultResp string
	if Bossresp.Result == 0{
		resultResp = "success"
	}else{
		resultResp = "fail"
	}
*/
	tmp := GetValue{
		DeviceId : Bossresp.DeviceId,
		Value : fmt.Sprintf("%v",Bossresp.Value),
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetQuietZone) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.IntegerValue{ Value : int32(options.Args.Value)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_IntegervalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetQuietZone(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetQuietZone) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetQuietZone(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETVALUE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	tmp := GetValue{
		DeviceId: Bossresp.DeviceId,
		Value : strconv.Itoa(int(Bossresp.Value)),
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetFecMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var directionParam int
	if options.Args.Direction == "downstream"{
		directionParam = 0
	}else if options.Args.Direction =="upstream"{
		directionParam = 1
	}else{
		fmt.Println("Validate Error")
		return nil
	}
	param := bossopenolt.SetDirectionMode{ Direction : int32(directionParam), Mode : int32(options.Args.Mode)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetdirectiommodeParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetFecMode(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetFecMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var direct int32

	if options.Args.Direction == "downstream"{
		direct = 0
	}else if options.Args.Direction =="upstream"{
		direct = 1
	}else{
		fmt.Println("Validate Error")
		return nil
	}
	param := bossopenolt.GetDirectionValue{Direction:direct}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetdirectionvalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetFecMode(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETMODE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Mode ==0{
		response="off"
	}else{
		response="on"
	}
	tmp := GetMode{
		DeviceId: Bossresp.DeviceId,
		Mode : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *AddOnuRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	param := bossopenolt.OnuCtrl{OnuId:int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.AddOnu(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_ADDONU_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	
	tmp := AddOnu{
		DeviceId: Bossresp.DeviceId,
		OnuId : strconv.Itoa(int(Bossresp.OnuId)),
		Result : Bossresp.Result,
		Rate : Bossresp.Rate,
		VendorId : Bossresp.VendorId,
		Vssn : Bossresp.Vssn,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *DeleteOnu) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.DeleteOnu25G(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *AddOnuSla) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.AddOnuSla{OnuId : int32(options.Args.OnuId), Tcont : int32(options.Args.Tcont), Type: int32(options.Args.Type), Si: int32(options.Args.Si), Abmin:int32(options.Args.Abmin), Absur : int32(options.Args.Absur)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_AddonuslaParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.AddOnuSla(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *ClearOnuSla) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.ClearOnuSla{OnuId : int32(options.Args.OnuId), Tcont : int32(options.Args.Tcont)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_ClearonuslaParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.ClearOnuSla(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func (options *GetSlaTable) Execute(args []string) error {

	conn, err := NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	ports, err := client.GetSlaTable(ctx, &request)
	if err != nil {
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETSLATABLE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	result := CommandResult{
		Format:    format.Format(outputFormat),
		Filter:    options.Filter,
		OrderBy:   orderBy,
		OutputAs:  toOutputType(options.OutputAs),
		NameLimit: options.NameLimit,
		Data:      ports.Resp,
	}

	GenerateOutput(&result)
	return nil
}

func(options *SetOnuAllocid) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
 	var allocId int64
	allocId,_ =strconv.ParseInt(options.Args.AllocId, 0,16)
	param := bossopenolt.SetOnuAllocid{OnuId : int32(options.Args.OnuId), AllocId : int32(allocId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetonuallocidParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetOnuAllocid(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *DeleteOnuAllocid) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var allocId int64
	allocId,_ =strconv.ParseInt(options.Args.AllocId, 0,16)

	param := bossopenolt.SetOnuAllocid{OnuId : int32(options.Args.OnuId), AllocId : int32(allocId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetonuallocidParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.DelOnuAllocid(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetVssn) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var vssn int64
	vssn, _ = strconv.ParseInt(options.Args.Vssn, 0,16)
	param := bossopenolt.SetOnuVssn{OnuId : int32(options.Args.OnuId), Vssn : int32(vssn)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetonuvssnParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetOnuVssn(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetOnuVssn) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetOnuVssn(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETONUVSSN_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetOnuDistance) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetOnuDistance(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETONUDISTANCE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetBurstDelimiter) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.SetBurstDelimit{Length : int32(options.Args.Length), Delimiter : options.Args.Delimiter}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetburstdelimitParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetBurstDelimiter(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetBurstDelimiter) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetBurstDelimiter(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETBURSTDELIMITER_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetBurstPreamble) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var repeat int64
	repeat,_ = strconv.ParseInt(options.Args.Repeat, 0, 16)
	param := bossopenolt.SetBurstPreamble{Repeat : int32(repeat)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetburstpreambleParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetBurstPreamble(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetBurstPreamble) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetBurstPreamble(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETBURSTPREAMBLE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetBurstVersion) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var pontag int64
	pontag,_ = strconv.ParseInt(options.Args.Pontag, 0, 16)
	param := bossopenolt.SetBurstVersion{Version : options.Args.Version, Index: int32(options.Args.Index), Pontag : int64(pontag)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetburstversionParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetBurstVersion(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetBurstVersion) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetBurstVersion(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETBURSTVERSION_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetBurstProfile) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetBurstProfile(ctx, &request)
	if err != nil{
		return err
	}
	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}	
	GenerateOutput(&result)
	return nil;
}

func(options *GetBurstProfile) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetBurstProfile(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETBURSTPROFILE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetRegisterStatus) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetRegisterStatus(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETREGISTERSTATUS_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetOnuInfo) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetOnuInfo(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETONUINFO_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetOmciStatus) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var updown int32
	if options.Args.Direction =="downstream" {
		updown = 0
	}else if options.Args.Direction =="upstream"{
		updown = 1
	}else {
		fmt.Println("validate Error")
		return nil
	}
	
	param := bossopenolt.GetDirectionValue{Direction : updown}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetdirectionvalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetOmciStatus(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETOMCISTATUS_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetDsOmciOnu) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetDsOmciOnu(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetDsOmciData) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.SetDsOmciData{Control : int32(options.Args.Control), Data : options.Args.Data}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetdsomcidataParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetDsOmciData(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetUsOmciData) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetUsOmciData(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETOMCIDATA_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetTod) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.SetTod{Mode : int32(options.Args.Mode), Time : int32(options.Args.Time)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SettodParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetTod(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetTodRequest) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetTod(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETTOD_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var respMode string
	if Bossresp.Mode ==0 {
		respMode = "off"
	}else {
		respMode ="on"
	}
	Result :=GetTod{
		DeviceId : Bossresp.DeviceId,
		Mode : respMode,
		Time : strconv.Itoa(int(Bossresp.Time)),
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Result,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetDataMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var updown int
	if options.Args.Direction =="Scrambler"{
		updown=0
	}else if options.Args.Direction=="Descrambler"{
		updown=1
	}else{
		fmt.Println("Validate Error")
		return nil
	}
	param := bossopenolt.SetDirectionMode{Direction : int32(updown), Mode : int32(options.Args.Mode)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetdirectiommodeParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetDataMode(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetDataMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	var updown int
	if options.Args.Direction =="Scrambler"{
		updown=0
	}else if options.Args.Direction=="Descrambler"{
		updown=1
	}else{
		fmt.Println("Validate Error")
		return nil
	}

	param := bossopenolt.GetDirectionValue{Direction : int32(updown)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_GetdirectionvalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetDataMode(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETMODE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Mode == 0{
		response = "bypass"
	}else if Bossresp.Mode==1{
		response = "normal"
	}else{
		response = "unknown"
	}
	tmp := GetMode{
		DeviceId :  Bossresp.DeviceId,
		Mode : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetFecDecMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	param := bossopenolt.IntegerValue{Value : int32(options.Args.Value)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_IntegervalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetFecDecMode(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetFecDecMode) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetFecDecMode(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETMODE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Mode == 0{
		response = "on"
	}else if Bossresp.Mode ==1{
		response = "off"
	}else{
		response="unknown"
	}
	tmp := GetMode{
		DeviceId :  Bossresp.DeviceId,
		Mode : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetDelimiter) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	param := bossopenolt.StringValue{Value : options.Args.Value}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_StringvalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetDelimiter(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetDelimiter) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetDelimiter(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETVALUE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetErrorPermit) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	param := bossopenolt.IntegerValue{Value : int32(options.Args.Value)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_IntegervalueParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetErrorPermit(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetErrorPermit) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetErrorPermit(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETVALUE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetPmControl) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()

	param := bossopenolt.SetPmControl{OnuId : int32(options.Args.OnuId), Mode : int32(options.Args.Mode), PowerTime: int32(options.Args.PowerTime), AwareTime: int32(options.Args.AwareTime)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetpmcontrolParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetPmControl(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}
func(options *GetPmControl) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	param := bossopenolt.OnuCtrl{OnuId:int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetPmControl(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETPMCONTROL_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetPmTable) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	param := bossopenolt.OnuCtrl{OnuId:int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.GetPmTable(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETPMTABLE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetSAOn) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetSAOn(ctx, &request)

	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetSAOff) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.OnuCtrl{OnuId : int32(options.Args.OnuId)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_OnuctrlParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetSAOff(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func (options *CreateDeviceHandler) Execute(args []string) error {
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

		_, err := client.CreateDeviceHandler(ctx, &id)
		if err != nil {
			Error.Printf("Error while enabling '%s': %s\n", i, err)
			lastErr = err
			continue
		}
		fmt.Printf("%s\n", i)
	}

	if lastErr != nil {
		return NoReportErr
	}
	return nil
}

func(options *SetSliceBw) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.SetSliceBw{Slice: int32(options.Args.Slice), Bw: int32(options.Args.Bw)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_SetslicebwParam{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetSliceBw(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_RESPONSE_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}
	var response string
	if Bossresp.Result == 0{
		response = "success"
	}else{
		response = "fail"
	}
	tmp := GetResponse{
		Result : response,
	}
	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : tmp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetSliceBw) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetSliceBw(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_GETSLICEBW_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *GetSlaV2) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id)}
	Bossresp, err := client.GetSlaV2(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_SLAV2_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}

func(options *SetSlaV2) Execute(args []string)error{
	conn,err:=NewConnection()
	defer conn.Close()

	client := bossopenolt.NewBossOpenoltClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), GlobalConfig.Current().Grpc.Timeout)
	defer cancel()
	
	param := bossopenolt.SetSlaV2{OnuId: int32(options.Args.OnuId), Tcont: int32(options.Args.Tcont), Slice : int32(options.Args.Slice), CoDba: int32(options.Args.CoDba), Type : int32(options.Args.Type), Rf: int32(options.Args.Rf), Ra: int32(options.Args.Ra), Rn : int32(options.Args.Rn)}
	requestParam := bossopenolt.ParamFields{Data:&bossopenolt.ParamFields_Setslav2Param{&param} }
	request := bossopenolt.BossRequest{DeviceId:string(options.Args.Id), Param : &requestParam}
	Bossresp, err := client.SetSlaV2(ctx, &request)
	if err != nil{
		return err
	}

	outputFormat := CharReplacer.Replace(options.Format)
	if outputFormat == "" {
		outputFormat = GetCommandOptionWithDefault("device-ports", "format", DEFAULT_DEVICE_BOSS_SLAV2_FORMAT)
	}

	orderBy := options.OrderBy
	if orderBy == "" {
		orderBy = GetCommandOptionWithDefault("device-ports", "order", "")
	}

	result := CommandResult{
		Format : format.Format(outputFormat),
		OrderBy : orderBy,
		Data : Bossresp,
	}
	GenerateOutput(&result)
	return nil;
}


