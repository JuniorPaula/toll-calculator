obu:
		  @go build -o bin/obu obu/main.go
		  @./bin/obu

receive:
	    @go build -o bin/receive data_receive/main.go
	    @./bin/receive

.PHONY:	obu
