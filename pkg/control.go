package pkg

import (
	"log"
	"os"
	"os/signal"
)

// task contains task params
type task struct {
	dockerImage, bashCommand, cloudWatchGroup, cloudWatchStream, awsSecretAccessKey, awsAccessKeyID, awsRegion string
}

// NewTask returns new task
func NewTask(dockerImage, bashCommand, cloudWatchGroup, cloudWatchStream, awsSecretAccessKey, awsAccessKeyID, awsRegion string) *task {
	return &task{
		dockerImage:        dockerImage,
		bashCommand:        bashCommand,
		cloudWatchGroup:    cloudWatchGroup,
		cloudWatchStream:   cloudWatchStream,
		awsSecretAccessKey: awsSecretAccessKey,
		awsAccessKeyID:     awsAccessKeyID,
		awsRegion:          awsRegion,
	}
}

// Process the task
func (t *task) Process() {
	containerLogStream := t.runAndGetDockerContainerLogStream()
	t.streamLogsToAWSCloudWatchGrp(containerLogStream)
}

func (t *task) PrintLogsFromCloudWatch() {
	t.printLogsFromCW()
}

// ListenInterruptSignal listens the interrupt and runs cleanup (if required)
func (t *task) ListenInterruptSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		log.Println("listening to interrupt signal")
		select {
		case <-c:
			log.Println("got interrupt signal, exiting gracefully")
			os.Exit(1)
		}
	}()
}
