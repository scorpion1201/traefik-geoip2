# MaxMind GeoIP2 Plugin for Traefik #

This plugin is Traefik middleware that looks up IP information(Country, City, ASN)
from the MaxMind GeoIP DB.

This plugin is built using
[the Go GeoIP2 reader](https://github.com/oschwald/geoip2-golang).

## Usage ##

For a plugin to be active for a given Traefik instance,
it must be declared in the static configuration.

### Configuration ###

For each plugin, the Traefik static configuration must define
the module name (as is usual for Go packages)

The following declaration (given here in YAML) defines a plugin:

```yaml
# Static configuration

experimental:
  plugins:
    geoip2:
      moduleName: github.com/scorpion1201/traefik-geoip2
      version: 0.1.0
```

Here is an example of a file provider dynamic configuration (given here in YAML),
where the interesting part is the `http.middlewares` section:

```yaml
# Dynamic configuration
http:
  routers:
    my-router:
      rule: host(`foobar.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - geoip2-plugin
  services:
    service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  middlewares:
    geoip2-plugin:
      plugin:
        geoip2:
          ASNFileName: /path/to/GeoIP2-ASN.mmdb
          CityFileName: /path/to/GeoIP2-City.mmdb
          CountryFileName: /path/to/GeoIP2-Country.mmdb
```

## Contributing ##

Contributions welcome! Please fork the repository and open a pull request
with your changes.

## License ##

This is free software, licensed under the MIT license.