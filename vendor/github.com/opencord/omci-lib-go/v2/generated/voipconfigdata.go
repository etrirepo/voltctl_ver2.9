/*
 * Copyright (c) 2018 - present.  Boling Consulting Solutions (bcsw.net)
 * Copyright 2020-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * NOTE: This file was generated, manual edits will be overwritten!
 *
 * Generated by 'goCodeGenerator.py':
 *              https://github.com/cboling/OMCI-parser/README.md
 */

package generated

import "github.com/deckarep/golang-set"

// VoipConfigDataClassID is the 16-bit ID for the OMCI
// Managed entity VoIP config data
const VoipConfigDataClassID = ClassID(138) // 0x008a

var voipconfigdataBME *ManagedEntityDefinition

// VoipConfigData (Class ID: #138 / 0x008a)
//	The VoIP configuration data ME defines the configuration for VoIP in the ONU. The OLT uses this
//	ME to discover the VoIP signalling protocols and configuration methods supported by this ONU.
//	The OLT then uses this ME to select the desired signalling protocol and configuration method.
//	The entity is conditionally required for ONUs that offer VoIP services.
//
//	An ONU that supports VoIP services automatically creates an instance of this ME.
//
//	Relationships
//		One instance of this ME is associated with the ONU.
//
//	Attributes
//		Managed Entity Id
//			This attribute uniquely identifies each instance of this ME. There is only one instance, number
//			0. (R) (mandatory) (2-bytes)
//
//		Available Signalling Protocols
//			This attribute is a bit map that defines the VoIP signalling protocols supported in the ONU. The
//			bit value 1 specifies that the ONU supports the associated protocol.
//
//			1	(LSB)	SIP
//
//			2	ITU-T H.248
//
//			3	MGCP
//
//			(R) (mandatory) (1-byte)
//
//		Signalling Protocol Used
//			0xFF	Selected by non-OMCI management interface
//
//			(R,-W) (mandatory) (1-byte)
//
//			This attribute specifies the VoIP signalling protocol to use. Only one type of protocol is
//			allowed at a time. Valid values are:
//
//			0	None
//
//			1	SIP
//
//			2	ITU-T H.248
//
//			3	MGCP
//
//		Available Voip Configuration Methods
//			This attribute is a bit map that indicates the capabilities of the ONU with regard to VoIP
//			service configuration. The bit value 1 specifies that the ONU supports the associated
//			capability.
//
//			1 (LSB)	ONU capable of using the OMCI to configure its VoIP services.
//
//			2	ONU capable of working with configuration file retrieval to configure its VoIP services.
//
//			3	ONU capable of working with [BBF TR-069] to configure its VoIP services.
//
//			4	ONU capable of working with IETF sipping config framework to configure its VoIP services.
//
//			Bits 5..24 are reserved by ITU-T. Bits 25..32 are reserved for proprietary vendor configuration
//			capabilities. (R) (mandatory) (4-bytes)
//
//		Voip Configuration Method Used
//			Specifies which method is used to configure the ONU's VoIP service.
//
//			0	Do not configure - ONU default
//
//			1	OMCI
//
//			2	Configuration file retrieval
//
//			3	BBF TR-069
//
//			4	IETF sipping config framework
//
//			5..240	Reserved by ITU-T
//
//			241..255	Reserved for proprietary vendor configuration methods
//
//			(R,-W) (mandatory) (1-byte)
//
//		Voip Configuration Address Pointer
//			If this attribute is set to any value other than a null pointer, it points to a network address
//			ME, which indicates the address of the server to contact using the method indicated in the VoIP
//			configuration method used attribute. This attribute is only relevant for non-OMCI configuration
//			methods.
//
//			If this attribute is set to a null pointer, no address is defined by this attribute. However,
//			the address may be defined by other methods, such as deriving it from the ONU identifier
//			attribute of the IP host config data ME and using a well-known URI schema.
//
//			The default value is 0xFFFF (R,-W) (mandatory) (2-bytes)
//
//		Voip Configuration State
//			Indicates the status of the ONU VoIP service.
//
//			0	Inactive: configuration retrieval has not been attempted
//
//			1	Active: configuration was retrieved
//
//			2	Initializing: configuration is now being retrieved
//
//			3	Fault: configuration retrieval process failed
//
//			Other values are reserved. At ME instantiation, the ONU sets this attribute to 0. (R)
//			(mandatory) (1-byte)
//
//		Retrieve Profile
//			This attribute provides a means by which the ONU may be notified that a new VoIP profile should
//			be retrieved. By setting this attribute, the OLT triggers the ONU to retrieve a new profile. The
//			actual value in the set action is ignored because it is the action of setting that is important.
//			(W) (mandatory) (1-byte)
//
//		Profile Version
//			This attribute is a character string that identifies the version of the last retrieved profile.
//			(R) (mandatory) (25-bytes)
//
type VoipConfigData struct {
	ManagedEntityDefinition
	Attributes AttributeValueMap
}

