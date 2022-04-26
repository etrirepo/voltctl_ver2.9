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

// PhysicalPathTerminationPointXdslUniPart3ClassID is the 16-bit ID for the OMCI
// Managed entity Physical path termination point xDSL UNI part 3
const PhysicalPathTerminationPointXdslUniPart3ClassID = ClassID(427) // 0x01ab

var physicalpathterminationpointxdslunipart3BME *ManagedEntityDefinition

// PhysicalPathTerminationPointXdslUniPart3 (Class ID: #427 / 0x01ab)
//	This ME represents the point in the ONU where physical paths terminate on an xDSL CO modem
//	(xTU-C). Standards and chip sets support several forms of DSL, including VDSL2 and FAST, and the
//	xDSL ME family is used for all of them, with specific extensions for technology variations.
//
//	The ONU creates or deletes an instance of this ME at the same time it creates or deletes the
//	corresponding PPTP xDSL UNI part 1.
//
//	Relationships
//		An instance of this ME is associated with each instance of a real or preprovisioned xDSL port
//
//	Attributes
//		Managed Entity Id
//			This attribute uniquely identifies each instance of this ME. This 2 byte number indicates the
//			physical position of the UNI. The six LSBs of the first byte are the slot ID, defined in clause
//			9.1.5. The two MSBs indicate the channel number in some of the implicitly linked MEs, and must
//			be 0 in the PPTP itself. This reduces the possible number of physical slots to 64. The second
//			byte is the port ID, with range 1..255. (R) (mandatory) (2 bytes)
//
//		Fast Line Configuration Profile
//			This attribute points to an instance of the FAST line configuration profiles (part 1, 2, 3 and
//			4) MEs, also to FAST vectoring line configuration extension MEs. Upon ME instantiation, the ONU
//			sets this attribute to 0, a null pointer. (R, W) (mandatory) (2 bytes)
//
//		Fast Data Path Configuration Profile
//			This attribute points to an instance of the FAST data configuration profile that defines data
//			path parameters. Upon ME instantiation, the ONU sets this attribute to 0, a null pointer. (R, W)
//			(optional) (2 bytes)
//
//		Fast Channel Configuration Profile For Bearer Channel 0 Downstream
//			This attribute points to an instance of the FAST channel configuration profile that defines
//			channel parameters. Upon ME instantiation, the ONU sets this attribute to 0, a null pointer.
//			(R,-W) (optional) (2-bytes) (R,-W) (optional) (2-bytes)
//
//		Fast Xdsl Channel Configuration Profile For Bearer Channel 0 Upstream
//			This attribute points to an instance of the FAST channel configuration profile that defines
//			channel parameters. Upon ME instantiation, the ONU sets this attribute to 0, a null pointer
//			(R,-W) (optional) (2-bytes)
//
type PhysicalPathTerminationPointXdslUniPart3 struct {
	ManagedEntityDefinition
	Attributes AttributeValueMap
}

func init() {
	physicalpathterminationpointxdslunipart3BME = &ManagedEntityDefinition{
		Name:    "PhysicalPathTerminationPointXdslUniPart3",
		ClassID: 427,
		MessageTypes: mapset.NewSetWith(
			Get,
			Set,
		),
		AllowedAttributeMask: 0xf000,
		AttributeDefinitions: AttributeDefinitionMap{
			0: Uint16Field("ManagedEntityId", PointerAttributeType, 0x0000, 0, mapset.NewSetWith(Read), false, false, false, 0),
			1: Uint16Field("FastLineConfigurationProfile", UnsignedIntegerAttributeType, 0x8000, 0, mapset.NewSetWith(Read, Write), false, false, false, 1),
			2: Uint16Field("FastDataPathConfigurationProfile", UnsignedIntegerAttributeType, 0x4000, 0, mapset.NewSetWith(Read, Write), false, true, false, 2),
			3: Uint16Field("FastChannelConfigurationProfileForBearerChannel0Downstream", UnsignedIntegerAttributeType, 0x2000, 0, mapset.NewSetWith(Read, Write), false, true, false, 3),
			4: Uint16Field("FastXdslChannelConfigurationProfileForBearerChannel0Upstream", UnsignedIntegerAttributeType, 0x1000, 0, mapset.NewSetWith(Read, Write), false, true, false, 4),
		},
		Access:  CreatedByOnu,
		Support: UnknownSupport,
	}
}

// NewPhysicalPathTerminationPointXdslUniPart3 (class ID 427) creates the basic
// Managed Entity definition that is used to validate an ME of this type that
// is received from or transmitted to the OMCC.
func NewPhysicalPathTerminationPointXdslUniPart3(params ...ParamData) (*ManagedEntity, OmciErrors) {
	return NewManagedEntity(*physicalpathterminationpointxdslunipart3BME, params...)
}
