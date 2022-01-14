build_amd64: 
	docker build . -t dueruen/learning-kubernetes:simple -f Dockerfile

build_arm: 
	docker build . -t dueruen/learning-kubernetes:simple-arm -f Dockerfile.ARM

push_amd64_simple: 
	docker push dueruen/learning-kubernetes:simple
	
push_arm_simple: 
	docker push dueruen/learning-kubernetes:simple-arm

build_push_arm_helm: 
	docker build . -t dueruen/learning-kubernetes:helm-arm -f Dockerfile.ARM
	docker push dueruen/learning-kubernetes:helm-arm

run: 
	NAME=alpha go run .

run_beta: 
	NAME=beta HTTP_PORT=8081 GET_URI="http://localhost:8080" POST_URI="http://localhost:8082" go run .

run_gamma: 
	NAME=gamma HTTP_PORT=8082 GET_URI="http://localhost:8081" POST_SLEEP=3 go run .