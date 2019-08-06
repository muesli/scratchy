scratchy
========

Quickly bootstrap a Linux distro in a (non-Docker) container and interactively
execute something in it.

Note: this is early stage and automatic bootstrapping only supports ArchLinux
as of now.

## Installation

### From Source

Make sure you have a working Go environment (Go 1.11 or higher is required).
See the [install instructions](http://golang.org/doc/install.html).

Compiling scratchy is easy, simply run:

    git clone https://github.com/muesli/scratchy.git
    cd scratchy
    go build

## Usage

Bootstrap a new base ArchLinux install into a directory:
(make sure package `arch-install-scripts` is installed on your host!)

```
$ sudo scratchy bootstrap /some/new/root
```

```
Bootstrapping arch: /some/new/root
==> Creating install root at /some/new/root
==> Installing packages to /some/new/root

...

Successfully bootstrapped root: /some/new/root
```

Start a bash shell inside the container:

```
$ sudo scratchy run /some/new/root /bin/bash
```

```
Copying /etc/resolv.conf to /some/new/root/etc/resolv.conf
Executing in container (/some/new/root): /bin/bash
[root@container /]#
```

Start with a specific uid/gid inside container:

```
$ sudo scratchy -uid 1000 -gid 1000 run ./root /bin/ps aux
```

```
Executing in container (/some/new/root): /bin/ps aux
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.0 103048  1428 ?        Sl+  01:14   0:00 /proc/self/exe -uid 1000 -gid 1000 child ./root /bin/ps aux
1000         6  0.0  0.0   9724  2804 ?        R+   01:14   0:00 /bin/ps aux
```
