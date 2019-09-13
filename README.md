# neotama
Fast and Safe spider for Apache Directory Listings

## What?
A while back I needed a way to spider an apache directory listing and was surprised to find a lack of tools designed for the task. A regular directory spider would work, but those tend to spider every single URL they discover, instead of parsing the directory listing for file information. The closest command I was able to find was `lftp`, but it's also fairly "dumb", and isn't concurrent. And so, neotama was born.

Neotama is a spider specifically designed to parse directory listings, only querying directories and reading metadata (file size, modification time, etc) from the contents of the directory listing. The results are displayed in a `tree`-esque ASCII tree including color formatting based on file type. It's concurrent too, so it does this pretty damn fast. Finally, it's extensible, allowing spidering unsupported servers using a user-defined config file. The parsing is done using regex, so it's not too difficult to throw together a config file for an uncommon server.

## Install
You'll need to have a Go toolchain installed. On most ditributions the package is "golang".

Then, simply run `go get github.com/rsheasby/neotama`.

This will automatically get the latest version of neotama including dependencies, and compile it into your GOPATH.
You'll probably want to add the GOPATH binaries folder to your path too. By default this is simply `~/go/bin/`.

## Usage

```
usage: neotama [-h|--help] -u|--url "<value>" [-t|--threads <integer>]
               [-r|--retry <integer>] [-d|--depth <integer>]
               [--disable-sorting] [--color (auto|on|off|lol)] [-s|--server
               (auto|apache)] [-p|--parser-config "<value>"] [-o|--output
               (tree|list|urlencoded)]

               Safely and quickly crawls a directory listing, outputting a
               pretty tree.

Arguments:

  -h  --help             Print help information
  -u  --url              URL to crawl
  -t  --threads          Maximum number of concurrent connections. Default: 10
  -r  --retry            Maximum amount of times to retry a failed query.
                         Default: 3
  -d  --depth            Maximum depth to traverse. Depth of 0 means only query
                         the provided URL. Value of -1 means unlimited.
                         Default: -1
      --disable-sorting  Disables sorting. Default behavior is to sort by path
                         alphabetically, with files above directories
      --color            Whether to output color codes or not. Color codes will
                         be read from LS_COLORS if it exists, and will fallback
                         to some basic defaults otherwise. Default: auto
  -s  --server           Server type to use for parsing. Auto will detect the
                         server based on the HTTP headers. Default: auto
  -p  --parser-config    Config file to use for parsing the directory listing
  -o  --output           Output format of results. Default: tree
```

