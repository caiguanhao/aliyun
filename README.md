aliyun
======

Command-line tool for [Aliyun Cloud Services](http://www.aliyun.com/product/).

[![Circle CI](https://circleci.com/gh/caiguanhao/aliyun.svg?style=svg)](https://circleci.com/gh/caiguanhao/aliyun)

USAGE
-----

### ECS

```help
NAME:
   ecs - control Aliyun ECS instances

USAGE:
   ecs [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   list-instances, list, ls, l          list all ECS instances of all regions
   list-images, images, i               show info of all images
   list-regions, regions, n             list all available regions and zones
   list-instance-types, types, t        list all instance types
   list-security-groups, groups, g      list all security groups
   create-instance, create, c           create an instance
   allocate-public-ip, allocate, a      allocate an IP address for an instance
   start-instance, start, s             start an instance
   stop-instance, stop, S               stop an instance
   restart-instance, restart, r         restart an instance
   remove-instance, remove, rm, R       remove an instance
   update-instance, update, u           update attributes of an instance
   hide-instance, hide, h               hide instance from instance list
   unhide-instance, unhide, H           un-hide instance from instance list
   monitor-instance, monitor, m         show CPU and network usage history of an instance

GLOBAL OPTIONS:
   --quiet, -q		show only name or ID
   --verbose, -V	show more info
   --version, -v	print the version
```

### OSS

```help
NAME:
   oss - control Aliyun OSS

USAGE:
   oss [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   upload, up, u, put           upload local files to remote OSS
   download, down, dl, d, get   get remote OSS files to local
   list, ls, l                  show list of files on remote OSS
   diff                         show different files on local and remote OSS

GLOBAL OPTIONS:
   --bucket, -b "xxxxxxxxx"                                     bucket name
   --prefix, -p "https://%s.oss-cn-hangzhou.aliyuncs.com"       API prefix
   --key "xxxxxxxxxxxxxxxx"                                     access key [$ACCESS_KEY]
   --secret                                                     access key secret [$ACCESS_SECRET]
   --concurrency, -c "4"                                        job concurrency, defaults to number of CPU (4), max is 16
   --dry-run, -D                                                do not actually run
   --verbose, -V                                                show more info
   --generate-bash-completion
   --version, -v                                                print the version
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
