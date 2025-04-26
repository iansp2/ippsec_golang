# Golang for Hackers (ippsec)
Repo to save projects as I learn on IppSec's Golang for Hacking youtube course. 

Projects roadmap:
- LDAP injector (in progress)
- Automating boolean sql injections (originally in Python, I will implement in go)

## LDAP injector
The HTB machine used for this is Ghost (https://app.hackthebox.com/machines/616)

### EP01: https://www.youtube.com/watch?v=uJFW4c4QE0U
- Start project creating injector struct
- Charset prunning by determining which chars are in password

### EP02: https://www.youtube.com/watch?v=BhLpqRev80s
- Dependency injection so that injector only handles injection and takes in object that handles http client. Interface is used for this.
- Create NetHTTP and FastHTTP implementations to illustrate how you can have different objects used as parameters on implementation without having to change inhector object

## Automating boolean sql injections
- Python implementation created
