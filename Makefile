check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/cmd/swagger
swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models
swagger-client:
	mkdir sdk
	cd sdk;swagger generate client -f ../swagger.yaml -A product-api
clean:
	rm -rf sdk

