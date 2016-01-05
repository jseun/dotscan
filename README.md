# A Simple Network Port Scanner for DigitalOcean folks

Hi! Welcome to my intimate attempt at getting a job at DigitalOcean.
As part of the interview process, I have been asked to build an efficient port scanner with the following requirements.

 * Ability to scan a single host with TCP from 1 to 65535.
 * Write quality code, as in production, and test it.

Additionally, the following features would earn me good points.

 * Concurrency in the port scanning process
 * Support for IP version 6 hosts
 * Support for UDP scan
 * Map port to services file for service guess
 * Ability to scan blocks of hosts

This project contains Go source to build an efficient port scanner.

## Wait, Sam! What features did you implement so far?

While I tried to keep that project as complete as possible yet delivering it as fast as possible, this is the features I could get done in 3 days of work.

 * Concurrency in the port scanning process
 * Support for IP version 6 hosts
 * Ability to scan blocks of hosts
 * Map port to services file for service guess

Also, I planned to support UDP port scan from the beginning but UDP protocol being stateless requires a way more complete port scan technique than TCP.  For that reason, I focused on writing good and simple code.

## And... Did you get the job?

Yes! I'll be working as a software engineer at DO starting on the 25th of January 2016.

## So what about this project?  Is it over?

Well, it is over for me.  I won't be adding new features, fixing bugs or maintaining this project at all.  However, if you feel a bit challenged about this project, you might want to look at the following tasks.

### 1. Fix UDP port scan

UDP port scan requires some payload to be sent at the host and there is no guarantee that a given dummy payload will make the remote host answer even if there's a service running at that port.

### 2. Write more tests

I love to write tests in Golang.  Other than the unit tests that can be found in `scanner/host_test.go` to cover 100% of the NewHosts() function, and TestRunAndWaitLocalhostTCP4 and TestRunAndWaitLocalhostTCP6 in `scanner/scanner_test.go`, a lot of other tests could be written.

### 3. Hunt the bugs down!

When scanning a host that does not exist, I've seen some erroneous results like the port x is open which cannot be real unless my network is haunted by ghost hosts. I guess I should call the ghost hosts buster.

# Enough.  How do I run the software?

OK.  Here's the deal.  Do you want a blazing fast port scan?  If the answer is yes, read on.

Else, you may go directly to the section called *I don't have time for configuration*.

In any case, you would need to install `doscan` using `go`.
```
$ go get github.com/jseun/doscan
```

## Run the software as root

What you need to realize is that if we can concurrently open 65535 ports to scan a host, it would take seconds to get results.

Here's how I achieved my best shots.

Running 20480 workers, each of them dialing a port for maximum 10 seconds, given that there is 65535 ports to dial, it would take around 3 rounds to cover them all.  Which gives us an execution time of about 30 seconds. Add some more if you need to resolve a hostname first.

```
$ sudo -s
# ulimit -n 100000
# doscan -w 20480 -t 10s google.com
Simple Network Port Scanner version 0.0.0

Copyright 2015 Samuel Jean. All rights reserved.
Use of this source code is governed by a Simplified BSD
license that can be found in the LICENSE file.

	206.126.112.187: 80 443.

Scan completed in 40.147542483s
```

However, if you are scanning a host on your local network, a timeout of 5s should be enough, thus reducing execution time by a power of 2.

```
# dotscan -w 20480 -t 5s 192.168.0.1
Simple Network Port Scanner version 0.0.0

Copyright 2015 Samuel Jean. All rights reserved.
Use of this source code is governed by a Simplified BSD
license that can be found in the LICENSE file.

	192.168.0.1: 80.

Scan completed in 20.094733178s
```

This is my wifi router.

## I don't have time for configuration

No problem, mate.  We will take it easy.
```
$ dotscan localhost
Simple Network Port Scanner version 0.0.0

Copyright 2015 Samuel Jean. All rights reserved.
Use of this source code is governed by a Simplified BSD
license that can be found in the LICENSE file.

	127.0.0.1: 22 631 6379.

Scan completed in 1.310837447s
```

Scanning localhost with the default parameters will give you fast results.

# Conclusion

Thank you for reading this far.  I've been very enthusiast about this homework and I am glad it gave me a chance to get a job at DigitalOcean.
