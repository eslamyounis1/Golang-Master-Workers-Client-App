#-ip1 http://1234.234:8091/fasta -ip2 http://1234.234:8092/fasta -ip3 http://1234.234:8093/fasta -mapr http://1234.234:8094/fasta
PORTS := 8090 8091 8092 8093 8094

main:
	go run master/main.go -port 8011 -ip1  https://8d93-105-35-226-5.eu.ngrok.io/fasta -ip2  https://f7c1-105-40-174-223.eu.ngrok.io/fasta

mapr:
	go run mapreduce/main.go 8010

slaves:
	go run slave/main.go 8097 slave1 &
	go run slave/main.go 8098 slave2 &
	go run slave/main.go 8099 slave3 &

slave1:
	go run slave/main.go 8097 slave1

slave2:
	go run slave/main.go 8098 slave2

slave3:
	go run slave/main.go 8099 slave3

cli:
	go run client/main.go

nuke:
	@for port in $(PORTS); do \
		pid=$$(lsof -t -i :$$port); \
		if [ "$$pid" != "" ]; then \
			kill -9 $$pid; \
			echo "Service running on port $$port has been killed."; \
		else \
			echo "No service running on port $$port."; \
		fi \
	done