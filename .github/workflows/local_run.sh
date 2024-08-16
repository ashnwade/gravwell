#!/bin/bash

set -e

govulncheck -test ./netflow/...
govulncheck -test ./manager/...
govulncheck -test ./ipexist/...
govulncheck -test ./generators/ipgen/...
govulncheck -test ./generators/base/...
govulncheck -test ./generators/gravwellGenerator/...
govulncheck -test ./filewatch/...
govulncheck -test ./client/...
govulncheck -test ./chancacher/...
govulncheck -test ./timegrinder/...
govulncheck -test ./ingest
govulncheck -test ./ingest/entry/...
govulncheck -test ./ingest/config/...
govulncheck -test ./ingest/log
govulncheck -test ./ingest/processors
govulncheck -test ./ingest/processors/tags
govulncheck -test ./ingest/processors/plugin
govulncheck -test ./tools/...
govulncheck -test ./kitctl/...
govulncheck -test ./migrate/...
govulncheck -test ./ingesters/s3Ingester
govulncheck -test ./ingesters/HttpIngester
govulncheck -test ./ingesters/pcapFileIngester
govulncheck -test ./ingesters/collectd
govulncheck -test ./ingesters/hackernews_ingester
govulncheck -test ./ingesters/base
govulncheck -test ./ingesters/massFile
govulncheck -test ./ingesters/MSGraphIngester
govulncheck -test ./ingesters/kafka_consumer
govulncheck -test ./ingesters/reddit_ingester
govulncheck -test ./ingesters/O365Ingester
govulncheck -test ./ingesters/args
govulncheck -test ./ingesters/version
govulncheck -test ./ingesters/sqsIngester
govulncheck -test ./ingesters/diskmonitor
govulncheck -test ./ingesters/session
govulncheck -test ./ingesters/snmp
govulncheck -test ./ingesters/xlsxIngester
govulncheck -test ./ingesters/multiFile
govulncheck -test ./ingesters/Shodan
govulncheck -test ./ingesters/reimport
govulncheck -test ./ingesters/SimpleRelay
govulncheck -test ./ingesters/KinesisIngester
govulncheck -test ./ingesters/netflow
govulncheck -test ./ingesters/AzureEventHubs
govulncheck -test ./ingesters/utils
govulncheck -test ./ingesters/IPMIIngester
govulncheck -test ./ingesters/regexFile
govulncheck -test ./ingesters/PacketFleet
govulncheck -test ./ingesters/canbus
govulncheck -test ./ingesters/GooglePubSubIngester
govulncheck -test ./ingesters/fileFollow
govulncheck -test ./ingesters/singleFile
GOOS=windows govulncheck -test ./ingesters/winevents
GOOS=windows govulncheck -test ./winevent/...

go test -v ./generators/ipgen
go test -v ./chancacher
go test -v ./ingest
go test -v ./ingest/entry
go test -v ./ingest/processors
go test -v ./ingest/processors/plugin
go test -v ./ingest/config
go test -v ./ingest/log
go test -v ./timegrinder
go test -v ./filewatch
go test -v ./ingesters/utils
go test -v ./ingesters/kafka_consumer
go test -v ./ingesters/SimpleRelay
go test -v ./ipexist
go test -v ./netflow
go test -v ./client/...

go build -o /dev/null ./generators/gravwellGenerator
go build -o /dev/null ./manager
go build -o /dev/null ./migrate
go build -o /dev/null ./tools/timetester
go build -o /dev/null ./timegrinder/cmd
go build -o /dev/null ./ipexist/textinput
go build -o /dev/null ./kitctl
GOOS=windows go build -o /dev/null ./ingesters/fileFollow
GOOS=windows go build -o /dev/null ./ingesters/winevents
GOOS=windows go build ./generators/windowsEventGenerator
go build -o /dev/null ./ingesters/massFile
go build -o /dev/null ./ingesters/diskmonitor
go build -o /dev/null ./ingesters/xlsxIngester
go build -o /dev/null ./ingesters/reimport
go build -o /dev/null ./ingesters/version
go build -o /dev/null ./ingesters/canbus
go build -o /dev/null ./ingesters/reddit_ingester
go build -o /dev/null ./ingesters/hackernews_ingester
go build -o /dev/null ./ingesters/multiFile
go build -o /dev/null ./ingesters/session
go build -o /dev/null ./ingesters/regexFile
go build -o /dev/null ./ingesters/Shodan
go build -o /dev/null ./ingesters/singleFile
go build -o /dev/null ./ingesters/pcapFileIngester
GOOS=darwin GOARCH=amd64 go build -o /dev/null ./ingesters/fileFollow
GOOS=darwin GOARCH=arm64 go build -o /dev/null ./ingesters/fileFollow
GOOS=linux GOARCH=arm64 go build -o /dev/null ./ingesters/fileFollow

/bin/bash ./ingesters/test/build.sh ./ingesters/GooglePubSubIngester ./ingesters/test/configs/pubsub_ingest.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/AzureEventHubs ingesters/test/configs/azure_event_hubs.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/HttpIngester ingesters/test/configs/gravwell_http_ingester.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/collectd ingesters/test/configs/collectd.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/netflow ingesters/test/configs/netflow_capture.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/KinesisIngester ingesters/test/configs/kinesis_ingest.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/kafka_consumer ingesters/test/configs/kafka.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/MSGraphIngester ingesters/test/configs/msgraph_ingest.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/IPMIIngester ingesters/test/configs/ipmi.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/fileFollow ingesters/test/configs/file_follow.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/s3Ingester ingesters/test/configs/s3.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/snmp ingesters/test/configs/snmp.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/sqsIngester ingesters/test/configs/sqs.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/networkLog ingesters/test/configs/network_capture.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/SimpleRelay ingesters/test/configs/simple_relay.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/O365Ingester ingesters/test/configs/o365_ingest.conf
/bin/bash ./ingesters/test/build.sh ./ingesters/PacketFleet ingesters/test/configs/packet_fleet.conf

