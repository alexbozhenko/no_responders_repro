# no_responders_repro

Reproducer for a bug where reconnection of queue group susbscriber causes
"no respoders available" error for the clients(publishers)

Steps to reproduce:

```
# 0. 
git clone https://github.com/alexbozhenko/no_responders_repro.git
cd no_responders_repro

# 1. Start the publishers in a separate terminal, and observe the output:
cd publishers
go run  . -s "tls://eu.geo.ngs.global" -creds path/to/user.creds -n 5   subject messag


# 2. start 3(three) instances of subscribers in separate terminals
# Important detail. If number of calls to nc.QueueSubscribe=1, 
# the issue could not be reproroduced. So n must be n>=2
# The issue also could not be reproduced if publishers are connected to the same server as the subscribers.
cd subscribers
go run  . -s "tls://east.us.geo.ngs.global"  -creds  path/to/user.creds -n 2  subject queue

# 3. 
# press ctrl+c to froce reconnection in the subscribers.
# You may have to do it several times in different subscribers terminals
# Eventually , you will start to consistently see the following errors in the publisher window:

# `2025/01/16 15:05:59.414344 main.go:48: Error: nats: no responders available for request`

# this seem like a bug, since responders are actually available.

```
