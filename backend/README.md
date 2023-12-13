# IDS Backend

## Configuration

### cache

The IDS utilizes caching to reduce the number of calls to the SQL Database. There are two options, which will be explained under the `driver` setting.

#### driver

There are two drivers available for caching:

- memory -- This setting will utilize system memory within the IDS backend application to store the cache. This is useful if you do not have Redis available or do not want to use it. This setting is not recommended for instances where the backend may, will be, or is horizontally scaled as this will lead to data inconsistencies.
- redis -- This setting will utilize [Redis](https://redis.io/) to store the cache. This is recommended for environments where the backend may, will be, or is horizontally scaled as it will allow for data consistency as all instances will be pulling from the same cache. **NOTE** sentinel is not yet supported.

#### host (only if driver=redis)

Host of the Redis instance.

#### port (only if driver=redis)

Port of the Redis instance.

#### password (only if driver=redis)

Password of the Redis instance.

#### db (only if driver=redis)

Database number of the Redis instance.

#### default_expiration

There are multiple domains of data within the IDS. This setting will allow a per-domain specification of
how long we should cache that data for.

- airports -- Airport data, including METARs, TAFs, ATIS information, runway assignments, etc.
- charts -- Chart data

An update of the information, IE, ATIS information, will cause the airport cache to be invalidated and repulled on the next request.

### database

This is the database configuration used by [Gorm](https://gorm.io) to connect to the RDMS. Gorm supports
and we have implemented MySQL/MariaDB and PostgreSQL drivers. The configuration is the same for both.

#### driver

This is the driver to use. The supported drivers are:

- mysql
- postgres

#### host

This is the host of the database server.

#### port

This is the port of the database server.

#### username

This is the username to use to connect to the database server.

#### password

This is the password to use to connect to the database server.

#### database_name

This is the name of the database to connect to.

#### auto_migrate

This is a boolean value that will cause the IDS to automatically migrate the database schema on startup.

### facility

This is the facility configuration. Primarily used for display purposes.

#### identifier

This is the identifier of the facility used for display purposes. Generally this should be the FAA ID
of the overlying enroute airspace.

#### name

This is the name of the facility used for display purposes.

#### adh

For facilities that are part of the ADH Partnership, or utilize its webstack, we have the ability to restrict
access to the IDS to rostered controllers only.

##### api_base

The base of the API. For example, https://api.zanartcc.org. This is used to validate rostered controllers, if set.

##### rostered

A boolean value to determine whether or not to restrict access to rostered controllers. A call will be made to the api_base to look at the `controller_type` field for the controller. Where a 404 is returned, or the `controller_type` field is `none`, and this field is true the IDS will restrict access for this user.

### oauth

The OAuth configuration for the IDS.

#### provider

The IDS supports two OAuth2 providers. These are:

- adh -- The ADH stack includes an OAuth2 provider that is shared among the services of the webstack. This builds upon the VATSIM OAuth2 provider but allows restrictions to be placed on who can authenticate.
- vatsim -- For instances that do not utilize the ADH stack, we support the VATSIM OAuth2 provider. Please note that we do not provide any authorization restrictions for this provider, so any VATSIM member will be able to login.

#### base_url

Base URL of the OAuth2 provider. This is used to build the OAuth2 URL.

#### client_id

Client ID of the OAuth2 provider. This will be retrieved from the OAuth2 provider.

#### client_secret

Client Secret of the OAuth2 provider. This will be retrieved from the OAuth2 provider.

#### my_base_url

Base URL of the IDS backend. This is used to build the OAuth2 redirection URL.

#### endpoints

##### authorize

This is the authorize endpoint of the OAuth2 provider. This will be retrieved from the OAuth2 provider.

Generally:

- ADH: /oauth/authorize
- VATSIM: /oauth/authorize

##### token

This is the token endpoint of the OAuth2 provider. This will be retrieved from the OAuth2 provider.

Generally:

- ADH: /oauth/token
- VATSIM: /oauth/token

##### userinfo

This is the userinfo endpoint of the OAuth2 provider. This will be retrieved from the OAuth2 provider.

Generally:

- ADH: /v1/info
- VATSIM: /api/user

### server

#### port

This is the port the IDS backend will listen on.

#### mode

This is the mode the IDS backend will run in. The supported modes are:

- plain -- This will run the IDS backend in plain HTTP mode. This will be limited to HTTP/1.1 only.
- tls -- This will run the IDS backend in TLS mode. Supports HTTP/1.1 and HTTP/2.
- h2c -- This will run in plain text mode, but support both HTTP/1.1 and HTTP/2. Useful for running behind load balancers or ingress controllers that handle TLS termination without losing HTTP/2 functionality.

**NOTE** When running in TLS mode, the SSL_CERT and SSL_KEY environment variables must be set as described later in this document.

### session

The ADH IDS utilizes JSON Web Tokens for authentication. We do not support Refresh Tokens as the current
intention is to require authentication on every visit... provided it is outside the configured AccessExpire window.

#### algorithm

There are a number of algorithms supported by the underlying dependency. At the moment, we have
selected to only support symmetrical algorithms. We may consider expanding to using JWKs in the future.

The supported algorithms are ([Reference](https://pkg.go.dev/github.com/lestrrat-go/jwx/v2/jwa#SignatureAlgorithm)):

- HS256
- HS384
- HS512

### secret

This is a key used to sign the JWT. It is recommended to be a randomly generated string of at least 32 characters.
The string should be a mix of characters, numbers, letters, symbols, and casing.

### issuer

This is a set string that identifies the issuer of the JWT. It has no meaning outside of validation
and can be any string. Generally, though, it is an identifier of the issuer ([Reference](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1)).

### audience

This is a string that is converted into an array of a single string. It is generally used to identify the
audience of the JWT. We use it as part of the validation process. ([Reference](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3))

### access_expire

This is an integer that represents the number of seconds an access_token is valid for. A recommended setting
is 3600, which equates to 1 hour. Access tokens will be refreshed by the frontend while the sessions is
active. Lower numbers are better, but shouldn't be less than 15 minutes. ([Reference](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4))

## Environment Variables

### General

There are some environment variables that are required to be passed which generally contain 
sensitive information that should not be part of a ConfigMap. There are additional items in the
config, such as passwords, that generally should not be specified either. However, these have 
been left for ease of use given the nature of the project and the typical deployment environments
and teams.

Implementations can be made to have the configuration generated by init containers with passwords
populated at that point from secrets, vaults, etc. as desired.

### SSL_CERT

This should be the path to the SSL certificate to be used for the server *if* the server mode is set
to TLS. Many deployment environments put the IDS behind a Kubernetes Ingress Controller that can handle
TLS termination, where this may not be needed.

### SSL_KEY

This should be the path to the SSL key to be used for the server *if* the server mode is set to TLS.
