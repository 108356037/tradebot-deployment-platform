#!/bin/bash

docker login -u 108356037 -p ${IMAGE_REPO_PWD}
echo ${OPENFAAS_BASIC_AUTH} | faas-cli login -u admin --password-stdin 
exec /app/main