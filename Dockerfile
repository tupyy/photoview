### Build UI ###
FROM --platform=${BUILDPLATFORM:-linux/amd64} node:18 as ui

ARG REACT_APP_API_ENDPOINT
ENV REACT_APP_API_ENDPOINT=${REACT_APP_API_ENDPOINT}

# Set environment variable UI_PUBLIC_URL from build args, uses "/" as default
ARG UI_PUBLIC_URL
ENV UI_PUBLIC_URL=${UI_PUBLIC_URL:-/}

ARG VERSION
ENV VERSION=${VERSION:-undefined}
ENV REACT_APP_BUILD_VERSION=${VERSION:-undefined}

ARG BUILD_DATE
ENV BUILD_DATE=${BUILD_DATE:-undefined}
ENV REACT_APP_BUILD_DATE=${BUILD_DATE:-undefined}

ARG GIT_COMMIT
ENV COMMIT_SHA=${GIT_COMMIT:-}
ENV REACT_APP_BUILD_COMMIT_SHA=${GIT_COMMIT:-}

RUN mkdir -p /app
WORKDIR /app

# Download dependencies
COPY ui/package*.json /app/
RUN npm install --omit=dev --ignore-scripts

# Build frontend
COPY ui /app
RUN npm run build -- --base=$UI_PUBLIC_URL

### Build API ###
FROM golang:1.20 AS api

WORKDIR /app

COPY api/ .
RUN if [ ! -d "./vendor" ]; then go mod vendor; fi

RUN GOOS=linux GOARCH=amd64 go build -o photoview server.go

### Copy api and ui to production environment ###
FROM fedora:39

WORKDIR /app

RUN dnf update -y && dnf install -y ffmpeg-free perl-Image-ExifTool
RUN useradd -u 1000 photoview

COPY --from=ui /app/dist /app/ui
COPY --from=api /app/photoview /app/photoview
RUN chown -R photoview /app

USER photoview

ENV PHOTOVIEW_LISTEN_IP 127.0.0.1
ENV PHOTOVIEW_LISTEN_PORT 8080

ENV PHOTOVIEW_SERVE_UI 1
ENV PHOTOVIEW_UI_PATH /app/ui

EXPOSE 8080

ENTRYPOINT ["/app/photoview"]
