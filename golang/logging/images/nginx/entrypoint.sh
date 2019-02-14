#!/bin/bash

cd /service

bash ./wait-for-it.sh logging_server:40000 --timeout=300

#Uncomment to scale up
#bash ./wait-for-it.sh logging_1:40000 --timeout=300
#bash ./wait-for-it.sh logging_2:40000 --timeout=300
#bash ./wait-for-it.sh logging_3:40000 --timeout=300

exec nginx -g 'daemon off;'
