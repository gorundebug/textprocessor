#!/bin/bash

CMD="$1"
CONTAINER_NAME="textprocessor_${CMD%.*}"
IMAGE_NAME="textprocessor"
WORK_DIR="/textprocessor"

if [ "$2" == "reset" ]; then
    if [ "$(docker ps -a -q -f name="${CONTAINER_NAME}")" ]; then
        echo "Remove container ${CONTAINER_NAME}"

        if ! docker rm -f "${CONTAINER_NAME}"; then
            echo "docker rm -f ${CONTAINER_NAME}"
            exit 1
        fi
    fi
    if [ "$(docker images -q ${IMAGE_NAME})" ]; then
        echo "Remove image ${IMAGE_NAME}"
        if ! docker rmi "$IMAGE_NAME"; then
            echo "docker rmi ${IMAGE_NAME} failed"
            exit 1
        fi
    fi
fi

if [ "$(docker ps -a -q -f name="${CONTAINER_NAME}")" ]; then
    echo "Container ${CONTAINER_NAME} exists. Starting it..."
    if ! docker start -i "${CONTAINER_NAME}"; then
      exit 1
    fi
else
    echo "Container ${CONTAINER_NAME} does not exist. Creating and starting it..."
    if ! bash build_docker.sh; then
        exit 1
    fi

    if ! docker run -it -v "$(pwd)":${WORK_DIR} -w ${WORK_DIR} -e CMD="${CMD}" --name "${CONTAINER_NAME}" ${IMAGE_NAME}; then
      exit 1
    fi
fi

exit 0