#!/bin/bash
set -e

# https://github.com/dicksontung/docker-db2/blob/master/ibm-db2/config/entrypoint.sh

if [ -z "$DB2INST1_PASSWORD" ]; then
  echo ""
  echo >&2 'error: DB2INST1_PASSWORD not set'
  echo >&2 'Did you forget to add -e DB2INST1_PASSWORD=... ?'
  exit 1
else
  echo -e "$DB2INST1_PASSWORD\n$DB2INST1_PASSWORD" | passwd db2inst1
fi

if [ -z "$LICENSE" ];then
   echo ""
   echo >&2 'error: LICENSE not set'
   echo >&2 "Did you forget to add '-e LICENSE=accept' ?"
   exit 1
fi

if [ "${LICENSE}" != "accept" ];then
   echo ""
   echo >&2 "error: LICENSE not set to 'accept'"
   echo >&2 "Please set '-e LICENSE=accept' to accept License before use the DB2 software contained in this image."
   exit 1
fi

if [[ $1 = "db2start" ]]; then
  su - db2inst1 -c "db2start"
  nohup /usr/sbin/sshd -D 2>&1 > /dev/null &
  while true; do sleep 1000; done
elif [[ $1 = "db2sampl" ]]; then
  echo "Starting DB instance..."
  su - db2inst1 -c "db2start"
  echo "Creating sample DB..."
  su - db2inst1 -c "db2 create db sample"
  echo "Initialize complete"
  nohup /usr/sbin/sshd -D 2>&1 > /dev/null &
  while true; do sleep 1000; done
else
  exec "$1"
fi