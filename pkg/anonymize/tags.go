package anonymize

import (
	"github.com/suyashkumar/dicom/pkg/tag"
)

type AnonymizeAction uint8

const (
	ReplaceAction AnonymizeAction = iota
	EmptyAction
	DeleteAction
	ReplaceUIDAction
	EmptyOrReplaceAction
	DeleteOrEmptyAction
	DeleteOrReplaceAction
	DeleteOrEmptyOrReplaceAction
	DeleteOrEmptyOrReplaceUIDAction
)

// To be anonymized tag and its anonymization action
type TBA struct {
	tag   tag.Tag
	value AnonymizeAction
}

// Tags's value to be replaced
var D = []TBA{
	{tag.GraphicAnnotationSequence, ReplaceAction},        // Graphic Annotation Sequence
	{tag.PersonIdentificationCodeSequence, ReplaceAction}, // Person Identification Code Sequence
	{tag.PersonName, ReplaceAction},                       // Person Name
	{tag.VerifyingObserverName, ReplaceAction},            // Verifying Observer Name
	{tag.VerifyingObserverSequence, ReplaceAction},        // Verifying Observer Sequence
}

// Tags's value to be emptied
var Z = []TBA{
	{tag.AccessionNumber, EmptyAction},    // Accession Number
	{tag.ContentCreatorName, EmptyAction}, // Content Creator's Name
	{
		tag.FillerOrderNumberImagingServiceRequest,
		EmptyAction,
	}, // Filler Order Number / Imaging Service Request
	{tag.PatientID, EmptyAction},        // Patient ID
	{tag.PatientBirthDate, EmptyAction}, // Patient's Birth Date
	{tag.PatientName, EmptyAction},      // Patient's Name
	{tag.PatientSex, EmptyAction},       // Patient's Sex
	{
		tag.PlacerOrderNumberImagingServiceRequest,
		EmptyAction,
	}, // Placer Order Number / Imaging Service Request
	{tag.ReferringPhysicianAddress, EmptyAction}, // Referring Physician's Name
	{tag.StudyDate, EmptyAction},                 // Study Date
	{tag.StudyID, EmptyAction},                   // Study ID
	{tag.StudyTime, EmptyAction},                 // Study Time
	{
		tag.VerifyingObserverIdentificationCodeSequence,
		EmptyAction,
	}, // Verifying Observer Identification Code Sequence
	{tag.ReferringPhysicianName, EmptyAction}, // Referring Physician's Name
}

