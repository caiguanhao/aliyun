#!/bin/bash

set -e

function str_to_array {
  eval "local input=\"\$$1\""
  input="$(echo "$input" | awk '
  {
    split($0, chars, "")
    for (i = 1; i <= length($0); i++) {
      if (i > 1) {
        printf(", ")
      }
      printf("'\''%s'\''", chars[i])
    }
  }
  ')"
  eval "$1=\"$input\""
}

function update_access_key {
  str_to_array BUCKET
  str_to_array ALIYUN_ACCESS_KEY
  str_to_array ALIYUN_ACCESS_SECRET
  awk "
  /DEFAULT_BUCKET/ {
    print \"var DEFAULT_BUCKET = []byte{${BUCKET}}\"
    next
  }
  /KEY/ {
    print \"var KEY = []byte{${ALIYUN_ACCESS_KEY}}\"
    next
  }
  /SECRET/ {
    print \"var SECRET = []byte{${ALIYUN_ACCESS_SECRET}}\"
    next
  }
  {
    print
  }
  " access.go > _access.go

  mv _access.go access.go
}

if test -z "$BUCKET"; then
  echo -n "Please enter default bucket name: "
  read BUCKET
fi
if test -z "$ALIYUN_ACCESS_KEY"; then
  echo -n "Please paste your access key ID: (will not be echoed) "
  read -s ALIYUN_ACCESS_KEY
  echo
fi
if test -z "$ALIYUN_ACCESS_SECRET"; then
  echo -n "Please paste your access key SECRET: (will not be echoed) "
  read -s ALIYUN_ACCESS_SECRET
  echo
fi
update_access_key

go build

BUCKET="bucket"
ALIYUN_ACCESS_KEY="key"
ALIYUN_ACCESS_SECRET="secret"
update_access_key
