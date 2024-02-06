FROM node:20-bookworm AS node_builder

WORKDIR /src
COPY frontend/ /src
RUN bash scripts/build.sh

FROM golang:1.21-bookworm AS go_builder

WORKDIR /src
COPY ./ /src
RUN CGO_ENABLED=0 GOOS=linux go build -o api backend/cmd/api/main.go

FROM nginx:1-alpine AS node_final

COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=node_builder /src/dist/ /app

CMD ["nginx", "-g", "daemon off;"]

FROM gcr.io/distroless/static-debian12 AS go_final

COPY --from=go_builder /src/api /api

CMD ["/api"]