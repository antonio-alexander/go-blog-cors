ARG GIT_BRANCH=no_git_branch_provided
ARG GIT_COMMIT=no_git_commit_provided
ARG PLATFORM=linux/amd64
ARG GO_ARCH=amd64
ARG GO_ARM=7

FROM golang:alpine AS build

ENV GOPROXY=https://proxy.golang.org,direct

ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GO_ARCH
ARG GO_ARM

WORKDIR /go/src/go-blog-cors

COPY . /go/src/go-blog-cors

RUN go mod download
RUN VERSION=`cat ./version.json| grep Version | sed 's/"//g' | sed 's/  Version: //g'` \
    && cd cmd/example \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux go build -ldflags \
    "-X github.com/antonio-alexander/go-blog-cors/internal/internal.Version=$VERSION \
    -X github.com/antonio-alexander/go-blog-cors/internal/internal.GitCommit=$GIT_COMMIT \
    -X github.com/antonio-alexander/go-blog-cors/internal/internal.GitBranch=$GIT_BRANCH" \
    -o go-blog-cors \
    && chmod +X go-blog-cors \
    && tar -czvf go-blog-cors.tar.gz ./go-blog-cors

FROM --platform=${PLATFORM} alpine:latest

ARG GIT_BRANCH
ARG GIT_COMMIT

WORKDIR /app

COPY --from=build /go/src/go-blog-cors/cmd/example/go-blog-cors.tar.gz /app/go-blog-cors.tar.gz

RUN apk add curl

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-blog-cors

EXPOSE 8080 2345

CMD tar -xzf go-blog-cors.tar.gz && ./go-blog-cors

