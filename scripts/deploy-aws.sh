#!/bin/bash
aws ecr describe-repositories --repository-names "$GITHUB_REPOSITORY" || aws ecr create-repository --repository-name "$GITHUB_REPOSITORY"
docker pull $DOCKER_REGISTRY/$GITHUB_REPOSITORY:latest || true
docker build --cache-from $DOCKER_REGISTRY/$GITHUB_REPOSITORY:latest -t $DOCKER_REGISTRY/$GITHUB_REPOSITORY:$GITHUB_SHA .
docker tag "$DOCKER_REGISTRY/$GITHUB_REPOSITORY:$GITHUB_SHA" $DOCKER_REGISTRY/$GITHUB_REPOSITORY:latest
docker push "$DOCKER_REGISTRY/$GITHUB_REPOSITORY" -a