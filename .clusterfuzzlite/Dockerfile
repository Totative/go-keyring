FROM gcr.io/oss-fuzz-base/base-builder-go:v1@sha256:68eb389f79e1b58277412ef425fab4c412cd4f61c521ce355a49c8cdd1a33e1e
COPY . $SRC/go-keyring
WORKDIR $SRC/go-keyring
COPY ./.clusterfuzzlite/build.sh $SRC/
