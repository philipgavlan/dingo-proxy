# dingo-proxy
DinGo (Dns over tls IN Go) Proxy is a proxy listener that will accept DNS requests and route them to a 3rd party DNS service, defined in a JSON config file.

# running the service

## defaults and flags
All flags have defaults, so running the `main` binary without any flags will give the following options:
* Port (`-p`): `:53` - The listening port for the service.
* Protocol (`-t`): `tcp` - Listen for TCP or UDP connections **UDP in development**
* Config file (`-c`): `cloudflare-secure.json` - Config file that defines the 3rd party DNS service

## config file
The config file, written in JSON, must define the connection details for the 3rd party DNS service in key value pairs. The options are of the `Protocol`, whether TCP or UDP (**UDP in development**), the `ConnectionString`, the address:port of the service's DNS listener, and `TLS`, a boolean defining whether or not to use TLS encryption.
