aliyun
======

Command-line tool for [Aliyun Cloud Services](http://www.aliyun.com/product/).

[![Circle CI](https://circleci.com/gh/caiguanhao/aliyun.svg?style=svg)](https://circleci.com/gh/caiguanhao/aliyun)

USAGE
-----

ECS:

```help
NAME:
   ecs - control Aliyun ECS instances

USAGE:
   ecs [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   list-instances, list, l		list all ECS instances of all regions
   list-images, images, i		show info of all images
   list-regions, regions, n		list all available regions
   list-instance-types, types, t	list all instance types
   list-security-groups, groups, g	list all security groups
   create-instance, create, c		create an instance
   allocate-public-ip, allocate, a	allocate an IP address for an instance
   start-instance, start, s		start an instance
   stop-instance, stop, k		stop an instance
   restart-instance, restart, r		restart an instance
   remove-instance, remove, d		remove an instance
   update-instance, update, e		update attributes of an instance
   hide-instance, hide, h		hide instance from instance list
   unhide-instance, unhide, u		un-hide instance from instance list

GLOBAL OPTIONS:
   --quiet, -q		show only name or ID
   --verbose, -V	show more info
   --version, -v	print the version
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