// Tags to be deleted
var X = []TBA{
	{tag.ACR_NEMA_AcquisitionComments, DeleteAction},        // Acquisition Comments
	{tag.AcquisitionContextSequence, DeleteAction},          // Acquisition Context Sequence
	{tag.AcquisitionProtocolDescription, DeleteAction},      // Acquisition Protocol Description
	{tag.ActualHumanPerformersSequence, DeleteAction},       // Actual Human Performers Sequence
	{tag.AdditionalPatientHistory, DeleteAction},            // Additional Patient's History
	{tag.AdmissionID, DeleteAction},                         // Admission ID
	{tag.AdmittingDate, DeleteAction},                       // Admitting Date
	{tag.AdmittingDiagnosesCodeSequence, DeleteAction},      // Admitting Diagnoses Code Sequence
	{tag.AdmittingDiagnosesDescription, DeleteAction},       // Admitting Diagnoses Description
	{tag.AdmittingTime, DeleteAction},                       // Admitting Time
	{tag.Tag{Group: 0x0000, Element: 0x1000}, DeleteAction}, // Affected SOP Instance UID
	{tag.Tag{Group: 0x0010, Element: 0x2110}, DeleteAction}, // Allergies
	{tag.Tag{Group: 0x4000, Element: 0x0010}, DeleteAction}, // Arbitrary
	{tag.Tag{Group: 0x0040, Element: 0xA078}, DeleteAction}, // Author Observer Sequence
	{tag.Tag{Group: 0x0010, Element: 0x1081}, DeleteAction}, // Branch of Service
	{tag.Tag{Group: 0x0018, Element: 0x1007}, DeleteAction}, // Cassette ID
	{
		tag.Tag{Group: 0x0040, Element: 0x0280},
		DeleteAction,
	}, // Comments on the Performed Procedure Step
	{
		tag.Tag{Group: 0x0040, Element: 0x3001},
		DeleteAction,
	}, // Confidentiality Constraint on Patient Data Description
	{
		tag.Tag{Group: 0x0070, Element: 0x0086},
		DeleteAction,
	}, // Content Creator's Identification Code Sequence
	{tag.Tag{Group: 0x0040, Element: 0xA730}, DeleteAction}, // Content Sequence
	{tag.Tag{Group: 0x0018, Element: 0xA003}, DeleteAction}, // Contribution Description
	{tag.Tag{Group: 0x0010, Element: 0x2150}, DeleteAction}, // Country of Residence
	{tag.Tag{Group: 0x0038, Element: 0x0300}, DeleteAction}, // Current Patient Location
	{tag.Tag{Group: 0x0008, Element: 0x0025}, DeleteAction}, // Curve Date
	{tag.Tag{Group: 0x0008, Element: 0x0035}, DeleteAction}, // Curve Time
	{tag.Tag{Group: 0x0040, Element: 0xA07C}, DeleteAction}, // Custodial Organization Sequence
	{tag.Tag{Group: 0xFFFC, Element: 0xFFFC}, DeleteAction}, // Data Set Trailing Padding
	{tag.Tag{Group: 0x0008, Element: 0x2111}, DeleteAction}, // Derivation Description
	{tag.Tag{Group: 0x0400, Element: 0x0100}, DeleteAction}, // Digital Signature UID
	{tag.Tag{Group: 0xFFFA, Element: 0xFFFA}, DeleteAction}, // Digital Signatures Sequence
	{tag.Tag{Group: 0x0038, Element: 0x0040}, DeleteAction}, // Discharge Diagnosis Description
	{tag.Tag{Group: 0x4008, Element: 0x011A}, DeleteAction}, // Distribution Address
	{tag.Tag{Group: 0x4008, Element: 0x0119}, DeleteAction}, // Distribution Name
	{tag.Tag{Group: 0x0010, Element: 0x2160}, DeleteAction}, // Ethnic Group
	{tag.Tag{Group: 0x0020, Element: 0x9158}, DeleteAction}, // Frame Comments
	{tag.Tag{Group: 0x0018, Element: 0x1008}, DeleteAction}, // Gantry ID
	{tag.Tag{Group: 0x0018, Element: 0x1005}, DeleteAction}, // Generator ID
	{tag.Tag{Group: 0x0040, Element: 0x4037}, DeleteAction}, // Human Performers Name
	{tag.Tag{Group: 0x0040, Element: 0x4036}, DeleteAction}, // Human Performers Organization
	{tag.Tag{Group: 0x0088, Element: 0x0200}, DeleteAction}, // Icon Image Sequence
	{tag.Tag{Group: 0x0008, Element: 0x4000}, DeleteAction}, // Identifying Comments
	{tag.Tag{Group: 0x0020, Element: 0x4000}, DeleteAction}, // Image Comments
	{tag.Tag{Group: 0x0028, Element: 0x4000}, DeleteAction}, // Image Presentation Comments
	{tag.Tag{Group: 0x0040, Element: 0x2400}, DeleteAction}, // Imaging Service Request Comments
	{tag.Tag{Group: 0x4008, Element: 0x0300}, DeleteAction}, // Impressions
	{tag.Tag{Group: 0x0008, Element: 0x0081}, DeleteAction}, // Institution Address
	{tag.Tag{Group: 0x0008, Element: 0x1040}, DeleteAction}, // Institutional Department Name
	{tag.Tag{Group: 0x0010, Element: 0x1050}, DeleteAction}, // Insurance Plan Identification
	{
		tag.Tag{Group: 0x0040, Element: 0x1011},
		DeleteAction,
	}, // Intended Recipients of Results Identification Sequence
	{tag.Tag{Group: 0x4008, Element: 0x0111}, DeleteAction}, // Interpretation Approver Sequence
	{tag.Tag{Group: 0x4008, Element: 0x010C}, DeleteAction}, // Interpretation Author
	{tag.Tag{Group: 0x4008, Element: 0x0115}, DeleteAction}, // Interpretation Diagnosis Description
	{tag.Tag{Group: 0x4008, Element: 0x0202}, DeleteAction}, // Interpretation ID Issuer
	{tag.Tag{Group: 0x4008, Element: 0x0102}, DeleteAction}, // Interpretation Recorder
	{tag.Tag{Group: 0x4008, Element: 0x010B}, DeleteAction}, // Interpretation Text
	{tag.Tag{Group: 0x4008, Element: 0x010A}, DeleteAction}, // Interpretation Transcriber
	{tag.Tag{Group: 0x0038, Element: 0x0011}, DeleteAction}, // Issuer of Admission ID
	{tag.Tag{Group: 0x0010, Element: 0x0021}, DeleteAction}, // Issuer of Patient ID
	{tag.Tag{Group: 0x0038, Element: 0x0061}, DeleteAction}, // Issuer of Service Episode ID
	{tag.Tag{Group: 0x0010, Element: 0x21D0}, DeleteAction}, // Last Menstrual Date
	{tag.Tag{Group: 0x0400, Element: 0x0404}, DeleteAction}, // MAC
	{tag.Tag{Group: 0x0010, Element: 0x2000}, DeleteAction}, // Medical Alerts
	{tag.Tag{Group: 0x0010, Element: 0x1090}, DeleteAction}, // Medical Record Locator
	{tag.Tag{Group: 0x0010, Element: 0x1080}, DeleteAction}, // Military Rank
	{tag.Tag{Group: 0x0400, Element: 0x0550}, DeleteAction}, // Modified Attributes Sequence
	{tag.Tag{Group: 0x0020, Element: 0x3406}, DeleteAction}, // Modified Image Description
	{tag.Tag{Group: 0x0020, Element: 0x3401}, DeleteAction}, // Modifying Device ID
	{tag.Tag{Group: 0x0020, Element: 0x3404}, DeleteAction}, // Modifying Device Manufacturer
	{tag.Tag{Group: 0x0008, Element: 0x1060}, DeleteAction}, // Name of Physician study
	{
		tag.Tag{Group: 0x0040, Element: 0x1010},
		DeleteAction,
	}, // Names of Intended Recipient of Results
	{tag.Tag{Group: 0x0010, Element: 0x2180}, DeleteAction}, // Occupation
	{tag.Tag{Group: 0x0400, Element: 0x0561}, DeleteAction}, // Original Attributes Sequence
	{tag.Tag{Group: 0x0040, Element: 0x2010}, DeleteAction}, // Order Callback Phone Number
	{tag.Tag{Group: 0x0040, Element: 0x2008}, DeleteAction}, // Order Entered By
	{tag.Tag{Group: 0x0040, Element: 0x2009}, DeleteAction}, // Order Enterer Location
	{tag.Tag{Group: 0x0010, Element: 0x1000}, DeleteAction}, // Other Patient IDs
	{tag.Tag{Group: 0x0010, Element: 0x1002}, DeleteAction}, // Other Patient IDs Sequence
	{tag.Tag{Group: 0x0010, Element: 0x1001}, DeleteAction}, // Other Patient Names
	{tag.Tag{Group: 0x0008, Element: 0x0024}, DeleteAction}, // Overlay Date
	{tag.Tag{Group: 0x0008, Element: 0x0034}, DeleteAction}, // Overlay Time
	{tag.Tag{Group: 0x0040, Element: 0xA07A}, DeleteAction}, // Participant Sequence
	{tag.Tag{Group: 0x0010, Element: 0x1040}, DeleteAction}, // Patient Address
	{tag.Tag{Group: 0x0010, Element: 0x4000}, DeleteAction}, // Patient Comments
	{tag.Tag{Group: 0x0038, Element: 0x0500}, DeleteAction}, // Patient State
	{tag.Tag{Group: 0x0040, Element: 0x1004}, DeleteAction}, // Patient Transport Arrangements
	{tag.Tag{Group: 0x0010, Element: 0x1010}, DeleteAction}, // Patient's Age
	{tag.Tag{Group: 0x0010, Element: 0x1005}, DeleteAction}, // Patient's Birth Name
	{tag.Tag{Group: 0x0010, Element: 0x0032}, DeleteAction}, // Patient's Birth Time
	{tag.Tag{Group: 0x0038, Element: 0x0400}, DeleteAction}, // Patient's Institution Residence
	{
		tag.Tag{Group: 0x0010, Element: 0x0050},
		DeleteAction,
	}, // Patient's Insurance Plan Code Sequence
	{tag.Tag{Group: 0x0010, Element: 0x1060}, DeleteAction}, // Patient's Mother's Birth Name
	{
		tag.Tag{Group: 0x0010, Element: 0x0101},
		DeleteAction,
	}, // Patient's Primary Language Code Sequence
	{
		tag.Tag{Group: 0x0010, Element: 0x0102},
		DeleteAction,
	}, // Patient's Primary Language Modifier Code Sequence
	{tag.Tag{Group: 0x0010, Element: 0x21F0}, DeleteAction}, // Patient's Religious Preference
	{tag.Tag{Group: 0x0010, Element: 0x1020}, DeleteAction}, // Patient's Size
	{tag.Tag{Group: 0x0010, Element: 0x2154}, DeleteAction}, // Patient's Telephone Numbers
	{tag.Tag{Group: 0x0010, Element: 0x1030}, DeleteAction}, // Patient's Weight
	{tag.Tag{Group: 0x0040, Element: 0x0243}, DeleteAction}, // Performed Location
	{tag.Tag{Group: 0x0040, Element: 0x0254}, DeleteAction}, // Performed Procedure Step Description
	{tag.Tag{Group: 0x0040, Element: 0x0250}, DeleteAction}, // Performed Procedure Step End Date
	{tag.Tag{Group: 0x0040, Element: 0x0251}, DeleteAction}, // Performed Procedure Step End Time
	{tag.Tag{Group: 0x0040, Element: 0x0253}, DeleteAction}, // Performed Procedure Step ID
	{tag.Tag{Group: 0x0040, Element: 0x0244}, DeleteAction}, // Performed Procedure Step Start Date
	{tag.Tag{Group: 0x0040, Element: 0x0245}, DeleteAction}, // Performed Procedure Step Start Time
	{tag.Tag{Group: 0x0040, Element: 0x0241}, DeleteAction}, // Performed Station AE Title
	{
		tag.Tag{Group: 0x0040, Element: 0x4030},
		DeleteAction,
	}, // Performed Station Geographic Location Code Sequence
	{tag.Tag{Group: 0x0040, Element: 0x0242}, DeleteAction}, // Performed Station Name
	{tag.Tag{Group: 0x0040, Element: 0x4028}, DeleteAction}, // Performed Station Name Code Sequence
	{
		tag.Tag{Group: 0x0008, Element: 0x1052},
		DeleteAction,
	}, // Performing Physician Identification Sequence
	{tag.Tag{Group: 0x0008, Element: 0x1050}, DeleteAction}, // Performing Physicians' Name
	{tag.Tag{Group: 0x0040, Element: 0x1102}, DeleteAction}, // Person Address
	{tag.Tag{Group: 0x0040, Element: 0x1103}, DeleteAction}, // Person Telephone Numbers
	{tag.Tag{Group: 0x4008, Element: 0x0114}, DeleteAction}, // Physician Approving Interpretation
	{
		tag.Tag{Group: 0x0008, Element: 0x1062},
		DeleteAction,
	}, // Physician study Identification Sequence
	{tag.Tag{Group: 0x0008, Element: 0x1048}, DeleteAction}, // Physician record
	{
		tag.Tag{Group: 0x0008, Element: 0x1049},
		DeleteAction,
	}, // Physician record Identification Sequence
	{tag.Tag{Group: 0x0018, Element: 0x1004}, DeleteAction}, // Plate ID
	{tag.Tag{Group: 0x0040, Element: 0x0012}, DeleteAction}, // Pre-Medication
	{tag.Tag{Group: 0x0010, Element: 0x21C0}, DeleteAction}, // Pregnancy Status
	{
		tag.Tag{Group: 0x0040, Element: 0x2001},
		DeleteAction,
	}, // Reason for the Imaging Service Request
	{tag.Tag{Group: 0x0032, Element: 0x1030}, DeleteAction}, // Reason for Study
	{
		tag.Tag{Group: 0x0400, Element: 0x0402},
		DeleteAction,
	}, // Referenced Digital Signature Sequence
	{tag.Tag{Group: 0x0038, Element: 0x0004}, DeleteAction}, // Referenced Patient Alias Sequence
	{tag.Tag{Group: 0x0008, Element: 0x1120}, DeleteAction}, // Referenced Patient Sequence
	{tag.Tag{Group: 0x0400, Element: 0x0403}, DeleteAction}, // Referenced SOP Instance MAC Sequence
	{tag.Tag{Group: 0x0008, Element: 0x0092}, DeleteAction}, // Referring Physician's Address
	{
		tag.Tag{Group: 0x0008, Element: 0x0096},
		DeleteAction,
	}, // Referring Physician's Identification Sequence
	{
		tag.Tag{Group: 0x0008, Element: 0x0094},
		DeleteAction,
	}, // Referring Physician's Telephone Numbers
	{tag.Tag{Group: 0x0010, Element: 0x2152}, DeleteAction}, // Region of Residence
	{tag.Tag{Group: 0x0040, Element: 0x0275}, DeleteAction}, // Request Attributes Sequence
	{tag.Tag{Group: 0x0032, Element: 0x1070}, DeleteAction}, // Requested Contrast Agent
	{tag.Tag{Group: 0x0040, Element: 0x1400}, DeleteAction}, // Requested Procedure Comments
	{tag.Tag{Group: 0x0040, Element: 0x1001}, DeleteAction}, // Requested Procedure ID
	{tag.Tag{Group: 0x0040, Element: 0x1005}, DeleteAction}, // Requested Procedure Location
	{tag.Tag{Group: 0x0032, Element: 0x1032}, DeleteAction}, // Requesting Physician
	{tag.Tag{Group: 0x0032, Element: 0x1033}, DeleteAction}, // Requesting Service
	{tag.Tag{Group: 0x0010, Element: 0x2299}, DeleteAction}, // Responsible Organization
	{tag.Tag{Group: 0x0010, Element: 0x2297}, DeleteAction}, // Responsible Person
	{tag.Tag{Group: 0x4008, Element: 0x4000}, DeleteAction}, // Results Comments
	{tag.Tag{Group: 0x4008, Element: 0x0118}, DeleteAction}, // Results Distribution List Sequence
	{tag.Tag{Group: 0x4008, Element: 0x0042}, DeleteAction}, // Results ID Issuer
	{tag.Tag{Group: 0x0040, Element: 0x4034}, DeleteAction}, // Scheduled Human Performers Sequence
	{
		tag.Tag{Group: 0x0038, Element: 0x001E},
		DeleteAction,
	}, // Scheduled Patient Institution Residence
	{
		tag.Tag{Group: 0x0040, Element: 0x000B},
		DeleteAction,
	}, // Scheduled Performing Physician Identification Sequence
	{tag.Tag{Group: 0x0040, Element: 0x0006}, DeleteAction}, // Scheduled Performing Physician Name
	{tag.Tag{Group: 0x0040, Element: 0x0004}, DeleteAction}, // Scheduled Procedure Step End Date
	{tag.Tag{Group: 0x0040, Element: 0x0005}, DeleteAction}, // Scheduled Procedure Step End Time
	{tag.Tag{Group: 0x0040, Element: 0x0007}, DeleteAction}, // Scheduled Procedure Step Description
	{tag.Tag{Group: 0x0040, Element: 0x0011}, DeleteAction}, // Scheduled Procedure Step Location
	{tag.Tag{Group: 0x0040, Element: 0x0002}, DeleteAction}, // Scheduled Procedure Step Start Date
	{tag.Tag{Group: 0x0040, Element: 0x0003}, DeleteAction}, // Scheduled Procedure Step Start Time
	{tag.Tag{Group: 0x0040, Element: 0x0001}, DeleteAction}, // Scheduled Station AE Title
	{
		tag.Tag{Group: 0x0040, Element: 0x4027},
		DeleteAction,
	}, // Scheduled Station Geographic Location Code Sequence
	{tag.Tag{Group: 0x0040, Element: 0x0010}, DeleteAction}, // Scheduled Station Name
	{tag.Tag{Group: 0x0040, Element: 0x4025}, DeleteAction}, // Scheduled Station Name Code Sequence
	{tag.Tag{Group: 0x0032, Element: 0x1020}, DeleteAction}, // Scheduled Study Location
	{tag.Tag{Group: 0x0032, Element: 0x1021}, DeleteAction}, // Scheduled Study Location AE Title
	{tag.Tag{Group: 0x0008, Element: 0x103E}, DeleteAction}, // Series Description
	{tag.Tag{Group: 0x0038, Element: 0x0062}, DeleteAction}, // Service Episode Description
	{tag.Tag{Group: 0x0038, Element: 0x0060}, DeleteAction}, // Service Episode ID
	{tag.Tag{Group: 0x0010, Element: 0x21A0}, DeleteAction}, // Smoking Status
	{tag.Tag{Group: 0x0038, Element: 0x0050}, DeleteAction}, // Special Needs
	{tag.Tag{Group: 0x0032, Element: 0x4000}, DeleteAction}, // Study Comments
	{tag.Tag{Group: 0x0008, Element: 0x1030}, DeleteAction}, // Study Description
	{tag.Tag{Group: 0x0032, Element: 0x0012}, DeleteAction}, // Study ID Issuer
	{tag.Tag{Group: 0x4000, Element: 0x4000}, DeleteAction}, // Text Comments
	{tag.Tag{Group: 0x2030, Element: 0x0020}, DeleteAction}, // Text String
	{tag.Tag{Group: 0x0008, Element: 0x0201}, DeleteAction}, // Timezone Offset From UTC
	{tag.Tag{Group: 0x0088, Element: 0x0910}, DeleteAction}, // Topic Author
	{tag.Tag{Group: 0x0088, Element: 0x0912}, DeleteAction}, // Topic Keywords
	{tag.Tag{Group: 0x0088, Element: 0x0906}, DeleteAction}, // Topic Subject
	{tag.Tag{Group: 0x0088, Element: 0x0904}, DeleteAction}, // Topic Title
	{tag.Tag{Group: 0x0040, Element: 0xA027}, DeleteAction}, // Verifying Organization
	{tag.Tag{Group: 0x0038, Element: 0x4000}, DeleteAction}, // Visit Comments
}

