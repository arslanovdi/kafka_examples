gen:
	protoc --go_out=. kafkaExample.v1.proto

PROTOC_ZIP:="protoc-28.3-linux-x86_64.zip"
install:
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v28.3/$(PROTOC_ZIP)
	sudo unzip -o $(PROTOC_ZIP) -d /usr/local bin/protoc
	sudo unzip -o $(PROTOC_ZIP) -d /usr/local 'include/*'
	rm -f $(PROTOC_ZIP)