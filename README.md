# 一半一半 
一半一半 (pronounced _i-ban i-ban_) is a simple webserver with a REST API that validates an [International Bank Account Number (IBAN)](https://sv.wikipedia.org/wiki/International_Bank_Account_Number).

## Usage
一半一半 has a single endpoint, `/validate`

```bash
GET /validate?iban=<IBAN>
```

For example, using curl:
```bash
curl localhost:8080/validate?iban=AL86751639367318444714198669
```


### More examples:


| Request                                          | Response                                                                                                                                                                 |
|--------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| <span style="color: green">**Valid IBAN**</span> |                                                                                                                                                                          |
| GET /?iban=AL86751639367318444714198669          | <pre>{<br>    iban: "AL86751639367318444714198669",<br>    valid: true,<br>    message: "OK"<br>}</pre>                                                                  |
| <span style="color: red">**Invalid IBAN**</span> |                                                                                                                                                                          |
| GET /?iban=GB8                                   | <pre>{<br>    iban: "GB8",<br>    valid: false,<br>    message: "IBAN is too short, expected > 4"<br>}</pre>                                                             |
| GET /?iban=GB83WEST12345698765432                | <pre>{<br>    iban: "GB83WEST12345698765432",<br>    valid: false,<br>    message: "checksum is incorrect"<br>}</pre>                                                    |
| GET /?iban=XX82WEST12345698765432                | <pre>{<br>    iban: "XX82WEST12345698765432",<br>    valid: false,<br>    message: "country code invalid"<br>}</pre>                                                     |
| GET /?iban=GB82WEST12345698765432                | <pre>{<br>    iban: "GB82WEST12345698765432*",<br>    valid: false,<br>    message: "Invalid characters, allowed are alphanumeric (A-Z, 0-9) and space (' ')"<br>}</pre> |
| GET /?iban=AL86751639367318444714198669AL86751   | <pre>{<br>    iban: "AL86751639367318444714198669AL86751",<br>    valid: false,<br>    message: "IBAN is too long, expected < 34"<br>}</pre>                             |

## Todo
- [x] Function to validate an IBAN
- [x] Basic webserver
- [ ] Webserver with timeouts
- [ ] Makefile for ~~build~~, ~~test~~ and load test
- [ ] Additional IBAN numbers for testing
- [x] Additional examples
- [ ] Clean up code
- [ ] Clean up README
- [ ] Move main package to cmd
- [ ] Flag to set port