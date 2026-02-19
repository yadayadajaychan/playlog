#!/usr/bin/env sh

if [ $# -ne 2 ]
then
	echo "usage: $(basename "$0") <SEGAID> <PASSWORD>" >&2
	exit 1
fi

SEGAID="$1"
PASSWORD="$2"

curl 'https://lng-tgk-aime-gw.am-all.net/common_auth/login?site_id=maimaidxex&redirect_url=https://maimaidx-eng.com/maimai-mobile/&back_url=https://maimai.sega.com/' \
  -c cookies > /dev/null

sleep 1

curl 'https://lng-tgk-aime-gw.am-all.net/common_auth/login/sid/' \
  -X POST \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-raw 'retention=1&sid='"$SEGAID"'&password='"$PASSWORD" \
  -b cookies \
  -c cookies \
  -L > /dev/null

sleep 1

curl 'https://maimaidx-eng.com/maimai-mobile/home/' \
  -b cookies \
  -c cookies \
  -L > /dev/null

sleep 1

curl 'https://maimaidx-eng.com/maimai-mobile/record/' \
  -b cookies \
  -c cookies
