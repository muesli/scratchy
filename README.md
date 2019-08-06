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
$ sudo scratchy bootstrap /some/root
```

```
Bootstrapping arch: /some/root
==> Creating install root at /some/root
==> Installing packages to /some/root

...

Successfully bootstrapped root: /some/root
```

Start a bash shell inside the container:

```
$ sudo scratchy run /some/root /bin/bash
```

```
Copying /etc/resolv.conf to /some/root/etc/resolv.conf
Executing in container (/some/root): /bin/bash
[root@container /]#
```

Start with a specific uid/gid inside container:

```
$ sudo scratchy -uid 1000 -gid 1000 run /some/root /bin/ps ax
```

```
Executing in container (/some/root): /bin/ps ax
  PID TTY      STAT   TIME COMMAND
    1 ?        Sl+    0:00 /proc/self/exe -uid 1000 -gid 1000 child /some/root /bin/ps ax
    6 ?        R+     0:00 /bin/ps ax
```
