#!/usr/bin/env bash

if [ $# -ne 1 ]
then
	echo "usage: $(basename "$0") <idx>" >&2
	exit 1
fi

IDX=$(echo -n "$1" | sed 's/,/%2C/')

curl "https://maimaidx-eng.com/maimai-mobile/record/playlogDetail/?idx=$IDX" \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0' \
  -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' \
  -H 'Accept-Language: en-US,en;q=0.9' \
  -H 'Accept-Encoding: gzip, deflate, br, zstd' \
  -H 'Connection: keep-alive' \
  -H 'Upgrade-Insecure-Requests: 1' \
  -H 'Sec-Fetch-Dest: document' \
  -H 'Sec-Fetch-Mode: navigate' \
  -H 'Sec-Fetch-Site: cross-site' \
  -H 'Sec-Fetch-User: ?1' \
  -H 'DNT: 1' \
  -H 'Sec-GPC: 1' \
  -H 'Priority: u=0, i' \
  -H 'Pragma: no-cache' \
  -H 'Cache-Control: no-cache' \
  -b cookies \
  -c cookies
