ceratopoGon
=======================================


# About

ceratopoGon is MQTT-SN Gateway written in golang.

# Requirements

 * golang > 1.6
 * eclipse/paho.mqtt.golang
 * go-yaml/yaml

# How to use
## get dependent library

```
$ git clone https://github.com/Kmotiko/ceratopoGon.git
$ cd ceratopoGon
$ go get github.com/Kmotiko/ceratopoGon/gateway
```

## execute command

Command format is like bellow.

```
go run ceratopogon.go -c <configfile> -t <predefined topic file>
```

For exmaple, you are able to launch ceratopoGon with following command.

```
$ go run ceratopogon.go -c ceratopogon.conf -t sample_topics.yml
```

## config file

Config file format is json style.
For example, like as bellow.

```
{
  "IsAggregate" : true,
  "Host" : "127.0.0.1",
  "Port" : 1884,
  "BrokerHost" : "127.0.0.1",
  "BrokerPort" : 1883,
  "BrokerUser" : "",
  "BrokerPassword" : "",
  "LogFilePath" : "ceratopogon.log",
  "MessageQueueSize" : 100
}
```

Following config parameters is available.

|parameter          |description                                    |Value Type  |exmaple            |
|-------------------|-----------------------------------------------|------------|-------------------|
|IsAggregate        |execute gw as aggregate mode or not            |bool        |true               |
|Host               |listen host address                            |string      |"127.0.0.1"        |
|Port               |listen port number                             |string      |"1884"             |
|BrokerHost         |broker host address                            |string      |"127.0.0.1"        |
|BrokerPort         |broker port number                             |string      |"1883"             |
|BrokerUser         |username to use broker authentication          |string      |"user"             |
|BrokerPassword     |password to use broker authentication          |string      |"pass"             |
|LogFilePath        |path of log file                               |string      |"ceratopogon.log"  |
|MessageQueueSize   |Size of queue only used in Transparent Mode    |int         |100                |

## predefined topic file

Predefined Topic file format is yaml style.

```
---
<TOPIC NAME>: <ID>
<TOPIC NAME>: <ID>
```


Following is example of predefined topic file.

```
---
topic1: 1
topic2: 2
```


# Support Status
## Transparent Gateway
### Message

 - [ ] ADVERTISE
 - [ ] SEARCHGW
 - [ ] GWINFO
 - [ ] CONNECT
 - [ ] CONNACK
 - [ ] WILLTOPICREQ
 - [ ] WILLTOPIC
 - [ ] WILLMSGREQ
 - [ ] WILLMSG
 - [ ] REGISTER
 - [ ] REGACK
 - [ ] PUBLISH
 - [ ] PUBACK
 - [ ] PUBCOMP
 - [ ] PUBREC
 - [ ] PUBREL
 - [ ] SUBSCRIBE
 - [ ] SUBACK
 - [ ] UNSUBSCRIBE
 - [ ] UNSUBACK
 - [ ] PINGREQ
 - [ ] PINGRESP
 - [ ] DISCONNECT
 - [ ] WILLTOPICUPD
 - [ ] WILLTOPICRESP
 - [ ] WILLMSGUPD
 - [ ] WILLMSGRESP

## Aggregate Gateway
### Message

 - [ ] ADVERTISE
 - [ ] SEARCHGW
 - [ ] GWINFO
 - [ ] CONNECT
 - [ ] CONNACK
 - [ ] WILLTOPICREQ
 - [ ] WILLTOPIC
 - [ ] WILLMSGREQ
 - [ ] WILLMSG
 - [ ] REGISTER
 - [ ] REGACK
 - [ ] PUBLISH
 - [ ] PUBACK
 - [ ] PUBCOMP
 - [ ] PUBREC
 - [ ] PUBREL
 - [ ] SUBSCRIBE
 - [ ] SUBACK
 - [ ] UNSUBSCRIBE
 - [ ] UNSUBACK
 - [ ] PINGREQ
 - [ ] PINGRESP
 - [ ] DISCONNECT
 - [ ] WILLTOPICUPD
 - [ ] WILLTOPICRESP
 - [ ] WILLMSGUPD
 - [ ] WILLMSGRESP

## Other features

  - [ ] CleanSession
  - [ ] Will Topic/Msg
  - [ ] Sleep State
  - [ ] Predefined Topic

## License

 This software is distributed with MIT license.

 ```
 Copyright (c) 2017 Kmotiko
 ```

 This software use eclipse/paho.mqtt.golang which is licensed under the Eclipse Public License 1.0.

 This software use go-yaml/yaml which is licensed under the Apache License 2.0.
