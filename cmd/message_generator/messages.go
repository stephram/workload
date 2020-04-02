package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/rand"
	"reflect"
	"strconv"
	"time"
	"workload/internal/utils"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

type SalesRecord struct {
}

var (
	numberOfMessages int
	queueUrl         string
	messageTemplates string
	storeIDs         string
	keyStartID       int
	awsProfile       string
)

func init() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	utils.GetLogger()            //.Infof("initialised logger")
}

func main() {
	flag.StringVar(&awsProfile,
		"profile",
		"innovation",
		"AWS Profile name")
	flag.StringVar(&queueUrl,
		"queue-url",
		"https://ap-southeast-2.queue.amazonaws.com/531004612469/api-pre-s2c-inbound",
		"URL of the SQS queue to send messages to")
	flag.IntVar(&numberOfMessages,
		"number-of-messages",
		10,
		"Number of messages to generate")
	flag.StringVar(&messageTemplates,
		"message-templates",
		// "/Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/15827625363.json",
		"/Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/1808221.json, /Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/15827625345.json",
		"JSON message template(s), comma separated list of filenames")
	flag.StringVar(&storeIDs,
		"store-numbers",
		"A399, A301",
		"Comma separated store numbers")
	flag.IntVar(&keyStartID,
		"key-start-id",
		333000,
		"Start ID for TRS_KEY")

	flag.Parse()

	storeNumbers := utils.ParseCommaSeparatedStrings(storeIDs)
	messageTmpls := utils.ParseCommaSeparatedFiles(messageTemplates)

	sess, serr := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
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

	keyID := keyStartID

	for i := 1; i <= numberOfMessages; i++ {
		// Read JSON template.
		filename := utils.SelectRandomString(messageTmpls)
		salesFile, _ := ioutil.ReadFile(filename)

		var salesMap map[string]interface{}
		jsonErr := json.Unmarshal([]byte(salesFile), &salesMap)
		if jsonErr != nil {
			log.WithError(jsonErr).Errorf("failed to unmarshal file to map: %s", filename)
			return
		}

		// Data setup.
		storeNumber := utils.SelectRandomString(storeNumbers)
		trsKey := strconv.FormatInt(int64(keyID), 10)
		messageHeader := uuid.New().String()

		sendMessageInput, smErr := createSendMessageInput(salesMap, storeNumber, trsKey, messageHeader, i)
		if smErr != nil {
			log.WithError(smErr).Errorf("failed to createMessage with file: %s", filename)
			continue
		}
		sendMessage(svc, sendMessageInput)
		keyID++
	}
}

func createSendMessageInput(salesMap map[string]interface{}, storeNumber, trsKey, messageHeader string, seqNum int) (*sqs.SendMessageInput, error) {
	updateSale(trsKey, storeNumber, messageHeader, salesMap)

	jsonArr, jsonErr := json.Marshal(salesMap)
	if jsonErr != nil {
		return nil, jsonErr
	}
	sendMessageInput := sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"KEY_ID": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(trsKey),
			},
			"STORE_REF": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(storeNumber),
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
				StringValue: aws.String(strconv.Itoa(seqNum)),
			},
		},
		MessageBody: aws.String(string(jsonArr)),
		QueueUrl:    aws.String(queueUrl),
	}
	return &sendMessageInput, nil
}

func sendMessage(sqsClient *sqs.SQS, sendMessageInput *sqs.SendMessageInput) {
	sendMessageOutput, err := sqsClient.SendMessage(sendMessageInput)
	if err != nil {
		log.WithError(err).Errorf("failed to send Message: ")
		return
	}
	log.Infof("sent MessageId: %s, MessageHeader: %s, STORE_REF: %s, TRS_KEY: %s, SEQUENCE_NUMBER: %s",
		*sendMessageOutput.MessageId,
		*sendMessageInput.MessageAttributes["MESSAGE_HEADER"].StringValue,
		*sendMessageInput.MessageAttributes["STORE_REF"].StringValue,
		*sendMessageInput.MessageAttributes["KEY_ID"].StringValue,
		*sendMessageInput.MessageAttributes["SEQUENCE_NUMBER"].StringValue)

}

func updateSale(trsKey string, storeRef string, messageHeader string, salesMap map[string]interface{}) {
	salesMap["MessageHeader"] = messageHeader
	salesMap["KeyID"] = trsKey
	salesMap["SAPStoreReference"] = storeRef

	updateTrsKeyFields(trsKey, salesMap)
}

func updateTrsKeyFields(trsKey string, salesMap map[string]interface{}) {
	v := reflect.ValueOf(salesMap["TableReference"])
	if v.Kind() != reflect.Map {
		log.Errorf("TableReference is not a Map")
		return
	}
	v.SetMapIndex(reflect.ValueOf("ReferenceKeyId"), reflect.ValueOf(trsKey))

	s := reflect.ValueOf(salesMap["SourceData"])
	if s.Kind() != reflect.Slice {
		log.Errorf("not a slice")
		return
	}
	sd, ok := s.Interface().([]interface{})
	if !ok {
		return
	}
	log.Debugf("sd = %s", sd)

	for _, sdi := range sd {
		log.Debugf("sdi = %+v", sdi)

		sdim, ok := sdi.([]interface{})
		if !ok {
			continue
		}
		log.Debugf("sdim = %+v", sdim)

		for _, sdimi := range sdim {
			sdimim, ok := sdimi.(map[string]interface{})
			if !ok {
				continue
			}
			log.Debugf("sdimim = %+v", sdimim)
			sdimim["TRS_KEY"] = trsKey
		}
	}
	return
}
