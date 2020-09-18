## Big file uploader

Web app for upload big files by secret http link. It is simply web without https protocol.

### Prerequisits:

- os: linux
- sdk: golang
- svc: git
- free space for UPLOAD_DIRECTORY storage

### Installation:

You may:

- make TMP directory in /home and set environment in .env file for: 
```bash
TMPDIR=/home/tmp

```
- install golang environment in linux
- open incoming port 3000/tcp on the router
- open incoming connection with port 3000/tcp in your system firewall
- get this codebase from $GOPATH/[your_own_dir]:

```bash
git clone git@github.com:tvitcom/go-bigfileuploader.git bigfileuploader && cd bigfileuploader

```

- copy _env to .env file and make necessary settings (SECRET_LINK, UPLOAD_DIRECTORY, TMPDIR)
- build that:

```bash 
	go build -o bigfileuploader
```

- run as current user:

```bash
./run.sh
```

- send secret link to your uploading respondent as: http://[YOUP.IP]:3000/[SECRET_LINK]


