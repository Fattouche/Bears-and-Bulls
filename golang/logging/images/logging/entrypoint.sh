#!/bin/bash

bash /startup/wait-for-it.sh logging_db:3306 --timeout=300

exec "$@"
