# golang-test-task
### prerequisite
- go1.17
- macos/linux (tested on macos)
- docker

### Usage
```text
This program creates a Docker container using the given Docker image name,
and the given bash command. This program handles the output logs of the container and send them to the
given AWS CloudWatch group/stream using the given AWS credentials. If the
corresponding AWS CloudWatch group or stream does not exist, it creates it
using the given AWS credentials.

Usage:
  golang-test-task [flags]

Flags:
      --aws-access-key-id string       aws access key id
      --aws-region string              aws region
      --aws-secret-access-key string   aws-secret-access-key
      --bash-command string            bash command
      --cloudwatch-group string        cloudwatch group
      --cloudwatch-stream string       cloudwatch stream
      --docker-image string            docker image name
  -h, --help                           help for golang-test-task
      --print-logs-from-cloudwatch     print log from aws cloudwatch/stream
```

- Example command :
```shell
golang-test-task \
--docker-image python \
--bash-command  $'pip install pip -U && pip install tqdm && python -c \"import time\ncounter = 0\nwhile True:\n\tprint(counter)\n\tcounter = counter + 1\n\ttime.sleep(0.1)"' \
--cloudwatch-group mytest123456789 \
--cloudwatch-stream mytest987654321 \
--aws-access-key-id <aws-access-key-id> \
--aws-secret-access-key <aws-secret-access-key> \
--aws-region us-east-2
```

- Example output
```log
2021/11/17 23:14:02 Docker Image :  bash
2021/11/17 23:14:02 Bash Command :  i=0; while [ $i -ne 3 ]; do i=$(($i+1)); echo "$i"; sleep 2; done
2021/11/17 23:14:02 AWS Access Key Id (encoded) :  xxxxxxxxx
2021/11/17 23:14:02 AWS Secret Access Key (encoded) :  xxxxxxxxx
2021/11/17 23:14:02 CloudWatch group :  mytest12345678
2021/11/17 23:14:02 CloudWatch stream :  mytest87654321
2021/11/17 23:14:02 AWS Region :  us-east-2
2021/11/17 23:14:02 docker client initialized
2021/11/17 23:14:02 listening to interrupt signal
2021/11/17 23:14:05 docker image pulled successfully or already exists dockerImage bash
2021/11/17 23:14:05 container created successfully container id 594d60391f48426bb11e881a3bc356fc8f038e197b416cd79283f51810249e6a bash command i=0; while [ $i -ne 3 ]; do i=$(($i+1)); echo "$i"; sleep 2; done
2021/11/17 23:14:05 container started successfully container id 594d60391f48426bb11e881a3bc356fc8f038e197b416cd79283f51810249e6a
2021/11/17 23:14:05 fetching container logs successfully container id 594d60391f48426bb11e881a3bc356fc8f038e197b416cd79283f51810249e6a
2021/11/17 23:14:05 aws client initialized
2021/11/17 23:14:06 aws cloud watch group already exists mytest12345678
2021/11/17 23:14:07 pushed log to cloudwatch logStr :  !2021-11-17T17:44:05.708584200Z 1
2021/11/17 23:14:08 pushed log to cloudwatch logStr :  !2021-11-17T17:44:07.711011200Z 2
2021/11/17 23:14:09 pushed log to cloudwatch logStr :  !2021-11-17T17:44:09.687765000Z 3
2021/11/17 23:14:11 going to print logs from cloud watch
2021/11/17 23:14:11 *****
2021/11/17 23:14:11 aws client initialized
2021/11/17 23:14:12 aws cloud watch group already exists mytest12345678
2021/11/17 23:14:12 ->   !2021-11-17T17:44:05.708584200Z 1
```

## Build
```shell
make build
```

### Install
```shell
make install
```

### Lint-Fix
```shell
make lint-fix
```

### Test
```shell
make test
```