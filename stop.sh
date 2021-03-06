#!/bin/sh

startwork() {
  rm -rf $PWD/log/*
}

stopwork() {
  SERVERLIST='rcenterserver logicserver roomserver'

  for serv in $SERVERLIST; do
    echo "stop $serv"
    ps aux | grep "/$serv" | sed -e "/grep/d" | awk '{print $2}' | xargs kill -9 2 &>/dev/null
  done

  echo "running server:"$(ps x | grep "server -c" | sed -e '/grep/d' | wc -l)
}

case $1 in
stop)
  stopwork
  ;;
start)
  startwork
  ;;
*)
  stopwork
  sleep 1s
  startwork
  ;;
esac
