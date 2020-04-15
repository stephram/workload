package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"workload/internal/utils"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

type messageInfo struct {
	Filename    string
	Data        map[string]interface{}
	StoreNumber string
	KeyID       int
	SeqNum      int
	Mq          *sqs.SQS
}

var (
	// app              = kingpin.New("messages", "S2C store side message generator for posting to SQS queue")
	numberOfMessages = kingpin.Flag("number-of-messages", "Number of messages to send").Default(
		"5").Short('n').Int()
	queueUrl = kingpin.Flag("queue-url", "URL of the SQS queue").Default(
		"https://sqs.ap-southeast-2.amazonaws.com/712510509017/api-dev-s2c-inbound").Short('q').URL()
	messageTemplates = kingpin.Flag("message-template", "Filenames containing payloads to use as templates").Short('t').Default(
		"test-data/sales-messages/1808712-body.json",
		"test-data/sales-messages/1808713-body.json",
		"test-data/sales-messages/1808714-body.json").ExistingFiles()
	storeIDs = kingpin.Flag("store-number", "Store numbers use").Short('s').Default(
		"A399", "A301").Strings()
	keyStartID = kingpin.Flag("key-start-id", "Reference key start index").Short('k').Short('k').Default(
		"333000").Int()
	numberOfWorkers = kingpin.Flag("number-of-workers", "Number of workers").Short('w').Default("100").Int()
	awsProfile      = kingpin.Flag("profile", "AWS profile name").Short('p').Default("api-dev").String()
	wDir            string
)

func init() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	utils.GetLogger()            //.Infof("initialised logger")
}

