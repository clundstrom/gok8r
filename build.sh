#!/bin/bash

api(){
    echo "Building api"
    docker build ./api -t api

    if [[ $1 == '-r' ]]; then
    echo "Starting api in detached mode.."
    docker run --name api -d -t -p "8000:8000" api 
    fi
}

frontend(){
    echo "Building frontend"
    docker build ./frontend -t frontend

    if [[ $1 == '-r' ]]; then
    echo "Starting frontend in detached mode.."
    docker run --name frontend -d -t -p "80:80" frontend
    fi
}

broker(){
    echo "Building event broker"
    docker build ./rabbitmq -t rabbitmq

    if [[ $1 == '-r' ]]; then
    echo "Starting frontend in detached mode.."
    docker run --name rabbitmq -d -t -p "5672:5672" rabbitmq
    fi
}

package(){
  echo "Packaging helm charts"
  DIR=$PWD
  cd "$PWD/gok8r/packages"
  helm package ../.
  cd "$DIR"
}

deploy(){
  helm install gok8r --generate-name
}



if [[ $# -eq 0 ]]; then
  api "-r"
  frontend "-r"
  broker "-r"
  exit 0
fi

for i in "$@"; do

  if [[ "$1" == '-r' ]]; then
    run='-r'
  fi

  case "$i" in
    "-r");;
    "api")
      api $run
      ;;
    "frontend") 
      frontend $run
      ;;
    "package") 
    package
    ;;
    "deploy") 
    deploy
    ;;
    *)
      echo 'One or more invalid arguments.'
      echo "Usage: $0 [OPTIONAL] {db|frontend|broker|charts|{empty}}"
      echo 'Optional: < -r > - Build and run container >'
      echo "Example: $0 -r api frontend"
      exit 1
      ;;
  esac
done





