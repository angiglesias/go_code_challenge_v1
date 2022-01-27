FROM golang:1.16-alpine AS build

WORKDIR /app
COPY . .
RUN apk --no-cache add make
RUN make build

FROM alpine:3.15 as release

WORKDIR /app 
COPY --from=build /app/counter /app/

EXPOSE 3000
ENTRYPOINT [ "/app/counter" ]
CMD [ "-cors" ]