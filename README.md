## HTTP log monitoring console program

Consumes an actively written-to w3c-formatted HTTP access log (https://www.w3.org/Daemon/User/Config/Logging.html). It should default to reading /tmp/access.log and be overrideable

Example log lines:

```
127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123
127.0.0.1 - jill [09/May/2018:16:00:41 +0000] "GET /api/user HTTP/1.0" 200 234
127.0.0.1 - frank [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 200 34
127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 503 12
```

1. Display stats every 10s about the traffic during those 10s: the sections of the web site with the most hits, as well as interesting summary statistics on the traffic as a whole. A section is defined as being what's before the second '/' in the resource section of the log line. For example, the section for "/pages/create" is "/pages"
2. Make sure a user can keep the app running and monitor the log file continuously
3. Whenever total traffic for the past 2 minutes exceeds a certain number on average, add a message saying that “High traffic generated an alert - hits = {value}, triggered at {time}”. The default threshold should be 10 requests per second, and should be overridable.
4. Whenever the total traffic drops again below that value on average for the past 2 minutes, add another message detailing when the alert recovered.
5. Make sure all messages showing when alerting thresholds are crossed remain visible on the page for historical reasons.

## Writer
I wrote a helper that generates log files at https://github.com/veverkap/logtop/blob/master/writer/writer.go

```
Usage of ./writer:
  -file string
    	Location of log file (default "/tmp/access.log")
  -rate int
    	Number of requests per second to write (default 10)
```

## Reader

The reader lives at https://github.com/veverkap/logtop/blob/master/reader/reader.go

```
Usage of ./reader:
  -logFileLocation string
    	Location of log file to parse (default "/tmp/access.log")
  -threshold int
    	Number of requests per second maximum for alert (default 10)
  -thresholdDuration int
    	Duration in seconds of sampling period for alerts (default 120)
```