// Tags's value to be replaced with fake UID
var U = []TBA{
	{tag.Tag{Group: 0x0020, Element: 0x9161}, ReplaceUIDAction}, // Concatenation UID
	{
		tag.Tag{Group: 0x0008, Element: 0x010D},
		ReplaceUIDAction,
	}, // Context Group Extension Creator UID
	{tag.Tag{Group: 0x0008, Element: 0x9123}, ReplaceUIDAction}, // Creator Version UID
	{tag.Tag{Group: 0x0018, Element: 0x1002}, ReplaceUIDAction}, // Device UID
	{tag.Tag{Group: 0x0020, Element: 0x9164}, ReplaceUIDAction}, // Dimension Organization UID
	{tag.Tag{Group: 0x300A, Element: 0x0013}, ReplaceUIDAction}, // Dose Reference UID
	{tag.Tag{Group: 0x0008, Element: 0x0058}, ReplaceUIDAction}, // Failed SOP Instance UID List
	{tag.Tag{Group: 0x0070, Element: 0x031A}, ReplaceUIDAction}, // Fiducial UID
	{tag.Tag{Group: 0x0020, Element: 0x0052}, ReplaceUIDAction}, // Frame of Reference UID
	{tag.Tag{Group: 0x0008, Element: 0x0014}, ReplaceUIDAction}, // Instance Creator UID
	{tag.Tag{Group: 0x0008, Element: 0x3010}, ReplaceUIDAction}, // Irradiation Event UID
	{
		tag.Tag{Group: 0x0028, Element: 0x1214},
		ReplaceUIDAction,
	}, // Large Palette Color Lookup Table UID
	{tag.Tag{Group: 0x0002, Element: 0x0003}, ReplaceUIDAction}, // Media Storage SOP Instance UID
	{tag.Tag{Group: 0x0028, Element: 0x1199}, ReplaceUIDAction}, // Palette Color Lookup Table UID
	{
		tag.Tag{Group: 0x3006, Element: 0x0024},
		ReplaceUIDAction,
	}, // Referenced Frame of Reference UID
	{
		tag.Tag{Group: 0x0040, Element: 0x4023},
		ReplaceUIDAction,
	}, // Referenced General Purpose Scheduled Procedure Step Transaction UID
	{tag.Tag{Group: 0x0008, Element: 0x1155}, ReplaceUIDAction}, // Referenced SOP Instance UID
	{
		tag.Tag{Group: 0x0004, Element: 0x1511},
		ReplaceUIDAction,
	}, // Referenced SOP Instance UID in File
	{tag.Tag{Group: 0x3006, Element: 0x00C2}, ReplaceUIDAction}, // Related Frame of Reference UID
	{tag.Tag{Group: 0x0000, Element: 0x1001}, ReplaceUIDAction}, // Requested SOP Instance UID
	{tag.Tag{Group: 0x0020, Element: 0x000E}, ReplaceUIDAction}, // Series Instance UID
	{tag.Tag{Group: 0x0008, Element: 0x0018}, ReplaceUIDAction}, // SOP Instance UID
	{tag.Tag{Group: 0x0088, Element: 0x0140}, ReplaceUIDAction}, // Storage Media File-set UID
	{tag.Tag{Group: 0x0020, Element: 0x000D}, ReplaceUIDAction}, // Study Instance UID
	{
		tag.Tag{Group: 0x0020, Element: 0x0200},
		ReplaceUIDAction,
	}, // Synchronization Frame of Reference UID
	{tag.Tag{Group: 0x0040, Element: 0xDB0D}, ReplaceUIDAction}, // Template Extension Creator UID
	{
		tag.Tag{Group: 0x0040, Element: 0xDB0C},
		ReplaceUIDAction,
	}, // Template Extension Organization UID
	{tag.TransactionUID, ReplaceUIDAction},                      // Transaction UID
	{tag.Tag{Group: 0x0040, Element: 0xA124}, ReplaceUIDAction}, // UID
}

