#!/bin/bash

bash ./wait-for-it.sh queue:6379 --timeout=300

python ./day_trader/manage.py rqworker default
