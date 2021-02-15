package common

var (
  // Version of mubeng itself
  Version = ""
  // Email handles of developer
  Email = "infosec@kitabisa.com"
  // Banner of mubeng
  Banner = `
           _   ` + Version + `
 _____ _ _| |_ ___ ___ ___ 
|     | | | . | -_|   | . |
|_|_|_|___|___|___|_|_|_  |
                      |___|
 ` + Email
  // Usage of mubeng
  Usage = `
  mubeng [-c|-a :8080] -f file.txt [options...]

Options:
  -f, --file <FILE>                Proxy file.
  -a, --address <ADDR>:<PORT>      Run proxy server.
  -c, --check                      To perform proxy live check.
  -t, --timeout                    Max. time allowed for proxy server/live check (default: 30s).
  -r, --rotate <AFTER>             Rotate proxy IP for every AFTER request.
  -v, --verbose                    Dump HTTP request/responses or show died proxy checks.
  -o, --output <FILE>              Log output from proxy servers or live proxy checks.

Examples:
  mubeng -f proxies.txt --check --output live.txt
  mubeng -a localhost:8080 -f live.txt -r 10
  mubeng -f live.txt --output file.log -d

`
)