// Tags to be deleted or value emptied
var ZD = []TBA{
	{tag.Tag{Group: 0x0008, Element: 0x0023}, EmptyOrReplaceAction}, // Content Date
	{tag.Tag{Group: 0x0008, Element: 0x0033}, EmptyOrReplaceAction}, // Content Time
	{tag.Tag{Group: 0x0018, Element: 0x0010}, EmptyOrReplaceAction}, // Contrast Bolus Agent
}

// Tags to be deleted or value emptied or replaced
var XZ = []TBA{
	{tag.Tag{Group: 0x0008, Element: 0x0022}, DeleteOrEmptyAction}, // Acquisition Date
	{tag.Tag{Group: 0x0008, Element: 0x0032}, DeleteOrEmptyAction}, // Acquisition Time
	{tag.Tag{Group: 0x0010, Element: 0x2203}, DeleteOrEmptyAction}, // Patient Sex Neutered
	{tag.Tag{Group: 0x0008, Element: 0x1110}, DeleteOrEmptyAction}, // Referenced Study Sequence
	{
		tag.Tag{Group: 0x0032, Element: 0x1060},
		DeleteOrEmptyAction,
	}, // Requested Procedure Description
	{tag.Tag{Group: 0x300E, Element: 0x0008}, DeleteOrEmptyAction}, // Reviewer Name
}

