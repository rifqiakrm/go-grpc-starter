FROM rifqiakrm/grpc-go-base-builder:1.0.3-alpine AS base

# define timezone
ENV TZ Asia/Jakarta

# define work directory
WORKDIR /app

# copy the sourcecode
COPY . .

# generate protocol buffers
RUN make generate-pb

# build beedoor exec
RUN cd /app/cmd/server && go mod vendor && go build -o backend-service

FROM alpine:3.16

WORKDIR /app

COPY --from=base app/gcloud-credentials.json .
COPY --from=base app/cmd/server/backend-service .
COPY --from=base app/template .

# EXPOSE 8080 is the port that the REST API will be exposed on
EXPOSE 8080
# EXPOSE 8080 is the port that the GRPC will be exposed on. But if deployed in cloud run just use the 8080 port
EXPOSE 8081

CMD [ "./backend-service" ]