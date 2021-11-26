package cmd

import (
	"encoding/base64"
	"github.com/snghnaveen/golang-test-task/pkg"
	"github.com/spf13/cobra"
	"log"
)

var (
	dockerImage, bashCommand, cloudWatchGroup, cloudWatchStream, awsSecretAccessKey, awsAccessKeyID, awsRegion string
	printLogsFromCW                                                                                            bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "golang-test-task",
	Short: "Creates docker container and sends log to AWS group/stream",
	Long: `This program creates a Docker container using the given Docker image name,
and the given bash command. This program handles the output logs of the container and send them to the
given AWS CloudWatch group/stream using the given AWS credentials. If the
corresponding AWS CloudWatch group or stream does not exist, it creates it
using the given AWS credentials.`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Docker Image : ", dockerImage)
		log.Println("Bash Command : ", bashCommand)
		log.Println("AWS Access Key Id (encoded) : ", base64.StdEncoding.EncodeToString([]byte(awsAccessKeyID)))
		log.Println("AWS Secret Access Key (encoded) : ", base64.StdEncoding.EncodeToString([]byte(awsSecretAccessKey)))
		log.Println("CloudWatch group : ", cloudWatchGroup)
		log.Println("CloudWatch stream : ", cloudWatchStream)
		log.Println("AWS Region : ", awsRegion)

		control := pkg.NewTask(dockerImage, bashCommand, cloudWatchGroup, cloudWatchStream, awsSecretAccessKey, awsAccessKeyID, awsRegion)
		control.ListenInterruptSignal()
		control.Process()
		if printLogsFromCW {
			log.Println("going to print logs from cloud watch")
			log.Println("*****")
			control.PrintLogsFromCloudWatch()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVar(&dockerImage, "docker-image", "", "docker image name")
	if err := rootCmd.MarkFlagRequired("docker-image"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().StringVar(&bashCommand, "bash-command", "", "bash command")
	if err := rootCmd.MarkFlagRequired("bash-command"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().StringVar(&cloudWatchGroup, "cloudwatch-group", "", "cloudwatch group")
	if err := rootCmd.MarkFlagRequired("cloudwatch-group"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().StringVar(&cloudWatchStream, "cloudwatch-stream", "", "cloudwatch stream")
	if err := rootCmd.MarkFlagRequired("cloudwatch-stream"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().StringVar(&awsAccessKeyID, "aws-access-key-id", "", "aws access key id")
	if err := rootCmd.MarkFlagRequired("aws-access-key-id"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().StringVar(&awsSecretAccessKey, "aws-secret-access-key", "", "aws-secret-access-key")
	if err := rootCmd.MarkFlagRequired("aws-secret-access-key"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().StringVar(&awsRegion, "aws-region", "", "aws region")
	if err := rootCmd.MarkFlagRequired("aws-region"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().BoolVar(&printLogsFromCW, "print-logs-from-cloudwatch", false, "print log from aws cloudwatch/stream")
}
