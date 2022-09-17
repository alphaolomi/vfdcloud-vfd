package internal

import (
	"testing"
)

func TestProduceEFDMS(t *testing.T) {
	// EFDMS Example:
	//<?xml version="1.0"?>
	//<EFDMS>
	//	<RCT>
	//		<DATE>2019-08-27</DATE>
	//		<!-- Date of issue of receipt/invoice in the format YYYY-MM-DD-->
	//		<TIME>08:36:02</TIME>
	//		<!-- Time of issue of receipt/invoice in the format HH:MM:SS-->
	//		<TIN>222222286</TIN>
	//		<!-- Tin of the taxpayer -->
	//		<REGID>TZ0100089</REGID>
	//		<!-- Registration number of the VFD -->
	//		<EFDSERIAL>10TZ1000211</EFDSERIAL>
	//		<!-- Serial number of VFD -->
	//		<CUSTIDTYPE>1</CUSTIDTYPE>
	//		<!-- Customer ID type, values range from 1 to 6 as specified in API document-->
	//		<CUSTID>111111111</CUSTID>
	//		<!-- Customer ID based on CUSTIDTYPE specified above-->
	//		<CUSTNAME></CUSTNAME>
	//		<!-- Customer name-->
	//		<MOBILENUM></MOBILENUM>
	//		<!-- Customer mobile number-->
	//		<RCTNUM>380</RCTNUM>
	//		<!-- A receipt/invoice number which is same as GC. It should compose of digits alone i.e. without letters-->
	//		<DC>1</DC>
	//		<!-- Daily counter of recipt/invoice which increments for each receipt/invoice and reset to 1 on a new day-->
	//		<GC>380</GC>
	//		<!-- Global counter of receipt/invoice which increment throughout the life of the VFD. It has the same value as RCTNUM-->
	//		<ZNUM>20190827</ZNUM>
	//		<!-- ZNUM will be a date of receipt/invoice generated as number in format of (YYYYMMDD) -->
	//		<RCTVNUM>MFT7AB380</RCTVNUM>
	//		<!-- receipt/invoice verification number which is in the combination of RCTVCODE and GC and is unique for each transaction. If RCTVCODE=RMH3YK and GC=36 then RCTVNUM=RMH3YK36-->
	//		<ITEMS>
	//			<ITEM>
	//				<ID>101</ID>
	//				<!-- Item ID-->
	//				<DESC>Item desc</DESC>
	//				<!-- Description of the item-->
	//				<QTY>1</QTY>
	//				<!-- Quantity-->
	//				<TAXCODE>1</TAXCODE>
	//				<!-- Tax code applicable on the item 1 for taxbale items and 3 for non-taxable items-->
	//				<AMT>200.00</AMT>
	//				<!-- Amount inclusive of TAX for each item-->
	//			</ITEM>
	//		</ITEMS>
	//		<TOTALS>
	//			<TOTALTAXEXCL>169.49</TOTALTAXEXCL>
	//			<!-- Total of all the items exclusive of Tax-->
	//			<TOTALTAXINCL>200</TOTALTAXINCL>
	//			<!-- Total of all the items inclusive of Tax-->
	//			<DISCOUNT>0.00</DISCOUNT>
	//		</TOTALS>
	//		<PAYMENTS>
	//			<PMTTYPE>CASH</PMTTYPE>
	//			<!-- Mode of Payment can either be CASH, CHEQUE, EMONEY or CCARD if receipt is generated. In this case payment is already received-->
	//			<PMTAMOUNT>200.00</PMTAMOUNT>
	//			<!-- Payment amount-->
	//			<PMTTYPE>INVOICE</PMTTYPE>
	//			<!-- Mode of Payment can only be INVOICE if invoice is generated. In this case payment is not yet received that is why we use Invoice -->
	//			<PMTAMOUNT>200.00</PMTAMOUNT>
	//			<!-- Payment amount-->
	//		</PAYMENTS>
	//		<VATTOTALS>
	//			<VATRATE>A</VATRATE>
	//			<!-- Tax group applicable on the items for VAT items should A and for no VAT items should be C-->
	//			<NETTAMOUNT>169.49</NETTAMOUNT>
	//			<!-- Total of all the items exclusive of Tax-->
	//			<TAXAMOUNT>30.51</TAXAMOUNT>
	//			<!-- Tax amount paid-->
	//		</VATTOTALS>
	//	</RCT>
	//	<EFDMSSIGNATURE>hsjahkskskkaksasd+cVF1kZ/uuyuasdausyduyaus//+uS6GVIA9+obJUdb/sjkadkskaskjakdjkjahkhksd87w7qjlasdas9+skajsakjskajs//iKG+5UOR+86VgKNdVcjWuPzOhAO6b/+uywuygdhsyaydshahsgkjdfal+/5s84kz5EUJocHzLMrI0dbALUP8AgC97ZTUIFrM/jZUSd624MD26BHrjy5KTurhpS+HJlsotIZqxyPbaw==</EFDMSSIGNATURE>
	//</EFDMS>
	// Produce EFDMS similar to the one in the above snippet.

	efdms := EFDMS{
		RCT: RCT{
			DATE:       "2019-08-27",
			TIME:       "08:36:02",
			TIN:        "222222286",
			REGID:      "TZ0100089",
			EFDSERIAL:  "10TZ1000211",
			CUSTIDTYPE: "1",
			CUSTID:     "111111111",
			CUSTNAME:   "",
			MOBILENUM:  "",
			RCTNUM:     "380",
			DC:         "1",
			GC:         "380",
			ZNUM:       "20190827",
			RCTVNUM:    "MFT7AB380",
			ITEMS: []ITEM{
				ITEM{
					ID:      "101",
					DESC:    "Item desc",
					QTY:     "1",
					TAXCODE: "1",
					AMT:     "200.00",
				},
			},
			TOTALS: []TOTAL{
				TOTAL{
					TOTALTAXEXCL: "169.49",
					TOTALTAXINCL: "200",
					DISCOUNT:     "0.00",
				},
			},
			PAYMENTS: []PAYMENT{
				PAYMENT{
					PMTTYPE:   "CASH",
					PMTAMOUNT: "200.00",
				},
				PAYMENT{
					PMTTYPE:   "INVOICE",
					PMTAMOUNT: "200.00",
				},
			},
			VATTOTALS: []VATTOTAL{
				VATTOTAL{
					VATRATE:    "A",
					NETTAMOUNT: "169.49",
					TAXAMOUNT:  "30.51",
				},
			},
		},
		EFDMSSIGNATURE: "hsjahkskskkaksasd+cVF1kZ/uuyuasdausyduyaus//+uS6GVIA9+obJUdb/sjkadkskaskjakdjkjahkhksd87w7qjlasdas9+skajsakjskajs//iKG+5UOR+86VgKNdVcjWuPzOhAO6b/+uywuygdhsyaydshahsgkjdfal+/5s84kz5EUJocHzLMrI0dbALUP8AgC97ZTUIFrM/jZUSd624MD26BHrjy5KTurhpS+HJlsotIZqxyPbaw==",
	}

	// Marshal the EFDMS struct to XML
	produceEFDMS, err := ProduceEFDMS(efdms)
	if err != nil {
		t.Errorf("Error producing EFDMS: %s", err)
		return
	}

	//print the XML
	t.Logf("%s", string(produceEFDMS))
}