// Tags to be deleted or value emptied or replaced
var XD = []TBA{
	{tag.Tag{Group: 0x0008, Element: 0x002A}, DeleteOrReplaceAction}, // Acquisition DateTime
	{
		tag.Tag{Group: 0x0018, Element: 0x1400},
		DeleteOrReplaceAction,
	}, // Acquisition Device Processing Description
	{tag.Tag{Group: 0x0018, Element: 0x700A}, DeleteOrReplaceAction}, // Detector ID
	{
		tag.Tag{Group: 0x0008, Element: 0x1072},
		DeleteOrReplaceAction,
	}, // Operators' Identification Sequence
	{tag.Tag{Group: 0x0018, Element: 0x1030}, DeleteOrReplaceAction}, // Protocol Name
	{tag.Tag{Group: 0x0008, Element: 0x0021}, DeleteOrReplaceAction}, // Series Date
	{tag.Tag{Group: 0x0008, Element: 0x0031}, DeleteOrReplaceAction}, // Series Time
}

// Tags to be deleted or value emptied or replaced
var XZD = []TBA{
	{tag.Tag{Group: 0x0018, Element: 0x1000}, DeleteOrEmptyOrReplaceAction}, // Device Serial Number
	{
		tag.Tag{Group: 0x0008, Element: 0x0082},
		DeleteOrEmptyOrReplaceAction,
	}, // Institution Code Sequence
	{tag.Tag{Group: 0x0008, Element: 0x0080}, DeleteOrEmptyOrReplaceAction}, // Institution Name
	{tag.Tag{Group: 0x0008, Element: 0x1070}, DeleteOrEmptyOrReplaceAction}, // Operators' Name
	{
		tag.Tag{Group: 0x0008, Element: 0x1111},
		DeleteOrEmptyOrReplaceAction,
	}, // Referenced Performed Procedure Step Sequence
	{tag.Tag{Group: 0x0008, Element: 0x1010}, DeleteOrEmptyOrReplaceAction}, // Station Name
}

// To be cleaned with UI as VR.
// Or cleaned according to VR (XZU*)
var XZUStar = []TBA{
	{
		tag.Tag{Group: 0x0008, Element: 0x1140},
		DeleteOrEmptyOrReplaceUIDAction,
	}, // Referenced Image Sequence
	{
		tag.Tag{Group: 0x0008, Element: 0x2112},
		DeleteOrEmptyOrReplaceUIDAction,
	}, // Source Image Sequence
}

var tags = append([]TBA{}, D...)

func init() {
	tags = append(tags, Z...)
	tags = append(tags, X...)
	tags = append(tags, U...)
	tags = append(tags, ZD...)
	tags = append(tags, XZ...)
	tags = append(tags, XD...)
	tags = append(tags, XZD...)
	tags = append(tags, XZUStar...)
}
