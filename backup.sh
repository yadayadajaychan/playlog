#!/usr/bin/env bash

if [ "$#" -ne 2 ]
then
	echo 'usage: backup.sh <playdb> <backup_dir>' >&2
	exit 1
fi

sqlite3 "$1" ".backup \"$2/$1.$(date +%s).bak\""
