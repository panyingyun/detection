#!/bin/bash

echo "kill & cp & run ok"

sudo killall -9 detection

sudo cp -f detection /home/pi/rserver
sudo cp -f appserver.conf  /home/pi/rserver

sudo nohup /home/pi/rserver/detection -c  /home/pi/rserver/appserver.conf > /home/pi/rserver/server.log 2>&1 &

echo "kill & cp & run ok"