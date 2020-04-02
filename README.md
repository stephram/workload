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
$ build/messages --help
```
Help output
```
Usage of build/messages:
  -key-start-id int
    	Start ID for TRS_KEY (default 333000)
  -message-templates string
    	JSON message template(s), comma separated list of filenames (default "/Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/1808221.json, /Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/15827625345.json")
  -number-of-messages int
    	Number of messages to generate (default 10)
  -profile string
    	AWS Profile name (default "innovation")
  -queue-url string
    	URL of the SQS queue to send messages to (default "https://ap-southeast-2.queue.amazonaws.com/531004612469/api-pre-s2c-inbound")
  -store-numbers string
    	Comma separated store numbers (default "A399, A301")
```
#### Running
```
$ ./messages_go  --number-of-messages 5000 --store-numbers "A338, A394, A415, A440, A515, A545" --key-start-id 335000 --message-templates /Users/sg/Development/api/s2c/prescription-body2.json
{"level":"info","msg":"sent MessageId: bd83dd96-d74b-427b-8993-780f763186b5, MessageHeader: 5a6fbedd-9ce6-4ee2-a4c7-eafaee16396d, STORE_REF: A515, TRS_KEY: 335000, SEQUENCE_NUMBER: 1","time":"2020-03-31T17:19:14+11:00"}
{"level":"info","msg":"sent MessageId: 9311181a-fb11-47ea-ab03-1f3d9b047bca, MessageHeader: a51a4346-2062-4c19-87ed-fad6307b559c, STORE_REF: A440, TRS_KEY: 335001, SEQUENCE_NUMBER: 2","time":"2020-03-31T17:19:14+11:00"}
{"level":"info","msg":"sent MessageId: a7bd9334-d9ca-4521-b330-9c0cff95ec5f, MessageHeader: cdd0349b-2158-4234-a716-9b727f2e2279, STORE_REF: A338, TRS_KEY: 335002, SEQUENCE_NUMBER: 3","time":"2020-03-31T17:19:14+11:00"}
{"level":"info","msg":"sent MessageId: 96735a6c-9ec9-4862-9a23-507d2c746365, MessageHeader: 1ff2416b-745a-4934-bf48-7ca01dd9f714, STORE_REF: A545, TRS_KEY: 335003, SEQUENCE_NUMBER: 4","time":"2020-03-31T17:19:14+11:00"}
{"level":"info","msg":"sent MessageId: 1c638909-4815-4c43-a13c-0588aca6ca99, MessageHeader: fa27231d-c8cd-4a0e-9f02-df3d05ef7922, STORE_REF: A394, TRS_KEY: 335004, SEQUENCE_NUMBER: 5","time":"2020-03-31T17:19:14+11:00"}
...
```