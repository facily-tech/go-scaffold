#!/bin/sh

# Importing env from secrets
ENV_FILE=/vault/secrets/.env
if [ -f "$ENV_FILE" ]; then
    echo "INFO: exporting env file from vault"
    set -o allexport # enforce export
    . $ENV_FILE # source with export envs
    set +o allexport # turn back to not enforce export
else
    echo "ERROR: vault secrets env not found int $ENV_FILE"
fi

if ! [ -z ${GOOGLE_APP_CREDENTIALS+x} ];
then
    echo ${GOOGLE_APP_CREDENTIALS} | base64 -d > /app/credentials.json
    export GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json
fi

export KUBE_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)

# Running 
RUN=$@
$RUN