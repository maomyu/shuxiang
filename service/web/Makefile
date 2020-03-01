build:
	go build ./service/web
	docker build -t webservice .
	docker run --name sxwebservice -p 8100:8100 -d webservice
