# selfish

selfish is an individual gist client for managing GitHub Gists from the command line.

## Features

- Create, update, and delete GitHub Gists
- Use aliases to easily manage your gists
- Automatically open created gists in your browser (or use `--silent` to skip)
- Support for both text and binary files
- Detect file types automatically and handle appropriately

## Installation

```bash
go install github.com/podhmo/selfish@latest
```

## Initial Setup

Before using selfish, you need to configure your GitHub access token:

```bash
mkdir -p ~/.config/selfish
cat <<-EOS > ~/.config/selfish/config.json
{
  "access_token": "<your github access token>"
}
EOS
```

If you don't have a GitHub access token, you can create one at [GitHub Settings > Personal access tokens](https://github.com/settings/tokens). Make sure to grant the `gist` scope.

## Usage

### Creating a Gist

Create a new gist with an alias:

```bash
$ cat <<-EOS > /tmp/hello.md
# hello
EOS

$ selfish --alias=hello /tmp/hello.md
create success. (id="5639abca377b5c92061248666d38e6aa")
opening.. "https://gist.github.com/5639abca377b5c92061248666d38e6aa"
```

### Updating a Gist

Update an existing gist using the same alias:

```bash
$ cat <<-EOS > /tmp/hello.md
# hello
hello hello hello
EOS

$ selfish --alias=hello --silent /tmp/hello.md
update success. (id="5639abca377b5c92061248666d38e6aa")
```

### Deleting a Gist

Delete a gist by its alias:

```bash
$ selfish --alias=hello --delete
deleted. (id="5639abca377b5c92061248666d38e6aa")
```

### Multiple Files

You can upload multiple files to a single gist:

```bash
$ selfish --alias=myproject file1.txt file2.md file3.py
```

## Options

```
Usage of selfish:
      --alias string   ENV: ALIAS	alias name of uploaded gists
      --client         ENV: CLIENT	if =fake, doesn't request {github, fake} (default github)
      --debug          ENV: DEBUG	enable debug output
      --delete         ENV: DELETE	delete uploaded gists
      --silent         ENV: SILENT	don't open gist pages with browser, after uploading
```

## Binary Files

selfish can handle binary files (e.g., images) by automatically cloning the gist repository, adding the binary files via git, and pushing them. Files larger than 5MB or detected as binary content types are handled this way.

## License

MIT

