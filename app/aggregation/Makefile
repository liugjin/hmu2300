BUILD_VERSION   := v1.2.0-beta14
BUILD_TIME      := $(shell date "+%F %T")
# BUILD_NAME      := app_$(shell date "+%Y%m%d%H" )
BUILD_NAME      := aggregation
SOURCE          := ./*.go
TARGET_DIR      := .
COMMIT_SHA1     := $(shell git rev-parse HEAD )

all:
	go build -ldflags                           \
	"                                           \
	-X 'clc.hmu/app/public/appver.BuildVersion=${BUILD_VERSION}'     \
	-X 'clc.hmu/app/public/appver.BuildTime=${BUILD_TIME}'       \
	-X 'clc.hmu/app/public/appver.BuildName=${BUILD_NAME}'       \
	-X 'clc.hmu/app/public/appver.CommitID=${COMMIT_SHA1}'       \
	-s										\
	-w										\
	"                                           \
	-o ${BUILD_NAME} ${SOURCE}

clean:
	rm ${BUILD_NAME} -f

install:
	mkdir -p ${TARGET_DIR}
	cp ${BUILD_NAME} ${TARGET_DIR} -f

.PHONY: all clean install ${BUILD_NAME}
