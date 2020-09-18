## Big file uploader

Web app for upload big files by secret http link. It is simply web without https protocol.
You may :

- make TMP directory in /home and set environment in .env file for: 
```bash
TMPDIR=/home/tmp

```
- install golang environment in linux
- open incoming port 3000/tcp on the router
- open incoming connection with port 3000/tcp in your system firewall 
- change uploadsDir variable in main.go
- build that:

```bash 
	go build -o bigfileuploader
```

### prerequisits:

- os: linux
- sdk: golang
- free space for uploadsDir strage

