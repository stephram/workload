package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"reflect"
	"strconv"
	"workload/internal/utils"
	"workload/internal/utils/ulid"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type SalesRecord struct {
}

var (
	numberOfMessages *int
	queueUrl         *string
	messageTemplate  *string
)

func init() {
	utils.GetLogger().Infof("initialised logger")
}

func main() {
	numberOfMessages = flag.Int("number-of-messages", 10, "Number of messages to generate")
	queueUrl = flag.String("queue-url",
		"https://ap-southeast-2.queue.amazonaws.com/531004612469/api-pre-s2c-inbound",
		"URL of the SQS queue to send messages to")
	messageTemplate = flag.String("message-template",
		"/Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/15827625363.json",
		"JSON message template")

	flag.Parse()

	// Read JSON template.
	salesFile, _ := ioutil.ReadFile(*messageTemplate)
	var salesMap map[string]interface{}

	jsonErr := json.Unmarshal([]byte(salesFile), &salesMap)
	if jsonErr != nil {
		log.WithError(jsonErr).Errorf("failed to unmarshal file to map")
		return
	}

	sess, serr := session.NewSessionWithOptions(session.Options{
		Profile: "innovation",
		Config: aws.Config{
			Region: aws.String("ap-southeast-2"),
		},
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	})
	if serr != nil {
		log.WithError(serr).Error("failed to create session")
		return
	}
	svc := sqs.New(sess)

	for i := 1; i <= *numberOfMessages; i++ {
		trsKey := ulid.New()
		trsKey = trsKey[len(trsKey)-12:]
		messageHeader := uuid.New().String()

		generateSale(trsKey, messageHeader, salesMap)

		jsonArr, jsonErr := json.Marshal(salesMap)
		if jsonErr != nil {
			log.WithError(jsonErr).Errorf("failed to marshal message body: %+v", salesMap)
			continue
		}
		sendMessageInput := sqs.SendMessageInput{
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"KEY_ID": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(trsKey),
				},
				"STORE_REF": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("A399"),
				},
				"TABLE_REF": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("sales"),
				},
				"MESSAGE_HEADER": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(messageHeader),
				},
				"SEQUENCE_NUMBER": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(strconv.Itoa(i)),
				},
			},
			MessageBody: aws.String(string(jsonArr)),
			QueueUrl:    aws.String(*queueUrl),
		}
		sendMessage(svc, &sendMessageInput)
	}
}

func sendMessage(sqsClient *sqs.SQS, sendMessageInput *sqs.SendMessageInput) {
	// sendMessageOutput, err := sqsClient.SendMessage(sendMessageInput)
	// if err != nil {
	// 	log.WithError(err).Errorf("failed to send")
	// }
	sendMessageOutput := &sqs.SendMessageOutput{
		MessageId: aws.String(uuid.New().String()),
	}
	log.Infof("sent MessageId: %s, MessageHeader: %s, TRS_KEY: %s, SEQUENCE_NUMBER: %s",
		*sendMessageOutput.MessageId,
		*sendMessageInput.MessageAttributes["MESSAGE_HEADER"].StringValue,
		*sendMessageInput.MessageAttributes["KEY_ID"].StringValue,
		*sendMessageInput.MessageAttributes["SEQUENCE_NUMBER"].StringValue)
}

func generateSale(trsKey string, messageHeader string, salesMap map[string]interface{}) {
	salesMap["SAPStoreReference"] = "A399"
	salesMap["KeyID"] = trsKey
	salesMap["MessageHeader"] = messageHeader

	updateTrsKeyFields(trsKey, salesMap)
}

func updateTrsKeyFields(trsKey string, salesMap map[string]interface{}) {
	s := reflect.ValueOf(salesMap["SourceData"])
	if s.Kind() != reflect.Slice {
		log.Errorf("not a slice")
		return
	}

	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		m := reflect.ValueOf(v)
		if m.Kind() != reflect.Map {
			log.Infof("not a map. Is a %s", m.Type().String())
			continue
		}

		log.Infof("%v", m.MapKeys())

		// if mm, ok := m.MapIndex("TRS_KEY"); ok {
		// 	mm["TRS_KEY"] = trsKey
		// 	log.Infof("%02d: %+v", i, mm)
		// }
	}
}
