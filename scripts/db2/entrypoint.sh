#!/bin/bash
su - db2inst1 -c "db2start"
su - db2inst1 -c "db2sampl"
while true; do sleep 1000; done