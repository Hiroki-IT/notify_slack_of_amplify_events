#===================
# Global ARG
#===================
ARG GO_VERSION="1.16.5"
ARG LABEL="Hiroki <hasegawafeedshop@gmail.com>"

#===================
# Build Stage
#===================
FROM golang:${GO_VERSION} as build

WORKDIR /go/src

# インストールのキャッシュを活用するためにコピーしておく．
COPY go.mod go.sum ./

# ライブラリをインストールする．
RUN go get github.com/cosmtrek/air@v1.15.1

COPY . .

# ビルドのアーティファクトを/go/binに配置する．
# NOTE: netパッケージは標準で動的リンクのため，静的リンクを明示的に指定する必要がある．
RUN go build -x -a -tags netgo -installsuffix netgo -o /go/bin ./cmd

#===================
# Base Stage
#===================
FROM golang:${GO_VERSION}-alpine as base

WORKDIR /go/src

# /go/binにパスを通す．
ENV PATH $PATH:/go/bin

COPY --from=build /go ../

#===================
# Develop Stage
#===================
FROM base as dev
LABEL maintainer=${LABEL}

CMD ["/go/bin/cmd"]

#===================
# Staging Stage
#===================
FROM base as stg
LABEL maintainer=${LABEL}

CMD ["/go/bin/cmd"]

#===================
# Production Stage
#===================
FROM base as prd
LABEL maintainer=${LABEL}

CMD ["/go/bin/cmd"]
