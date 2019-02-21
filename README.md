# 一半<br>一半

一半一半 (pronounced _i-ban i-ban_) is a simple webserver with a REST API that validates an [International Bank Account Number (IBAN)](https://sv.wikipedia.org/wiki/International_Bank_Account_Number). It has a single endpoint `/validate`

```bash
GET /validate?iban=<IBAN>
```

The server tested for 2000 connected users and 20k queries per second (<a href="#load-test-results">details here</a>)

## Usage

1. Install and start server

```bash
> go get github.com/kyeett/yibanyiban/cmd/yibanyiban
> $GOBIN/yibanyiban -port 8080
Serving IBAN validation service on :8080
```

2. Validate your IBAN
 - [localhost:8080/validate?iban=AL86751639367318444714198669](localhost:8080/validate?iban=AL86751639367318444714198669)
 - `curl localhost:8080/validate?iban=AL86751639367318444714198669`

### More examples:

| Request                                                | Response                                                                                                                                                         |
| ------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| <span style="color: green">**Valid IBAN**</span>       |                                                                                                                                                                  |
| GET /validate?iban=AL86751639367318444714198669        | <pre>{<br> iban: "AL86751639367318444714198669",<br> valid: true,<br> message: "OK"<br>}</pre>                                                                   |
| <span style="color: red">**Invalid IBAN**</span>       |                                                                                                                                                                  |
| GET /validate?iban=GB83WEST12345698765432              | <pre>{<br> iban: "GB83WEST12345698765432",<br> valid: false,<br> message: "checksum is incorrect"<br>}</pre>                                                     |
| GET /validate?iban=GB82WEST12345698765432\*            | <pre>{<br> iban: "GB82WEST12345698765432\*",<br> valid: false,<br> message: "Invalid characters, allowed are alphanumeric (A-Z, 0-9) and space (' ')"<br>}</pre> |
| GET /validate?iban=GB8                                 | <pre>{<br> iban: "GB8",<br> valid: false,<br> message: "IBAN is too short, expected > 4"<br>}</pre>                                                              |
| GET /validate?iban=AL86751639367318444714198669AL86751 | <pre>{<br> iban: "AL86751639367318444714198669AL86751",<br> valid: false,<br> message: "IBAN is too long, expected < 34"<br>}</pre>                              |


<br><br><br>

# Miscellaneous

## Todo

- [x] Function to validate an IBAN
- [x] Basic webserver
- [x] Webserver with timeouts
- [x] Makefile for ~~build~~, ~~test~~ and ~~load~~ test
- [x] Additional IBAN numbers for testing
- [x] Additional examples
- [x] Clean up code
- [x] Clean up README
- [x] Move main package to cmd
- [x] Flag to set port

## Load test results

Running load test with 2000 of concurrent clients

`QPS` = Queries per second

| Target QPS | Actual QPS   | Total requests | Errors |
| ---------- | ------------ | -------------- | ------ |
| 1000       | 999.936652   | 30009          | 0      |
| 5000       | 4606.634259  | 138249         | 0      |
| 10000      | 9991.643456  | 299859         | 10     |
| 15000      | 14991.572884 | 449907         | 94     |
| 20000      | 19952.290100 | 598862         | 102    |
| 25000      | 23339.391439 | 700569         | 20     |
| 30000      | 26364.486095 | 791647         | 0      |

## Lessons learned

### Big ints
The integer generated from the IBAN is bigger than int32 and float64. Go supports handling arbitrarily big numbers using the [math/big](https://golang.org/pkg/math/big/) package. This create an big integer from the `stringInt` with base 10.

```go
numeric, _ := new(big.Int).SetString(stringInt, 10)
```

Then you can perform most other operations that can be done on regular ints, but using the `big.Int` methods.

```go
    func (z *Int) Mod(x, y *Int) *Int
    func (z *Int) ModInverse(g, n *Int) *Int
    func (z *Int) ModSqrt(x, p *Int) *Int
    func (z *Int) Mul(x, y *Int) *Int
    func (z *Int) Xor(x, y *Int) *Int
```

**Note**

There is an algorithm for calculating `mod 97` (used for IBAN check number), that can be used if `go` didn't support big numbers. If I would need to optimize further, this could be an option (see [Wikipedia: Validating_the_IBAN](https://en.wikipedia.org/wiki/International_Bank_Account_Number#Validating_the_IBAN)).

#### References

- [medium: Big integers in Go](https://medium.com/orbs-network/big-integers-in-go-14534d0e490d)
- [golang.org: math/big](https://golang.org/pkg/math/big/)
- [wiki: Validating_the_IBAN](https://en.wikipedia.org/wiki/International_Bank_Account_Number#Validating_the_IBAN)

### Load testing
I used `fasthttploader`([github](https://github.com/hagen1778/fasthttploader)) to load test my application.

#### Problems

```bash
accept tcp [::]:8080: accept: too many open files in system
```

Since both the server and load test client was running on the same, they both share cpu resources, as well as file descriptors. Obviously, it is not ideal to share for server and client, but a few ways of observing what is happening:

#### 1. Monitoring number of goroutines

```golang
fmt.Printf("num goroutines: %d\n", runtime.NumGoroutine())
```

#### 2. Change max number of file descriptors

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

#### 3. Change number of concurrent connections from client

```bash
fasthttploader -c <parallel connections> ...
```