func init() {
	voipconfigdataBME = &ManagedEntityDefinition{
		Name:    "VoipConfigData",
		ClassID: 138,
		MessageTypes: mapset.NewSetWith(
			Get,
			Set,
		),
		AllowedAttributeMask: 0xff00,
		AttributeDefinitions: AttributeDefinitionMap{
			0: Uint16Field("ManagedEntityId", PointerAttributeType, 0x0000, 0, mapset.NewSetWith(Read), false, false, false, 0),
			1: ByteField("AvailableSignallingProtocols", UnsignedIntegerAttributeType, 0x8000, 0, mapset.NewSetWith(Read), false, false, false, 1),
			2: ByteField("SignallingProtocolUsed", UnsignedIntegerAttributeType, 0x4000, 0, mapset.NewSetWith(Read, Write), false, false, false, 2),
			3: Uint32Field("AvailableVoipConfigurationMethods", UnsignedIntegerAttributeType, 0x2000, 0, mapset.NewSetWith(Read), false, false, false, 3),
			4: ByteField("VoipConfigurationMethodUsed", UnsignedIntegerAttributeType, 0x1000, 0, mapset.NewSetWith(Read, Write), false, false, false, 4),
			5: Uint16Field("VoipConfigurationAddressPointer", UnsignedIntegerAttributeType, 0x0800, 0, mapset.NewSetWith(Read, Write), false, false, false, 5),
			6: ByteField("VoipConfigurationState", UnsignedIntegerAttributeType, 0x0400, 0, mapset.NewSetWith(Read), false, false, false, 6),
			7: ByteField("RetrieveProfile", UnsignedIntegerAttributeType, 0x0200, 0, mapset.NewSetWith(Write), false, false, false, 7),
			8: MultiByteField("ProfileVersion", OctetsAttributeType, 0x0100, 25, toOctets("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="), mapset.NewSetWith(Read), true, false, false, 8),
		},
		Access:  CreatedByOnu,
		Support: UnknownSupport,
		Alarms: AlarmMap{
			0:  "VCD config server name",
			1:  "VCD config server reach",
			2:  "VCD config server connect",
			3:  "VCD config server validate",
			4:  "VCD config server auth",
			5:  "VCD config server timeout",
			6:  "VCD config server fail",
			7:  "VCD config file error",
			8:  "VCD subscription name",
			9:  "VCD subscription reach",
			10: "VCD subscription connect",
			11: "VCD subscription validate",
			12: "VCD subscription auth",
			13: "VCD subscription timeout",
			14: "VCD subscription fail",
			15: "VCD reboot request",
		},
	}
}

// NewVoipConfigData (class ID 138) creates the basic
// Managed Entity definition that is used to validate an ME of this type that
// is received from or transmitted to the OMCC.
func NewVoipConfigData(params ...ParamData) (*ManagedEntity, OmciErrors) {
	return NewManagedEntity(*voipconfigdataBME, params...)
}
