build_amd64: 
	docker build . -t dueruen/learning-kubernetes:simple -f Dockerfile

build_arm: 
	docker build . -t dueruen/learning-kubernetes:simple-arm -f Dockerfile.ARM

push_amd64: 
	docker push dueruen/learning-kubernetes:simple
	
push_arm: 
	docker push dueruen/learning-kubernetes:simple-arm

run: 
	NAME=alpha go run .

run_beta: 
	NAME=beta HTTP_PORT=8081 GET_PORT=8080 POST_PORT=8082 go run .

run_gamma: 
	NAME=gamma HTTP_PORT=8082 GET_PORT=8081 POST_SLEEP=3 go run .