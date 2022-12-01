package vfd

// ACKCODE	STATUS	DESCRIPTION	POSSIBLE REASON
// 0	SUCCESS	Success
// 1	FAIL	Invalid Signature	Signature generated not in correct format. Signature generated with missing nodes, signature generated with empty lines in XML or
// 3	FAIL	Invalid TIN	TIN specified with dash or wrong TIN specified
// 4	FAIL	VFD Registration Approval required	Request posted without Client details, which is WEBAPI.
// 5	FAIL	Unhandled Exception	Contact TRA for further troubleshooting
// 6		Invalid Serial or Serial not Registered to Web API/TIN	CERTKEY is not registered to TIN sending registration request. Use only TIN and CERTKEY provided by TRA
// 7	FAIL	Invalid client header	Wrong client value specified
// 8	FAIL	Wrong Certificate used to Register Web API	Wrong certificate used

const (
	SuccessCode          int64 = 0
	InvalidSignatureCode int64 = 1
	InvalidTaxID         int64 = 3
	ApprovalRequired     int64 = 4
	UnhandledException   int64 = 5
	InvalidSerial        int64 = 6
	InvalidClientHeader  int64 = 7
	InvalidCertificate   int64 = 8
)

type Error struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ParseErrorCode parses the error code and returns the corresponding error message.
func ParseErrorCode(code int64) string {
	switch code {
	case 0:
		return "SUCCESS"
	case 1:
		return "FAIL"
	case 3:
		return "Invalid TIN"
	case 4:
		return "VFD Registration Approval required"
	case 5:
		return "Unhandled Exception"
	case 6:
		return "Invalid Serial or Serial not Registered to Web API/TIN"
	case 7:
		return "Invalid client header"
	case 8:
		return "Wrong Certificate used to Register Web API"
	default:
		return "Unknown error"
	}
}
