A tool to massively request breached password from ProxyNova https://www.proxynova.com/tools/comb/

passwordnova.go
=====================
```
go build passwordnova.go
./passwordnove -u <username list file> -t
```
It outputs passwordnova_result.txt contains the result from proxynova in format of [user email addaress]:[password]

-t is used to remove the domain in email address. It generates password_trim.txt in format of [username]:[password]

-nonum skip the line if the password is numberonly

It shows the username with more then 4 password. 

Output files can be used by hydra to perform combination test
```
hydra -C <file> <dc> <service>
```


userpassworcount.go
=====================
count how many occurance for the same user name in password_trim. -t to specify the threasold. Before brute-forcing, it should make sure that such Brute-Force would not lock the account.
