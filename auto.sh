#!/bin/bash

echo "kill & cp & run ok"

sudo killall -9 detection

sleep 10 


sudo cp -f detection /home/pi/rserver
sudo cp -f appserver.conf  /home/pi/rserver

sudo cd /home/pi/rserver

sudo nohup /home/pi/rserver/detection -c  appserver_prod.conf > /home/pi/rserver/server.log 2>&1 &

echo "kill & cp & run ok"