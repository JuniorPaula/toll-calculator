obu:
		  @go build -o bin/obu obu/main.go
		  @./bin/obu

receive:
	    @go build -o bin/receive data_receive/*.go
	    @./bin/receive

calculator:
	    @go build -o bin/calculator distance_calculator/*.go
	    @./bin/calculator

.PHONY:	obu
