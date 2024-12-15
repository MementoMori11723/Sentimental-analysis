run: 
	@docker build -t sentimental-analysis .
	@docker run -it -p 11000:8080 sentimental-analysis

