#!/bin/bash

cd day_trader
#If you change the schema you need to run this
#easyjson -all main.go  

bash /startup/wait-for-it.sh daytrader_db:3306 --timeout=300

bash /startup/wait-for-it.sh daytrader_cache:11211 --timeout=300



exec "$@"
