deploy: docker-build docker-tag aws-login docker-push lambda-update

docker-build:
	docker build --platform linux/amd64 -t mtg-price-scrapper .

docker-tag:
	docker tag mtg-price-scrapper 206363131200.dkr.ecr.ap-southeast-1.amazonaws.com/mtg-price-scrapper:latest

docker-push:
	docker push 206363131200.dkr.ecr.ap-southeast-1.amazonaws.com/mtg-price-scrapper:latest

lambda-create:
	aws lambda create-function \
      --function-name mtg-price-scrapper \
      --package-type Image \
      --code ImageUri=206363131200.dkr.ecr.ap-southeast-1.amazonaws.com/mtg-price-scrapper:latest \
      --role arn:aws:iam::206363131200:role/lambda-mtg

lambda-update:
	aws lambda update-function-code \
      --function-name mtg-price-scrapper \
      --image-uri 206363131200.dkr.ecr.ap-southeast-1.amazonaws.com/mtg-price-scrapper:latest

aws-login:
	aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 206363131200.dkr.ecr.ap-southeast-1.amazonaws.com
