# Define variables for repeated use
BINARY_NAME=main
BINARY_PATH=./src/handlers/$(BINARY_NAME)
TEMPLATE=./template.yaml
BUCKET_NAME=weather-bucket
STACK_NAME=mystack
LOCALSTACK_URL=http://localhost:4566
LOCALSTACK_AUTH_TOKEN=ls-KocUduQo-5063-dile-GUFU-samIqiGi728f
FUNCTION_NAME=WeatherFunction


build-docker:
	docker run --rm -v "$(PWD)":/usr/src/weathermap -w /usr/src/weathermap golang:1.15 bash -c "CGO_ENABLED=0 GOOS=linux go build -o ./src/handlers/main ./src/handlers"


# Set execute permissions on the binary (might be necessary on some systems)
chmod:
	sudo chmod +x $(BINARY_PATH)

# Generate an API Gateway event and invoke the function locally
invoke-local: build-docker chmod
	sam local generate-event apigateway aws-proxy | sam local invoke $(FUNCTION_NAME) -t $(TEMPLATE)

# Create an S3 bucket in LocalStack
create-bucket:
	aws --endpoint-url=$(LOCALSTACK_URL) --profile localstack s3 mb s3://$(BUCKET_NAME)

## Package and upload the function to the S3 bucket in LocalStack
bucket-deploy:
	aws --endpoint-url=$(LOCALSTACK_URL) --profile localstack s3 cp $(BINARY_PATH) s3://$(BUCKET_NAME)/$(BINARY_NAME)


cloudformation-deploy:
	aws --endpoint-url=$(LOCALSTACK_URL) --profile localstack cloudformation deploy --stack-name $(STACK_NAME) --template-file $(TEMPLATE) --capabilities CAPABILITY_IAM

deploy: build-docker chmod bucket-deploy cloudformation-deploy
# Deploy to LocalStack
#deploy: build-docker chmod
#	# Package and upload the function to the S3 bucket in LocalStack
#	aws --endpoint-url=$(LOCALSTACK_URL) --profile localstack s3 cp $(BINARY_PATH) s3://$(BUCKET_NAME)/$(BINARY_NAME)
#	# Deploy using CloudFormation in LocalStack
#	aws --endpoint-url=$(LOCALSTACK_URL) --profile localstack cloudformation deploy --stack-name $(STACK_NAME) --template-file $(TEMPLATE) --capabilities CAPABILITY_IAM

# This target is just a convenience for setting up LocalStack S3 bucket initially
setup-localstack: create-bucket

.PHONY: build-docker chmod invoke-local create-bucket deploy setup-localstack
