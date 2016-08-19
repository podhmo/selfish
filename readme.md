# selfish

selfish is individual gist client.

## init setting

```bash
mkdir -p ~/.selfish
cat <<-EOS > ~/.selfish/config.json
{
  "access_token": "<your github access token>"
}
EOS
```

if you doesn't have github-access-token. [here](https://github.com/settings/tokens)

## uploading gists


```
$ cat <<-EOS > /tmp/hello.md
# hello
EOS

# create gist
$ ./bin/selfish -alias=hello hello.md
create success. (id="5639abca377b5c92061248666d38e6aa")
opening.. "https://gist.github.com/5639abca377b5c92061248666d38e6aa"

$ cat <<-EOS > /tmp/hello.md
# hello
hello hello hello
EOS

# update gist
$ ./bin/selfish -alias=hello -silent hello.md
update success. (id="5639abca377b5c92061248666d38e6aa")

# delete gist
$ ./bin/selfish -alias=hello -delete
deleted. (id="5639abca377b5c92061248666d38e6aa")
```

### help

```
$ ./bin/selfish -h
Usage of ./bin/selfish:
  -alias string
        alias name of uploaded gists
  -delete
        delete uploaded gists
  -silent
        deactivate webbrowser open, after gists uploading
```

