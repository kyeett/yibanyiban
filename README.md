<div style="font-weight: bold; font-size:42px; border:10px solid black; padding:5px 10px;display: inline-block">一半<br>一半</div>
<hr>

一半一半 (pronounced _i-ban i-ban_) is a simple webserver with a REST API that validates an [International Bank Account Number (IBAN)](https://sv.wikipedia.org/wiki/International_Bank_Account_Number).

## Usage

一半一半 has a single endpoint `/validate`

```bash
GET /validate?iban=<IBAN>
```

For example, using curl:

```bash
curl localhost:8080/validate?iban=AL86751639367318444714198669
```

### More examples:

| Request                                                | Response                                                                                                                                                         |
| ------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| <span style="color: green">**Valid IBAN**</span>       |                                                                                                                                                                  |
| GET /validate?iban=AL86751639367318444714198669        | <pre>{<br> iban: "AL86751639367318444714198669",<br> valid: true,<br> message: "OK"<br>}</pre>                                                                   |
| <span style="color: red">**Invalid IBAN**</span>       |                                                                                                                                                                  |
| GET /validate?iban=GB83WEST12345698765432              | <pre>{<br> iban: "GB83WEST12345698765432",<br> valid: false,<br> message: "checksum is incorrect"<br>}</pre>                                                     |
| GET /validate?iban=XX82WEST12345698765432              | <pre>{<br> iban: "XX82WEST12345698765432",<br> valid: false,<br> message: "country code invalid"<br>}</pre>                                                      |
| GET /validate?iban=GB82WEST12345698765432\*            | <pre>{<br> iban: "GB82WEST12345698765432\*",<br> valid: false,<br> message: "Invalid characters, allowed are alphanumeric (A-Z, 0-9) and space (' ')"<br>}</pre> |
| GET /validate?iban=GB8                                 | <pre>{<br> iban: "GB8",<br> valid: false,<br> message: "IBAN is too short, expected > 4"<br>}</pre>                                                              |
| GET /validate?iban=AL86751639367318444714198669AL86751 | <pre>{<br> iban: "AL86751639367318444714198669AL86751",<br> valid: false,<br> message: "IBAN is too long, expected < 34"<br>}</pre>                              |

## Todo

- [x] Function to validate an IBAN
- [x] Basic webserver
- [x] Webserver with timeouts
- [ ] Makefile for ~~build~~, ~~test~~ and load test
- [ ] Additional IBAN numbers for testing
- [x] Additional examples
- [ ] Clean up code
- [ ] Clean up README
- [ ] Move main package to cmd
- [ ] Flag to set port

## Lessons learned

https://stackoverflow.com/questions/10971800/golang-http-server-leaving-open-goroutines/10972453#10972453

`accept tcp [::]:8080: accept: too many open files in system`

https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/


Leak Detection.
The formula for detecting leaks in webservers is to add instrumentation endpoints and use them alongside load tests.

https://blog.minio.io/debugging-go-routine-leaks-a1220142d32c


### 

### Set max number of file descriptors

The max number of file descriptors can be set using 

```bash
ulimit -n <MAX_NUM>
```

Highest allowed value on my Mac is `10240`.

```bash
╰─$ ulimit -a
-t: cpu time (seconds)              unlimited
-f: file size (blocks)              unlimited
-d: data seg size (kbytes)          unlimited
-s: stack size (kbytes)             8192
-c: core file size (blocks)         0
-v: address space (kbytes)          unlimited
-l: locked-in-memory size (kbytes)  unlimited
-u: processes                       709
-n: file descriptors                10240
```

### Monitoring number of goroutines
Find number of goroutines `runtime.NumGoroutine()`.


accept tcp [::]:8080: accept: too many open files in system


## References
https://www.youtube.com/watch?v=hVFEV-ieeew
