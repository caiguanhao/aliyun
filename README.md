aliyun
======

Command-line tool for [Aliyun Cloud Services](http://www.aliyun.com/product/).

[![Circle CI](https://circleci.com/gh/caiguanhao/aliyun.svg?style=svg)](https://circleci.com/gh/caiguanhao/aliyun)

USAGE
-----

ECS:

```help
Usage: ecs [OPTION] [ACTION] [TARGET]

Using Access Key xxxxxxxxxxxxxxxx

Action                Alias        Description
list-instances        l, list      List all instances, show one if ID is specified
list-images           i, images    List all images
list-regions          n, regions   List all regions
list-instance-types   t, types     List all instance types
list-security-groups  g, groups    List all security groups
create-instance       c, create    Create an instance
allocate-public-ip    a, allocate  Allocate an IP address for an instance
start-instance        s, start     Start an instance
stop-instance         k, stop      Stop an instance
restart-instance      r, restart   Restart an instance
remove-instance       d, remove    Remove an instance
update-instance       e, update    Update attributes of an instance
hide-instance         h, hide      Hide instance from instance list
unhide-instance       u, unhide    Un-hide instance from instance list
```

To upload files:

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
Source: https://github.com/caiguanhao/aliyun
```

To get list of different files:

```help
oss-diff [OPTION] LOCAL-DIR  REMOTE-DIR
                  LOCAL-FILE REMOTE-FILE

Options:
    -r, --reverse  Print LOCAL file paths to stderr, REMOTE to stdout

    -m, --md5      Verify MD5 checksum besides file name and size
    -s, --shhh     Show only file path

Status code: 0 - local and remote are identical
             1 - local has different files
             2 - remote has different files
             3 - both local and remote have different files
```

To get a file from OSS:

```help
oss-get [--curl] REMOTE-FILE [LOCAL-FILE]

    --curl    generate curl script
```

BUILD
-----

Run `./build.sh` and then enter configs, key ID and secret to start.

If you are on Mac OS X and you want to build a Linux version,
you can run `BUILD_DOCKER=1 ./build.sh` to build in a Docker container.
You need to install `boot2docker` and `docker-compose`.

To continously build, you can also set environment variables:
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
