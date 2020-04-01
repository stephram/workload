# Messages receieved on the input queue from S2CTransport v5.0.0.1

Each file name represents a TRS_KEY and the seperate 'body' file is it's extracted body. 
The 'body' file for each is the value of the 'body' field of the incoming message extracted into a separate file.
It is the payload that is sent from S2CTransport.

## 1808713.json and 1808713-body.json

This is a prescription sale. Presecription sales cant be identifed by the presence of the ```TLog60D601ScriptInf1``` entry.

Example:
```json
        "IncludedTables": [
            "TLog010000ItemSale",
            "TLog60D601ScriptInf1",
            "TLog040000Tender",
            "TLog050000TicketEnd",
            "TLog210000TicketStart",
            "TLog_Header"
        ]
```

An non prescription type of `sales` transaction would not have this TLog entry.

Example:
```json
        "IncludedTables": [
            "TLog010000ItemSale",
            "TLog030000Disc",
            "TLog040000Tender",
            "TLog050000TicketEnd",
            "TLog110000ItemSaleExt1",
            "TLog210000TicketStart",
            "TLog_Header"
        ]
```
