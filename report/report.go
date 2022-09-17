package report

import (
	"encoding/xml"
)

type (

	// ZACK is a Z-Report acknowledgement message.
	//ZNUMBER Int
	//DATE DateAndTime
	//TIME DateAndTime
	//ACKCODE Int(1)
	//ACKMSG String(50)
	//Mandatory Mandatory Mandatory Mandatory Mandatory
	//Z-Report Acknowledgement Envelop Z-Report Number for the posted receipt Report Date YYYY-MM-DD
	//Report Time HH24:MI:SS
	//0 means success. Else it would be an error code Describes the ACKCODE above
	ZACK struct {
		XMLName xml.Name `xml:"ZACK"`
		Text    string   `xml:",chardata"`
		ZNUMBER int64    `xml:"ZNUMBER"`
		DATE    string   `xml:"DATE"`
		TIME    string   `xml:"TIME"`
		ACKCODE int64    `xml:"ACKCODE"`
		ACKMSG  string   `xml:"ACKMSG"`
	}

	ZACKEFDMS struct {
		XMLName        xml.Name `xml:"EFDMS"`
		Text           string   `xml:",chardata"`
		ZACK           ZACK     `xml:"ZACK"`
		EFDMSSIGNATURE string   `xml:"EFDMSSIGNATURE"`
	}
	EFDMS struct {
		XMLName        xml.Name `xml:"EFDMS"`
		Text           string   `xml:",chardata"`
		ZREPORT        ZREPORT  `xml:"ZREPORT"`
		EFDMSSIGNATURE string   `xml:"EFDMSSIGNATURE"`
	}

	Header struct {
		XMLName xml.Name `xml:"HEADER"`
		Text    string   `xml:",chardata"`
		ZREPORT ZREPORT  `xml:"ZREPORT"`
	}

	TOTALS struct {
		XMLName          xml.Name `xml:"TOTALS"`
		Text             string   `xml:",chardata"`
		DAILYTOTALAMOUNT float64  `xml:"DAILYTOTALAMOUNT"`
		GROSS            float64  `xml:"GROSS"`
		CORRECTIONS      float64  `xml:"CORRECTIONS"`
		DISCOUNTS        float64  `xml:"DISCOUNTS"`
		SURCHARGES       float64  `xml:"SURCHARGES"`
		TICKETSVOID      int64    `xml:"TICKETSVOID"`
		TICKETSVOIDTOTAL int64    `xml:"TICKETSVOIDTOTAL"`
		TICKETSFISCAL    int64    `xml:"TICKETSFISCAL"`
		TICKETSNONFISCAL int64    `xml:"TICKETSNONFISCAL"`
	}

	VATTOTALS struct {
		XMLName    xml.Name `xml:"VATTOTALS"`
		Text       string   `xml:",chardata"`
		VATRATE    string   `xml:"VATRATE"`
		NETTAMOUNT float64  `xml:"NETTAMOUNT"`
		TAXAMOUNT  float64  `xml:"TAXAMOUNT"`
	}

	Payment struct {
		XMLName   xml.Name `xml:"PAYMENTS"`
		Text      string   `xml:",chardata"`
		PMTTYPE   string   `xml:"PMTTYPE"`
		PMTAMOUNT float64  `xml:"PMTAMOUNT"`
	}

	ZREPORT struct {
		XMLName xml.Name `xml:"ZREPORT"`
		Text    string   `xml:",chardata"`
		DATE    string   `xml:"DATE"`
		TIME    string   `xml:"TIME"`
		HEADER  struct {
			Text string   `xml:",chardata"`
			LINE []string `xml:"LINE"`
		} `xml:"HEADER"`
		VRN              string    `xml:"VRN"`
		TIN              string    `xml:"TIN"`
		TAXOFFICE        string    `xml:"TAXOFFICE"`
		REGID            string    `xml:"REGID"`
		ZNUMBER          string    `xml:"ZNUMBER"`
		EFDSERIAL        string    `xml:"EFDSERIAL"`
		REGISTRATIONDATE string    `xml:"REGISTRATIONDATE"`
		USER             string    `xml:"USER"`
		SIMIMSI          string    `xml:"SIMIMSI"`
		TOTALS           TOTALS    `xml:"TOTALS"`
		VATTOTALS        VATTOTALS `xml:"VATTOTALS"`
		PAYMENTS         struct {
			Text      string   `xml:",chardata"`
			PMTTYPE   []string `xml:"PMTTYPE"`
			PMTAMOUNT []string `xml:"PMTAMOUNT"`
		} `xml:"PAYMENTS"`
		CHANGES struct {
			Text          string `xml:",chardata"`
			VATCHANGENUM  string `xml:"VATCHANGENUM"`
			HEADCHANGENUM string `xml:"HEADCHANGENUM"`
		} `xml:"CHANGES"`
		ERRORS     string `xml:"ERRORS"`
		FWVERSION  string `xml:"FWVERSION"`
		FWCHECKSUM string `xml:"FWCHECKSUM"`
	}
)