func main() {
	// kingpin.MustParse(app.Parse(os.Args[1:]))
	kingpin.Parse()

	storeNumbers := *storeIDs
	messageTmpls := *messageTemplates

	if *numberOfWorkers > *numberOfMessages {
		*numberOfWorkers = *numberOfMessages
	}
	wDir, _ = os.Getwd()

	sess, serr := session.NewSessionWithOptions(session.Options{
		Profile: *awsProfile,
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
	keyID := *keyStartID

	fmt.Printf("%20s : %d\n", "number-of-messages", *numberOfMessages)
	fmt.Printf("%20s : %d\n", "number-of-workers", *numberOfWorkers)
	fmt.Printf("%20s : (%d) %v\n", "message-templates", len(*messageTemplates), *messageTemplates)
	fmt.Printf("%20s : (%d) %v\n", "store-numbers", len(storeNumbers), storeNumbers)
	fmt.Printf("%20s : %d\n", "key-start-id", keyID)
	fmt.Printf("%20s : %s\n", "queue-url", (*queueUrl).String())

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n') // nolint

	taskInp := make(chan messageInfo, 50)
	taskOut := make(chan string, 50)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	log.Infof("Creating worker %d tasks", *numberOfWorkers)
	go createWorkerTasks(taskInp, taskOut, *numberOfWorkers)

	go waitForTasks(*numberOfMessages, taskOut, &waitGroup)

	defer stopWatch(time.Now(), "Waited %s")

	log.Infof("Send %d messages", *numberOfMessages)
	go sendMessageTasks(messageTmpls, storeNumbers, keyID, svc, taskInp, taskOut)

	waitGroup.Wait()
}

func stopWatch(start time.Time, message string) {
	log.Infof(message, time.Since(start))
}

func waitForTasks(numberOfMessages int, taskOut <-chan string, waitGroup *sync.WaitGroup) {
	// Wait for all of the routines to finish.
	log.Infof("Waiting for output")
	for msgCount := 1; msgCount < numberOfMessages; msgCount++ {
		s := <-taskOut
		log.Infof("%.8d : %s", msgCount, s)
	}
	log.Infof("Received %d messages", numberOfMessages)
	waitGroup.Done()
}

func sendMessageTasks(messageTmpls []string, storeNumbers []string, keyID int, svc *sqs.SQS, taskInp chan<- messageInfo, taskOut chan<- string) {
	for seqNum := 1; seqNum <= *numberOfMessages; seqNum++ {
		// Read JSON template.
		filename := wDir + string(filepath.Separator) + strings.TrimSpace(utils.SelectRandomString(messageTmpls))

		file, fileErr := os.Open(filename)
		if fileErr != nil {
			log.WithError(fileErr).Panicf("Failed to open file: %s", filename)
			// taskOut <- fmt.Sprintf("%s. Failed to open file: %s", errors.Unwrap(fileErr).Error(), filename)
			// return
		}
		// defer file.Close()

		salesFile, fileErr := ioutil.ReadAll(file)
		if fileErr != nil {
			log.WithError(fileErr).Panicf("Failed to read file: %s", filename)
			// taskOut <- fmt.Sprintf("%s. Failed to read file: %s from %s", errors.Unwrap(fileErr).Error(), filename, wDir)
			// return
		}

		fileErr = file.Close()
		if fileErr != nil {
			log.WithError(fileErr).Panicf("Failed to close file: %s", filename)
			// taskOut <- fmt.Sprintf("%s. Failed to close file: %s", errors.Unwrap(fileErr).Error(), filename)
			// return
		}
		// log.Infof("close file: %s", filename)

		var salesMap map[string]interface{}
		jsonErr := json.Unmarshal([]byte(salesFile), &salesMap)
		if jsonErr != nil {
			taskOut <- fmt.Sprintf("%s. Failed to unmarshal file to map: %s", errors.Unwrap(jsonErr).Error(), filename)
			return
		}

		mi := &messageInfo{
			StoreNumber: utils.SelectRandomString(storeNumbers),
			Data:        salesMap,
			SeqNum:      seqNum,
			KeyID:       keyID,
			Mq:          svc,
		}
		taskInp <- *mi
		keyID++
	}
}

func createWorkerTasks(taskInp <-chan messageInfo, taskOut chan<- string, taskCount int) {
	for i := 0; i < taskCount; i++ {
		go func() {
			for {
				sendMessageTask(<-taskInp, taskOut)
			}
		}()
	}
}

func sendMessageTask(mi messageInfo, taskOut chan<- string) {
	// Data setup.
	trsKey := strconv.FormatInt(int64(mi.KeyID), 10)
	messageHeader := uuid.New().String()

	sendMessageInput, cmErr := createSendMessageInput(mi.Data, mi.StoreNumber, trsKey, messageHeader, mi.SeqNum, (*queueUrl).String())
	if cmErr != nil {
		taskOut <- fmt.Sprintf("%s. Failed to Create Message from map", cmErr.Error())
		return
	}
	sendMessageOutput, smErr := sendMessage(mi.Mq, sendMessageInput)
	if smErr != nil {
		taskOut <- fmt.Sprintf("%s. Failed to send message from map", smErr.Error())
		return
	}
	taskOut <- fmt.Sprintf("sent MessageId: %s, MessageHeader: %s, STORE_REF: %s, TRS_KEY: %s, SEQUENCE_NUMBER: %s",
		*sendMessageOutput.MessageId,
		*sendMessageInput.MessageAttributes["MESSAGE_HEADER"].StringValue,
		*sendMessageInput.MessageAttributes["STORE_REF"].StringValue,
		*sendMessageInput.MessageAttributes["KEY_ID"].StringValue,
		*sendMessageInput.MessageAttributes["SEQUENCE_NUMBER"].StringValue)
}

func createSendMessageInput(salesMap map[string]interface{}, storeNumber, trsKey, messageHeader string, seqNum int, queueURL string) (*sqs.SendMessageInput, error) {
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
		QueueUrl:    aws.String(queueURL),
	}
	return &sendMessageInput, nil
}

func sendMessage(sqsClient *sqs.SQS, sendMessageInput *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	sendMessageOutput, err := sqsClient.SendMessage(sendMessageInput)
	if err != nil {
		log.WithError(err).Errorf("failed to send Message: ")
		return nil, err
	}
	return sendMessageOutput, nil
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
	return // nolint
}
