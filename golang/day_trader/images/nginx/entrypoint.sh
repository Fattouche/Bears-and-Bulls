#!/bin/bash

cd /service

bash ./wait-for-it.sh daytrader_web:41000 --timeout=300

#Uncomment to scale up
#bash ./wait-for-it.sh day_trader_web_1:41000 --timeout=300
#bash ./wait-for-it.sh day_trader_web_2:41000 --timeout=300
#bash ./wait-for-it.sh day_trader_web_3:41000 --timeout=300

exec nginx -g 'daemon off;'
