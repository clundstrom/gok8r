#!/bin/bash

api(){
    echo "Building api"
    docker build ./api -t api

    if [[ $1 == '-r' ]]; then
      echo "Starting api in detached mode.."
      docker run --name api -d -t -p "8000:8000" --network gok8r-net api
    fi
}

frontend(){
    echo "Building frontend"
    docker build ./frontend -t frontend

    if [[ $1 == '-r' ]]; then
      echo "Starting frontend in detached mode.."
      docker run --name frontend -d -t -p "80:80" --network gok8r-net frontend
    fi
}

broker(){
    echo "Building event broker"
    docker build ./broker -t rabbitmq

    if [[ $1 == '-r' ]]; then
      echo "Starting broker in detached mode.."
      docker run --name rabbitmq -d -t -p "5672:5672" --network gok8r-net -p "8080:15672" rabbitmq
    fi
}

package(){
  echo "Packaging helm charts"
  DIR=$PWD
  cd "$PWD/gok8r/packages"
  helm package ../.
  cd "$DIR"
}

network_create(){
  docker network create \
        --driver=bridge \
        --subnet=10.0.0.0/16 \
        --ip-range=10.0.0.0/24 \
        --gateway=10.0.0.1 \
        gok8r-net || true
}

clear_containers(){
  docker rm -f api rabbitmq frontend || echo "No containers found. Nothing cleared."
}

if [[ $# -eq 0 ]]; then
  api
  frontend
  broker
  exit 0
fi

if [[ "$1" == '-r' ]]; then
  run='-r'
  clear_containers
  network_create
fi

for i in "$@"; do
  case "$i" in
    "-r")
    ;;
    "api")
      api $run
      ;;
    "frontend") 
      frontend $run
      ;;
    "broker")
      broker $run
      ;;
    "package")
      package
      ;;
    "deploy")
      deploy
      ;;
    *)
      echo 'One or more invalid arguments.'
      echo "Usage: $0 [OPTIONAL] {frontend|broker|api|{empty}}"
      echo 'Optional: < -r > - Build and run container'
      echo "Example: $0 -r api frontend"
      exit 1
      ;;
  esac
done





