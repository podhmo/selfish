# selfish

selfish is individual gist client.

## install

```
go install github.com/podhmo/selfish/cmd/selfish@latest
```

## init setting

```bash
mkdir -p ~/.config/selfish
cat <<-EOS > ~/.config/selfish/config.json
{
  "access_token": "<your github access token>"
}
EOS
```

if you don't have github-access-token. [here](https://github.com/settings/tokens)

## uploading gists


```
$ cat <<-EOS > /tmp/hello.md
# hello
EOS

# create gist
$ ./bin/selfish --alias=hello hello.md
create success. (id="5639abca377b5c92061248666d38e6aa")
opening.. "https://gist.github.com/5639abca377b5c92061248666d38e6aa"

$ cat <<-EOS > /tmp/hello.md
# hello
hello hello hello
EOS

# update gist
$ ./bin/selfish --alias=hello --silent hello.md
update success. (id="5639abca377b5c92061248666d38e6aa")

# delete gist
$ ./bin/selfish --alias=hello --delete
deleted. (id="5639abca377b5c92061248666d38e6aa")
```

### help

```
$ ./bin/selfish -h
Usage of selfish:
      --alias string   ENV: ALIAS	alias name of uploaded gists
      --debug          ENV: DEBUG	-
      --delete         ENV: DELETE	delete uploaded gists
      --silent         ENV: SILENT	don't open gist pages with browser, after uploading
pflag: help requested
```

