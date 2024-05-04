A tool to  request breached password from ProxyNova https://www.proxynova.com/tools/comb/ and modify them for password spray/brute-force

passwordnova.go
=====================
```
go build passwordnova.go
./passwordnove -u <username list file> -t
```
It outputs passwordnova_result.txt contains the result from proxynova in format of [user email addaress]:[password]

![image](https://github.com/restdone/passwordnova/assets/42227817/baed1f09-4f30-433b-be9d-81ed4c51af7a)


-t is used to remove the domain in email address. It generates password_trim.txt in format of [username]:[password]

![image](https://github.com/restdone/passwordnova/assets/42227817/095d864c-7396-406d-a555-6d8ee76c74c4)


-nonum skip the line if the password is numberonly

It shows the username with more then 4 password. 

Output files can be used by hydra to perform combination test
```
hydra -C <file> <dc> <service>
```

userpassworcount.go
=====================
```
go build userpasswordcount
userpasswordcount -t x
```
count how many occurance for the same user name in password_trim. -t to specify the threasold. Before brute-forcing, it should make sure that such Brute-Force would not lock the account.

![image](https://github.com/restdone/passwordnova/assets/42227817/d0387a60-a211-4d4f-ba9b-df5b4a906e42)



removepassword.go
=================
After finding which account has too many password, this can be used to find the list of password belongs to this user and remove them from password_trim
```
go build removepassword.go
./removepassword <user name>
```
![image](https://github.com/restdone/passwordnova/assets/42227817/ea8129e6-c6e1-41f3-a99c-80057459d760)

Before:
![image](https://github.com/restdone/passwordnova/assets/42227817/bb19af8b-1047-4eb0-8c56-5667d5e2b75d)

After:
![image](https://github.com/restdone/passwordnova/assets/42227817/5c5a0b18-c130-4c45-a32f-eebe453f3929)



