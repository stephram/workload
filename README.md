# workload

### messages

Command line utility for sending messages to the S2C Input Queue.

#### Building

```
make install
make build
```
Build output will be similar to this.
```
CGO_ENABLED=0 GOARCH=amd64 \
	go build -a \
	-installsuffix cgo \
	-ldflags='-w -s -X "github.com/stephram/workload/pkg/app.Name=messages" -X "github.com/stephram/workload/pkg/app.Product=messages" -X "github.com/stephram/workload/pkg/app.Branch=master" -X "github.com/stephram/workload/pkg/app.BuildDate=2020-04-02 14:00:40" -X "github.com/stephram/workload/pkg/app.Commit=86e2b8670a49ea74500c393076358254d4f2dcf4" -X "github.com/stephram/workload/pkg/app.Version=0.0.1" -X "github.com/stephram/workload/pkg/app.Author=Stephen Graham" -X "github.com/stephram/workload/pkg/app.AuthorEmail=stephram@gmail.com" ' \
	-o /Users/sg/Development/api/s2c/workload/build/messages cmd/message_generator/messages.go
```
#### Tests
```
$ make test
Starting tests
```
Test run output
```
go test ./... -v -cover -coverprofile=coverage.out
=== RUN   TestParseStoreNumbers
=== RUN   TestParseStoreNumbers/Parse_store_numbers_string
--- PASS: TestParseStoreNumbers (0.00s)
    --- PASS: TestParseStoreNumbers/Parse_store_numbers_string (0.00s)
PASS
coverage: 17.9% of statements
ok  	workload/cmd/bot_controller	0.024s	coverage: 17.9% of statements
=== RUN   TestCreateSendMessageIinput
=== RUN   TestCreateSendMessageIinput/Create_messages_from_test_files
--- PASS: TestCreateSendMessageIinput (0.00s)
    --- PASS: TestCreateSendMessageIinput/Create_messages_from_test_files (0.00s)
PASS
coverage: 40.0% of statements
ok  	workload/cmd/message_generator	0.020s	coverage: 40.0% of statements
=== RUN   TestGetLogger
=== RUN   TestGetLogger/check_that_logger_is_initialised
--- PASS: TestGetLogger (0.00s)
    --- PASS: TestGetLogger/check_that_logger_is_initialised (0.00s)
=== RUN   TestParseCommaSeparatedFiles
=== RUN   TestParseCommaSeparatedFiles/test_for_correct_split
--- PASS: TestParseCommaSeparatedFiles (0.00s)
    --- PASS: TestParseCommaSeparatedFiles/test_for_correct_split (0.00s)
=== RUN   TestParseCommaSeparatedStrings
=== RUN   TestParseCommaSeparatedStrings/test_separation
=== RUN   TestParseCommaSeparatedStrings/test_empty_element
--- PASS: TestParseCommaSeparatedStrings (0.00s)
    --- PASS: TestParseCommaSeparatedStrings/test_separation (0.00s)
    --- PASS: TestParseCommaSeparatedStrings/test_empty_element (0.00s)
=== RUN   TestSelectRandomString
=== RUN   TestSelectRandomString/select_from_slice
=== RUN   TestSelectRandomString/select_from_empty_slice
--- PASS: TestSelectRandomString (0.00s)
    --- PASS: TestSelectRandomString/select_from_slice (0.00s)
    --- PASS: TestSelectRandomString/select_from_empty_slice (0.00s)
PASS
coverage: 88.5% of statements
ok  	workload/internal/utils	0.026s	coverage: 88.5% of statements
=== RUN   TestUlidPackage
=== RUN   TestUlidPackage/Valid_ULID
=== RUN   TestUlidPackage/Valid_Monotonic_time
--- PASS: TestUlidPackage (0.00s)
    --- PASS: TestUlidPackage/Valid_ULID (0.00s)
    --- PASS: TestUlidPackage/Valid_Monotonic_time (0.00s)
PASS
coverage: 100.0% of statements
ok  	workload/internal/utils/ulid	0.028s	coverage: 100.0% of statements
Tests complete
go tool cover -html=coverage.out -o ./coverage.html
Coverage report written to coverage.html
```
#### Usage:

