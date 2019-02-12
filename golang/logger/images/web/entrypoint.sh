#!/bin/bash

bash /startup/wait-for-it.sh db:3306 --timeout=300

exec "$@"
