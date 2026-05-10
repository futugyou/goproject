package domains

type Mode uint32

const (
	ModeSTANDARD              Mode = 0
	ModePROPAGATE_INHERITANCE Mode = 1
)

type SCIMSchemaAttributeMutabilities uint32

const (
	SCIMSchemaAttributeMutabilitiesREADONLY SCIMSchemaAttributeMutabilities = iota
	SCIMSchemaAttributeMutabilitiesREADWRITE
	SCIMSchemaAttributeMutabilitiesIMMUTABLE
	SCIMSchemaAttributeMutabilitiesWRITEONLY
)

func (s SCIMSchemaAttributeMutabilities) Name() string {
	switch s {
	case SCIMSchemaAttributeMutabilitiesREADONLY:
		return "readOnly"
	case SCIMSchemaAttributeMutabilitiesREADWRITE:
		return "readWrite"
	case SCIMSchemaAttributeMutabilitiesIMMUTABLE:
		return "immutable"
	case SCIMSchemaAttributeMutabilitiesWRITEONLY:
		return "writeOnly"
	}

	return ""
}

type SCIMSchemaAttributeReturned uint32

const (
	SCIMSchemaAttributeReturnedALWAYS SCIMSchemaAttributeReturned = iota
	SCIMSchemaAttributeReturnedNEVER
	SCIMSchemaAttributeReturnedDEFAULT
	SCIMSchemaAttributeReturnedREQUEST
)

func (s SCIMSchemaAttributeReturned) Name() string {
	switch s {
	case SCIMSchemaAttributeReturnedALWAYS:
		return "always"
	case SCIMSchemaAttributeReturnedNEVER:
		return "never"
	case SCIMSchemaAttributeReturnedDEFAULT:
		return "default"
	case SCIMSchemaAttributeReturnedREQUEST:
		return "request"
	}

	return ""
}

type SCIMSchemaAttributeTypes uint32

const (
	SCIMSchemaAttributeTypesSTRING SCIMSchemaAttributeTypes = iota
	SCIMSchemaAttributeTypesBOOLEAN
	SCIMSchemaAttributeTypesINTEGER
	SCIMSchemaAttributeTypesDATETIME
	SCIMSchemaAttributeTypesREFERENCE
	SCIMSchemaAttributeTypesCOMPLEX
	SCIMSchemaAttributeTypesDECIMAL
	SCIMSchemaAttributeTypesBINARY
)

type SCIMSchemaAttributeUniqueness uint32

const (
	SCIMSchemaAttributeUniquenessNONE SCIMSchemaAttributeUniqueness = iota
	SCIMSchemaAttributeUniquenessSERVER
	SCIMSchemaAttributeUniquenessGLOBAL
)

func (s SCIMSchemaAttributeUniqueness) Name() string {
	switch s {
	case SCIMSchemaAttributeUniquenessNONE:
		return "none"
	case SCIMSchemaAttributeUniquenessSERVER:
		return "server"
	case SCIMSchemaAttributeUniquenessGLOBAL:
		return "global"
	}

	return ""
}
