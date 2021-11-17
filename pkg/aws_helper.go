package pkg

import (
	"bufio"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/google/uuid"
	"io"
	"log"
	"time"
)

// streamLogsToAWSCloudWatchGrp streams logs to user aws cloud watch group
func (t *task) streamLogsToAWSCloudWatchGrp(containerLogStream io.ReadCloser) {
	sess := t.getAwsClientOrDie()
	cwl := t.getAwsCloudwatchSession(sess)
	t.createAWSCloudWatchGrpIfNotExist(cwl)
	t.streamToAWSCloudWatchGrp(containerLogStream, cwl)
}

// getAwsClientOrDie returns aws session
func (t *task) getAwsClientOrDie() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(t.awsRegion),
		Credentials: credentials.NewStaticCredentials(t.awsAccessKeyID, t.awsSecretAccessKey, ""),
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("aws client initialized")
	return sess
}

// getAwsCloudwatchSession returns cloudwatch session
func (t *task) getAwsCloudwatchSession(sess *session.Session) *cloudwatchlogs.CloudWatchLogs {
	return cloudwatchlogs.New(sess)
}

// createAWSCloudWatchGrpIfNotExist creates cloudwatchgroup if it doesn't exist
func (t *task) createAWSCloudWatchGrpIfNotExist(cwl *cloudwatchlogs.CloudWatchLogs) {
	if t.isAWSCloudWatchGroupExists(cwl) {
		return
	}

	t.createAWSCloudWatchLogGroup(cwl)
}

// isAWSCloudWatchGroupExists checks in CW group exists
func (t *task) isAWSCloudWatchGroupExists(cwl *cloudwatchlogs.CloudWatchLogs) bool {
	resp, err := cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, logGroup := range resp.LogGroups {
		if *logGroup.LogGroupName == t.cloudWatchGroup {
			log.Println("aws cloud watch group already exists", t.cloudWatchGroup)
			return true
		}
	}
	return false
}

// createAWSCloudWatchLogGroup creates aws cloud watch log group
func (t *task) createAWSCloudWatchLogGroup(cwl *cloudwatchlogs.CloudWatchLogs) {
	_, err := cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &t.cloudWatchGroup,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("aws cloud watch group created", t.cloudWatchGroup)
}

// streamToAWSCloudWatchGrp reads container log stream and pushes to cloudwatch
func (t *task) streamToAWSCloudWatchGrp(containerLogStream io.ReadCloser, cwl *cloudwatchlogs.CloudWatchLogs) {
	var logQueue []*cloudwatchlogs.InputLogEvent
	sequenceToken := ""
	scanner := bufio.NewScanner(containerLogStream)
	for scanner.Scan() {
		logStr := scanner.Text()
		logQueue = append(logQueue, &cloudwatchlogs.InputLogEvent{
			Message:   &logStr,
			Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
		})

		if len(logQueue) > 0 {
			input := cloudwatchlogs.PutLogEventsInput{
				LogEvents:    logQueue,
				LogGroupName: &t.cloudWatchGroup,
			}

			if sequenceToken == "" {
				err := t.createLogStream(cwl)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				input = *input.SetSequenceToken(sequenceToken)
			}

			input = *input.SetLogStreamName(t.cloudWatchStream)

			resp, err := cwl.PutLogEvents(&input)
			if err != nil {
				log.Fatal(err)
			}

			if resp != nil {
				sequenceToken = *resp.NextSequenceToken
			}
			logQueue = []*cloudwatchlogs.InputLogEvent{}
		}
		log.Println("pushed log to cloudwatch", "logStr : ", logStr)
	}
}

// createLogStream will make a new logStream with a random uuid as its name.
func (t *task) createLogStream(cwl *cloudwatchlogs.CloudWatchLogs) error {
	name := uuid.New().String()

	_, err := cwl.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &t.cloudWatchGroup,
		LogStreamName: &name,
	})

	t.cloudWatchStream = name
	return err
}

// printLogsFromCW fetches and prints logs from CW. Only for debugging purpose
func (t *task) printLogsFromCW() {
	sess := t.getAwsClientOrDie()
	cwl := t.getAwsCloudwatchSession(sess)
	if !t.isAWSCloudWatchGroupExists(cwl) {
		log.Fatal("unable to fetch logs from cloud watch; group does not exist")
	}

	resp := t.getAWSCloudWatchLogEvent(cwl)

	gotToken := ""
	nextToken := ""

	for _, event := range resp.Events {
		gotToken = nextToken
		nextToken = *resp.NextForwardToken

		if gotToken == nextToken {
			break
		}

		log.Println("->  ", *event.Message)
	}
	log.Println()
	log.Println("finished printing logs from aws cloud watch")
}

// getAWSCloudWatchLogEvent return log event output
func (t *task) getAWSCloudWatchLogEvent(cwl *cloudwatchlogs.CloudWatchLogs) *cloudwatchlogs.GetLogEventsOutput {
	resp, err := cwl.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(t.cloudWatchGroup),
		LogStreamName: aws.String(t.cloudWatchStream),
	})
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
