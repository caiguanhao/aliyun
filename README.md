oss
===

Run `./build.sh` and then enter settings and key ID and secret to build `oss`.

Run `BUILD_DOCKER=1 ./build.sh` if you want to build `oss` with Docker.

Or:

```
# store configs and keys in environment variables
for v in API_PREFIX BUCKET; do printf "$v: " && read $v && export $v; done && \
  for v in ALIYUN_ACCESS_KEY ALIYUN_ACCESS_SECRET; do printf "$v: " && read -s $v && echo && export $v; done

# build without asking
./build.sh

# clean
unset API_PREFIX BUCKET ALIYUN_ACCESS_KEY ALIYUN_ACCESS_SECRET
```

LICENSE: MIT
