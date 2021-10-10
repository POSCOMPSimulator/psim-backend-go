#!/bin/sh
echo $DATABASE_URL
goose -dir database/migrations postgres $DATABASE_URL up