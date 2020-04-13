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
   message-templates : (2) [test-data/sales-messages/1808712-body.json test-data/sales-messages/1808713-body.json]
       store-numbers : (2) [A399 A301]
        key-start-id : 333000
           queue-url : https://sqs.ap-southeast-2.amazonaws.com/712510509017/api-dev-s2c-inbound
   number-of-workers : 100
Press 'Enter' to continue...
{"level":"info","msg":"Creating worker 100 tasks","time":"2020-04-13T13:17:49+10:00"}
{"level":"info","msg":"Send 10 messages","time":"2020-04-13T13:17:49+10:00"}
{"level":"info","msg":"Waiting for output","time":"2020-04-13T13:17:49+10:00"}
{"level":"info","msg":"sent MessageId: 75658f72-de97-4f8c-844d-98b852c0f901, MessageHeader: 2bc70631-5921-44d8-9699-3ada11099e5f, STORE_REF: A399, TRS_KEY: 333004, SEQUENCE_NUMBER: 5","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: 490e9cf4-8712-4daa-ace9-1fd164bbf528, MessageHeader: e3cc9588-f047-4ff3-b14d-0d17d5632390, STORE_REF: A399, TRS_KEY: 333005, SEQUENCE_NUMBER: 6","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: 416f82ce-b576-4818-9593-af6260564dd2, MessageHeader: 789aae4c-9554-4ed2-87a5-10fae98f117c, STORE_REF: A301, TRS_KEY: 333003, SEQUENCE_NUMBER: 4","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: b79ba857-4f2a-4de9-987e-67ffa638a1e3, MessageHeader: 0d122e61-cbe5-4688-97a6-acd8307c3056, STORE_REF: A301, TRS_KEY: 333008, SEQUENCE_NUMBER: 9","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: 5cdd8628-ffaa-46bd-a020-dcae4ae39a77, MessageHeader: 38c93be0-758a-40c9-aa19-825eea4a2d57, STORE_REF: A301, TRS_KEY: 333000, SEQUENCE_NUMBER: 1","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: 9282d20b-049c-4f70-8b85-e0d88f6d6e90, MessageHeader: 3d202fa0-e8ae-43f5-a907-9d0374bb7c16, STORE_REF: A301, TRS_KEY: 333006, SEQUENCE_NUMBER: 7","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: 379eb510-603a-4c57-b5f0-7e2b8ac5b377, MessageHeader: 59a3f856-d65d-46a4-92ae-64214e0d1acc, STORE_REF: A399, TRS_KEY: 333001, SEQUENCE_NUMBER: 2","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: f9a8db30-247f-4add-8579-1b0ca60164c3, MessageHeader: e45fb08d-a27f-4aa6-8c1e-f5d82dbb55f9, STORE_REF: A301, TRS_KEY: 333009, SEQUENCE_NUMBER: 10","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: ab9d91f5-c9fa-48b1-bc1b-9ded67a08426, MessageHeader: e9602a88-7c38-4f70-8247-9fad6f214267, STORE_REF: A399, TRS_KEY: 333007, SEQUENCE_NUMBER: 8","time":"2020-04-13T13:17:50+10:00"}
{"level":"info","msg":"sent MessageId: 5ccab011-74e4-4386-a852-67e638422ad9, MessageHeader: 7085789e-af7f-44e6-9b06-469ced3c74b0, STORE_REF: A399, TRS_KEY: 333002, SEQUENCE_NUMBER: 3","time":"2020-04-13T13:17:50+10:00"}
```