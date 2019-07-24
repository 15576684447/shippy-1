FROM ubuntu:14.04
RUN mkdir -p /app
WORKDIR /app
ADD consignment.json /app/consignment.json
ADD consignment-cli /app/consignment-cli

CMD ["./consignment-cli"]