```
$ ./messages --help
usage: messages [<flags>]

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -n, --number-of-messages=5   Number of messages to send
  -q, --queue-url=https://sqs.ap-southeast-2.amazonaws.com/712510509017/api-dev-s2c-inbound
                               URL of the SQS queue
  -t, --message-template=test-data/sales-messages/1808712-body.json... ...
                               Filenames containing payloads to use as templates
  -s, --store-number=A399... ...
                               Store numbers use
  -k, --key-start-id=333000    Reference key start index
  -w, --number-of-workers=100  Number of workers
  -p, --profile="api-dev"      AWS profile name
```
#### Running
```
$ ./messages -s A399 -s A301 -n 10 -t test-data/sales-messages/1808712-body.json -t test-data/sales-messages/1808713-body.json
  number-of-messages : 10
   number-of-workers : 10
   message-templates : (2) [test-data/sales-messages/1808712-body.json test-data/sales-messages/1808713-body.json]
       store-numbers : (2) [A399 A301]
        key-start-id : 333000
           queue-url : https://sqs.ap-southeast-2.amazonaws.com/712510509017/api-dev-s2c-inbound
Press 'Enter' to continue...
{"level":"info","msg":"Creating worker 10 tasks","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"Send 10 messages","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"Waiting for output","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: 3e2ea685-72c1-44a6-acc7-95b90e6cfaa2, MessageHeader: 6da35958-10fb-4578-86bc-49cd0bf48058, STORE_REF: A301, TRS_KEY: 333005, SEQUENCE_NUMBER: 6","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: bad1d5fe-80eb-4082-ba63-790d67596d20, MessageHeader: ed8b4a07-8f50-46ef-9bd5-448892ee7205, STORE_REF: A301, TRS_KEY: 333008, SEQUENCE_NUMBER: 9","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: 5d593a4c-1028-498c-867e-9fac54b87f4f, MessageHeader: 6c0a43b2-fe12-4cdd-a3b1-ab15df2f3e04, STORE_REF: A301, TRS_KEY: 333006, SEQUENCE_NUMBER: 7","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: f324cf64-ba49-4d4e-b9c2-269596e5f358, MessageHeader: 2f9735b4-dfee-4a2e-a380-db1061573f73, STORE_REF: A399, TRS_KEY: 333003, SEQUENCE_NUMBER: 4","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: 47660a36-0cd7-4249-a4da-b0f782157072, MessageHeader: c933921c-5ec6-450a-bf73-695f2078387f, STORE_REF: A399, TRS_KEY: 333002, SEQUENCE_NUMBER: 3","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: a8571336-74d3-42e9-9ad2-64ef96591acc, MessageHeader: 8501b832-ff89-46dc-8f27-562a88b1e0c8, STORE_REF: A301, TRS_KEY: 333004, SEQUENCE_NUMBER: 5","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: 2ba7f5fc-7b89-4e56-9cfe-9faf510519bd, MessageHeader: 2264b888-0f77-4706-9730-23051179232b, STORE_REF: A399, TRS_KEY: 333009, SEQUENCE_NUMBER: 10","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: e5b10782-a67b-467b-a91d-d5ff9985d3ce, MessageHeader: 240f876c-7b27-4deb-9734-55287381dd0f, STORE_REF: A301, TRS_KEY: 333007, SEQUENCE_NUMBER: 8","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: b3300ed6-01de-4786-b34c-b400c07e8185, MessageHeader: 65370761-a978-49b1-bc4d-6da8610587dd, STORE_REF: A301, TRS_KEY: 333000, SEQUENCE_NUMBER: 1","time":"2020-04-13T14:58:51+10:00"}
{"level":"info","msg":"sent MessageId: 62a79112-4510-43a3-9fce-070a0d0a798e, MessageHeader: f2fd7416-8dcf-4b09-ac2e-bc70292a446c, STORE_REF: A399, TRS_KEY: 333001, SEQUENCE_NUMBER: 2","time":"2020-04-13T14:58:51+10:00"}```