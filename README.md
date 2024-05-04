A tool to massively request breached password from ProxyNova https://www.proxynova.com/tools/comb/

```
go build passwordnova.go
./passwordnove -u <username list file> -t
```
It outputs passwordnova_result.txt contains the result from proxynova in format of [user email addaress]:[password]

-t is used to remove the domain in email address. It generates password_trim.txt in format of [username]:[password]

Output files can be used by hydra to perform combination test
```
hydra -C <file> <dc> <service>
```

