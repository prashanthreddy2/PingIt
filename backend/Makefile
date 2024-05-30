DOCKER_IMAGE := pingit-deploy
DOCKER_IMAGE_TAG := code_latest

# Build the Go application
build:
	go build -o backend.exe .

# Build Docker image
docker-build:
	minikube image build -t $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) -f ./Dockerfile .

# Push Docker image to GCR
docker-push:
	minikube image load $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG)  

# Clean up
clean:
	del backend.exe
# Run all steps
all: build docker-build clean

#phoen1x0/origin_images:tagname