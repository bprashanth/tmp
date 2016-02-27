```console
$ cat /etc/lsb-release
DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=15.04
DISTRIB_CODENAME=vivid
DISTRIB_DESCRIPTION="Ubuntu 15.04"

$ sudo cat /lib/systemd/system/test.service
[Unit]
Description=test

[Service]
WatchdogSec=5s
ExecStart=/tmp/test
Restart=always
```
```c
$ cat test.c
#include "systemd/sd-daemon.h"
#include <fcntl.h>
#include <time.h>
#include <stdio.h>

int main ()
{
        while(1) {
            printf("ping\n");
            sd_notify(0, "WATCHDOG=1");
            fflush(stdout);
            sleep(20);
        }
        return 0;
}
```
```bash
$ sudo apt-get install libsystemd-dev
$ gcc -o test ./test.c -lsystemd
$ ./test
ping
$ sudo systemctl enable /lib/systemd/system/test.service
$ sudo systemctl start test.service
$ journalctl -f -u test.service
-- Logs begin at Sun 2016-02-07 21:30:17 UTC. --
Feb 27 03:47:19 webmaster systemd[1]: Starting test...
Feb 27 03:47:19 webmaster test[28668]: ping
Feb 27 03:47:25 webmaster systemd[1]: test.service watchdog timeout (limit 5s)!
Feb 27 03:47:25 webmaster systemd[1]: test.service: main process exited, code=dumped, status=6/ABRT
Feb 27 03:47:25 webmaster systemd[1]: Unit test.service entered failed state.
Feb 27 03:47:25 webmaster systemd[1]: test.service failed.
Feb 27 03:47:25 webmaster systemd[1]: test.service holdoff time over, scheduling restart.
Feb 27 03:47:25 webmaster systemd[1]: Started test.
Feb 27 03:47:25 webmaster systemd[1]: Starting test...
Feb 27 03:47:25 webmaster test[28672]: ping
Feb 27 03:47:30 webmaster systemd[1]: test.service watchdog timeout (limit 5s)!
Feb 27 03:47:30 webmaster systemd[1]: test.service: main process exited, code=dumped, status=6/ABRT
Feb 27 03:47:30 webmaster systemd[1]: Unit test.service entered failed state.
Feb 27 03:47:30 webmaster systemd[1]: test.service failed.
Feb 27 03:47:30 webmaster systemd[1]: test.service holdoff time over, scheduling restart.
Feb 27 03:47:30 webmaster systemd[1]: Started test.
Feb 27 03:47:30 webmaster systemd[1]: Starting test...
Feb 27 03:47:30 webmaster test[28677]: ping

$  sudo systemctl stop test

$ journalctl -f -u test.service
-- Logs begin at Sun 2016-02-07 21:30:17 UTC. --
Feb 27 03:47:47 webmaster systemd[1]: test.service watchdog timeout (limit 5s)!
Feb 27 03:47:47 webmaster systemd[1]: test.service: main process exited, code=dumped, status=6/ABRT
Feb 27 03:47:47 webmaster systemd[1]: Unit test.service entered failed state.
Feb 27 03:47:47 webmaster systemd[1]: test.service failed.
Feb 27 03:47:47 webmaster systemd[1]: test.service holdoff time over, scheduling restart.
Feb 27 03:47:47 webmaster systemd[1]: Started test.
Feb 27 03:47:47 webmaster systemd[1]: Starting test...
Feb 27 03:47:47 webmaster test[28703]: ping
Feb 27 03:47:50 webmaster systemd[1]: Stopping test...
Feb 27 03:47:50 webmaster systemd[1]: Stopped test.
```
