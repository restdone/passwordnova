A tool to massively request breached password from ProxyNova https://www.proxynova.com/tools/comb/

```
go build passwordnova.go
./passwordnove -u <username list file> -o <output> -d <filter>
```
-d is optional, can be used to filter targeted domain
Output file can be used by 
```
hydra -C <file> <dc> <service>
```
