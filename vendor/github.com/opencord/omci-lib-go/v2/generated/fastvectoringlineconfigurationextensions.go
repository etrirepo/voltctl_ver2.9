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

// FastVectoringLineConfigurationExtensionsClassID is the 16-bit ID for the OMCI
// Managed entity FAST vectoring line configuration extensions
const FastVectoringLineConfigurationExtensionsClassID = ClassID(434) // 0x01b2

var fastvectoringlineconfigurationextensionsBME *ManagedEntityDefinition

// FastVectoringLineConfigurationExtensions (Class ID: #434 / 0x01b2)
//	This ME extends FAST line configuration MEs with attributes that are specific to vectoring. An
//	instance of this ME is created and deleted by the OLT.
//
//	Relationships
//		An instance of this ME may be associated with zero or more instances of an xDSL UNI.////		The overall FAST line configuration MEs is modelled in several parts, all of which are
//		associated together through a common ME ID (the client PPTP xDSL UNI part 3 has a single
//		pointer, which refers to the entire set of line configuration parts).
//
//	Attributes
//		Managed Entity Id
//			This attribute uniquely identifies each instance of this ME. The value 0 is reserved. (R, set-
//			by-create) (mandatory) (2 bytes)
//
//		Fext Cancellation Enabling_Disabling Upstream Fext_To_Cancel_Enableus
//			FEXT cancellation enabling/disabling upstream (FEXT_TO_CANCEL_ENABLEus): A value of 1 enables
//			and a value of 0 disables FEXT cancellation in the upstream direction from all the other
//			vectored lines into the line in the vectored group. See clause 7.1.7.2 of [ITU-T G.997.2].
//			(R,-W) (mandatory) (1-byte)
//
//		Fext Cancellation Enabling_Disabling Downstream Fext_To_Cancel_Enableds
//			FEXT cancellation enabling/disabling downstream (FEXT_TO_CANCEL_ENABLEds): A value of 1 enables
//			and a value of 0 disables FEXT cancellation in the downstream direction from all the other
//			vectored lines into the line in the vectored group. See clause 7.1.7.1 of [ITUT-G.997.2]. (R,-W)
//			(mandatory) (1-byte)
//
type FastVectoringLineConfigurationExtensions struct {
	ManagedEntityDefinition
	Attributes AttributeValueMap
}

func init() {
	fastvectoringlineconfigurationextensionsBME = &ManagedEntityDefinition{
		Name:    "FastVectoringLineConfigurationExtensions",
		ClassID: 434,
		MessageTypes: mapset.NewSetWith(
			Create,
			Delete,
			Get,
			Set,
		),
		AllowedAttributeMask: 0xc000,
		AttributeDefinitions: AttributeDefinitionMap{
			0: Uint16Field("ManagedEntityId", PointerAttributeType, 0x0000, 0, mapset.NewSetWith(Read, SetByCreate), false, false, false, 0),
			1: ByteField("FextCancellationEnablingDisablingUpstreamFextToCancelEnableus", UnsignedIntegerAttributeType, 0x8000, 0, mapset.NewSetWith(Read, Write), false, false, false, 1),
			2: ByteField("FextCancellationEnablingDisablingDownstreamFextToCancelEnableds", UnsignedIntegerAttributeType, 0x4000, 0, mapset.NewSetWith(Read, Write), false, false, false, 2),
		},
		Access:  CreatedByOlt,
		Support: UnknownSupport,
	}
}

// NewFastVectoringLineConfigurationExtensions (class ID 434) creates the basic
// Managed Entity definition that is used to validate an ME of this type that
// is received from or transmitted to the OMCC.
func NewFastVectoringLineConfigurationExtensions(params ...ParamData) (*ManagedEntity, OmciErrors) {
	return NewManagedEntity(*fastvectoringlineconfigurationextensionsBME, params...)
}
