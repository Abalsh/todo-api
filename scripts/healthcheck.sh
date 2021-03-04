#!/bin/sh

wget --quiet --tries=1 --spider http://127.0.0.1:8080/todo-api-health || exit 1