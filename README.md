A tool to massively request breached password from ProxyNova https://www.proxynova.com/tools/comb/

```
go build passwordnova.go
./passwordnove -u <username list file> -t
```
It outputs passwordnova_result.txt and password_trim.txt, where password_trim remove doamin
-t is used to remove the domain in email address
Output files can be used by hydra to perform combination test
```
hydra -C <file> <dc> <service>
```
