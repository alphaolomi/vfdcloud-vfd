## INTEGRATION RULES
1. Registration is a onetime request. In order to be able to access EFDMS taxpayer (seller)
   must send TIN number so TRA can register VFD for use
2. Registration data must be saved to taxpayer system database for later use.
3. When successful registered VFD will not be required to send registration again when
   submitting receipts/invoices.
4. TIN belongs to the seller and not the customer (buyer). There is a parameter for customer
   i.e. CUSTID
5. Cert-Serial is a serial number of certificate file (pfx) used to sign request where private key
   will be used in this case. (TRA provides)
6. GC, RCTNUM and Dc must be maintained by taxpayer's system
7. GC must always be equal to RCNUM and starts from one (1) and always increment for
   each invoice/receipt issued. Numbers should follow sequence without skipping
8. DC starts with 1 and increment until 23:59:59 (midnight) and after midnight DC resets
   (starts with 1 again when first invoice/receipt for a new day is issued)
9. Different receipts/invoices can never have same GC they must always be unique
10. Different receipts/invoices can never have same RCTNUM they must always be unique
11. Different receipts/invoices issued in the same day can never have same DC
12. ZNUM must always be same as RCT_DATE but in a number format i.e. YYYYMMDD
13. Current receipt/invoice can never have old ZNUM, RCT_DATE than previous
    receipt/invoice
14. Future receipts/invoices can never have old ZNUM, RCT_DATE
15. Future dates are not allowed, therefore when VFD generates receipts/invoices it must
    synchronize with NTP server to take current timestamp
16. If transaction is cancelled, next transaction should take not reuse transaction of cancelled
    transaction it should take new number in the sequence
17. Token will be requested only after expiry of current one, so before posting to TRA check
    if current token is valid and only post if is valid otherwise request new one
18. Token value must be saved to taxpayer system database
19. If for some reason if receipt is issued with wrong data but received success response,
    when sending again with correct data assign a new receipt number don’t reuse same if
    you send same receipt the later will be treated as duplicate and wont saved because
    rctvnum is same.
20. If CUSTIDTYPE=1 i.e. TIN is chosen, we recommend to restrict input to only 9 digits
    meaning only numbers should be allowed because TIN is always a 9 digits number.
    CUSTIDTYPES for other IDs can remain open as string.
21. MOBILENUM should not contain + or spaces or dashes, it should in the format
    255712XXXXXX or 0712XXXXXX
22. If VFD get success response it should not resubmit same transaction again.
23. If for any reason VFD does not receive response at all or receiving negative response
    (ACKCODE which is not 0) for specific invoice/receipt then when resubmitting the same
    invoice/receipt to TRA the VFD should submit the original xml content and not the modified
    the content this include also ZNUM and RCT_DATE, RCT_TIME must always be date of
    the first attempt and not the current date/time. This is to say monitor status of each receipt
    and only when response with ACKCODE 0 returned consider receipt successful delivered
    to TRA
24. Print/send receipt/invoice to customer (do not wait for TRA response) and immediately
    send receipt to TRA (1 and 2 can either be concurrent or 2 can follow after 1)
25. For a specific receipt if no response is received VFD should keep try sending same request
    until it receives response.
26. Send one transaction at a time only send next transaction when current one has
    succeeded
27. For printed receipts/invoices, taxpayer must display verification information (QR and code)
    on the printout.
28. To avoid receipt/invoices being rejected ,escape special characters in receipt/invoice XML
    especially in customer name or items descriptions.
29. When TRA server is not accessible (OFFLINE), continue generating transactions as they
    occur but make sure you design a mechanism to save status of each transaction i.e.
    success or pending while keep checking for connection and later when TRA connection
    resumes automatically resend all pending transactions in the order.
30. Token will be requested only after expiry of current one, so before posting to TRA check
    if current token is valid and only post if is valid otherwise request new one
31. Token value must be saved to taxpayer billing system


## TEST CASES SCENARIOS
It is suggested that these are done during testing the VFD integration to determine if the implementation works as expected

1. Post receipts/invoices with different CUSTIDTYPES as indicated in the API
2. Post receipts/invoices with discounts if any
3. Post receipts/invoices with multiple items having different tax codes if any as per API
4. For receipt post transactions with different payment types as indicated in the API
5. We recommend posting transactions in daily basis (this helps us checking sequence of
   DC, GC, ZNUM).
6. Post as many transactions as possible preferably from 100 and above.