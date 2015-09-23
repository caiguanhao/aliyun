oss
===

Command-line tool to upload files to Aliyun's Open Storage Service ([OSS](http://www.aliyun.com/product/oss)).

USAGE
-----

```help
oss [OPTION] SOURCE ... TARGET

Options:
    -c <num>   Specify how many files to process concurrently, default is 2, max is 10

    -b <name>  Specify bucket other than: my-bucket
    -z <url>   Specify API URL prefix other than: https://%s.oss-cn-hangzhou.aliyuncs.com
       You can use custom domain or official URL like this:
       {http, https}://%s.oss-cn-{beijing, hangzhou, hongkong, qingdao, shenzhen}{, -internal}.aliyuncs.com
       Note: %s will be replaced with the bucket name if specified

    --parents  Use full source file name under TARGET

    -v  Be verbosive
    -d  Dry-run. See list of files that will be transferred,
        show full URL if -v is also set

Built with key ID abcdefghijklmnop on 2015-08-19 11:08:01 (8b72aaf)
API: https://my-bucket.oss-cn-hangzhou.aliyuncs.com
Source: https://github.com/caiguanhao/oss
```

BUILD
-----

Run `./build.sh` and then enter configs and key ID and secret to build `oss`.

If you are on Mac OS X and you want to build a Linux version,
you can run `BUILD_DOCKER=1 ./build.sh` to build `oss` in a Docker container.
You need to install `boot2docker` and `docker-compose`.

To continously build `oss`, you can also set environment variables:
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
