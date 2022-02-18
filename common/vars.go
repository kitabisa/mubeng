package common

var (
  // App name
  App = "mubeng"
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
  GENERAL
    -f, --file <FILE>                Proxy file
    -o, --output <FILE>              Write log output to FILE
    -t, --timeout <TIME>             Max. time allowed for connection (default: 30s)
    -u, --update                     Update mubeng to the latest stable version
    -v, --verbose                    Verbose mode
    -V, --version                    Show current mubeng version
  
  PROXY CHECKER
    -c, --check                      Perform proxy check
        --only-cc <AA>,<BB>          Only for specific country code (comma separated)
  
  IP ROTATOR
    -a, --address <ADDR>:<PORT>      Run proxy server
    -A, --auth <USER>:<PASS>         Set authorization for proxy server
    -d, --daemon                     Daemonize proxy server
    -m, --method <METHOD>            Rotation method (sequent/random) (default: sequent)
    -r, --rotate <N>                 Rotate proxy IP after N request (default: 1)
    -s, --sync                       Syncrounus mode

Examples:
  mubeng -f proxies.txt --check --output live.txt
  mubeng -a localhost:8080 -f live.txt -r 10

`
)
