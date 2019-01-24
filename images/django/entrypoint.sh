#!/bin/bash

bash /startup/wait-for-it.sh db:3306 --timeout=300

bash /startup/wait-for-it.sh cache:11211 --timeout=300

python ./day_trader/manage.py makemigrations
python ./day_trader/manage.py migrate

exec "$@"
