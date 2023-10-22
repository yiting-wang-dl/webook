.PHONY: docker
# create kubernetes container image
docker:
	# delete last complied file
	@rm webook || true
	@go mod tidy
	# specify to compile an executable file to run on the ARM architecture of the Linux operating system
	@@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f ytw/webook:v0.0.1
	@docker build -t ytw/webook:v0.0.1 .