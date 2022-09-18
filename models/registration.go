package models

import "encoding/xml"

type (
	RegistrationAck struct {
		XMLName        xml.Name             `xml:"EFDMS"`
		Text           string               `xml:",chardata"`
		EFDMSRESP      RegistrationResponse `xml:"EFDMSRESP"`
		EFDMSSIGNATURE string               `xml:"EFDMSSIGNATURE"`
	}
)

// RegistrationResponse is the response message received from the VFD
// after a successful registration.
type RegistrationResponse struct {
	XMLName     xml.Name `xml:"EFDMSRESP"`
	Text        string   `xml:",chardata"`
	ACKCODE     string   `xml:"ACKCODE"`
	ACKMSG      string   `xml:"ACKMSG"`
	REGID       string   `xml:"REGID"`
	SERIAL      string   `xml:"SERIAL"`
	UIN         string   `xml:"UIN"`
	TIN         string   `xml:"TIN"`
	VRN         string   `xml:"VRN"`
	MOBILE      string   `xml:"MOBILE"`
	ADDRESS     string   `xml:"ADDRESS"`
	STREET      string   `xml:"STREET"`
	CITY        string   `xml:"CITY"`
	COUNTRY     string   `xml:"COUNTRY"`
	NAME        string   `xml:"NAME"`
	RECEIPTCODE string   `xml:"RECEIPTCODE"`
	REGION      string   `xml:"REGION"`
	ROUTINGKEY  string   `xml:"ROUTINGKEY"`
	GC          string   `xml:"GC"`
	TAXOFFICE   string   `xml:"TAXOFFICE"`
	USERNAME    string   `xml:"USERNAME"`
	PASSWORD    string   `xml:"PASSWORD"`
	TOKENPATH   string   `xml:"TOKENPATH"`
	TAXCODES    TaxCodes `xml:"TAXCODES"`
}

type TaxCodes struct {
	XMLName xml.Name `xml:"TAXCODES"`
	Text    string   `xml:",chardata"`
	CODEA   string   `xml:"CODEA"`
	CODEB   string   `xml:"CODEB"`
	CODEC   string   `xml:"CODEC"`
	CODED   string   `xml:"CODED"`
}

type RegistrationRequest struct {
	XMLName        xml.Name         `xml:"EFDMS"`
	Text           string           `xml:",chardata"`
	Reg            RegistrationBody `xml:"REGDATA"`
	EFDMSSIGNATURE string           `xml:"EFDMSSIGNATURE"`
}

type RegistrationBody struct {
	XMLName xml.Name `xml:"REGDATA"`
	Text    string   `xml:",chardata"`
	TIN     string   `xml:"TIN"`
	CERTKEY string   `xml:"CERTKEY"`
}

type Error struct {
	XMLName xml.Name `xml:"Error"`
	Text    string   `xml:",chardata"`
	Message string   `xml:"Message"`
}

func (regBody *RegistrationBody) Validate() error {
	return nil
}
