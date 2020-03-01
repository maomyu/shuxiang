build:
	go build ./service/comentservice
	docker build -t commentservicev3 .
	docker tag commentservicev3:latest registry.cn-shenzhen.aliyuncs.com/mowuya/program:commentservicev3
	docker push registry.cn-shenzhen.aliyuncs.com/mowuya/program:commentservicev3


