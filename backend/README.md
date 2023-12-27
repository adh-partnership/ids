# IDS Backend

## Table of Contents

- [IDS Backend](#ids-backend)
  - [Table of Contents](#table-of-contents)
  - [Configuration](#configuration)
    - [cache](#cache)
      - [driver](#driver)
      - [host (only if driver=redis)](#host-only-if-driverredis)
      - [port (only if driver=redis)](#port-only-if-driverredis)
      - [password (only if driver=redis)](#password-only-if-driverredis)
      - [db (only if driver=redis)](#db-only-if-driverredis)
      - [default\_expiration](#default_expiration)
    - [database](#database)
      - [driver](#driver-1)
      - [host](#host)
      - [port](#port)
      - [username](#username)
      - [password](#password)
      - [database\_name](#database_name)
      - [auto\_migrate](#auto_migrate)
    - [facility](#facility)
      - [identifier](#identifier)
      - [name](#name)
      - [adh](#adh)
        - [api\_base](#api_base)
        - [rostered](#rostered)
    - [oauth](#oauth)
      - [provider](#provider)
      - [base\_url](#base_url)
      - [client\_id](#client_id)
      - [client\_secret](#client_secret)
      - [my\_base\_url](#my_base_url)
      - [endpoints](#endpoints)
        - [authorize](#authorize)
        - [token](#token)
        - [userinfo](#userinfo)
    - [server](#server)
      - [ip](#ip)
      - [port](#port-1)
      - [mode](#mode)
    - [session](#session)
      - [block\_key](#block_key)
    - [hash\_key](#hash_key)
    - [name](#name-1)
    - [path](#path)
    - [domain](#domain)
    - [max\_age](#max_age)
    - [secure](#secure)
    - [http\_only](#http_only)
  - [Environment Variables](#environment-variables)
    - [General](#general)
    - [SSL\_CERT](#ssl_cert)
    - [SSL\_KEY](#ssl_key)

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

---

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

---

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

---

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

---

### server

#### ip

The ip to bind to. This should be left blank to bind to all interfaces. This is not needed, except if developing
on Windows Subsystem for Linux, where dualstack is not supported correctly at the present moment. Binding to all
interfaces causes connections to fail outside of IPv6 from outside the WSL environment. In this instance, bind to
IP `127.0.0.1` unless you want to use IPv6-only.

#### port

This is the port the IDS backend will listen on.

#### mode

This is the mode the IDS backend will run in. The supported modes are:

- plain -- This will run the IDS backend in plain HTTP mode. This will be limited to HTTP/1.1 only.
- tls -- This will run the IDS backend in TLS mode. Supports HTTP/1.1 and HTTP/2.
- h2c -- This will run in plain text mode, but support both HTTP/1.1 and HTTP/2. Useful for running behind load balancers or ingress controllers that handle TLS termination without losing HTTP/2 functionality.

**NOTE** When running in TLS mode, the SSL_CERT and SSL_KEY environment variables must be set as described later in this document.

---

### session

The ADH IDS utilizes sessions via cookies for authentication. The only data we really care about is, are they
a VATSIM member, and if configured in the Facility section, are they rostered (ADH members only due to reliance
on the API).

#### block_key

The block key used for encryption. The length should be 16 bytes (AES-128), 20 bytes (AES-192), or 32 bytes (AES-256). The length dictates the encryption algorithm used. 32 bytes is highly encouraged.

---

### hash_key

The hash key used for encryption. This should be at least 32 bytes but may be longer.

---

### name

The name of the cookie to set. This should be something unique to the IDS, such as "ids" or "ids-session".

---

### path

Sets the path option of the cookie. This should be set to "/" to allow the cookie to be sent to all paths used by the IDS API and frontend.

---

### domain

This should be the domain of the IDS. Note that if the IDS backend and frontend run on different subdomains, 
you should set this to the common domain name used between them. For example if the IDS API is at
ids-api.zanartcc.org and the frontend is at ids.zanartcc.org, set this to zanartcc.org. While calls to the backend are the only ones that necessitate the cookie, cross-subdomain cookies are difficult and may lose support.

---

### max_age

This is how long the cookie should last. By default, we set this to 86400 seconds (24 hours). The cookie is refreshed
on every request, so this time will be reset. If you set this to 0, the cookie will be a session cookie and will be
removed when the browser is closed.

---

### secure

This is a boolean value that determines whether or not the cookie should be set as secure. This should be set to true.

---

### http_only

This is a boolean value that determines whether or not the cookie should be set as HTTP only. This should be set to true.

## Environment Variables

---

### General

There are some environment variables that are required to be passed which generally contain 
sensitive information that should not be part of a ConfigMap. There are additional items in the
config, such as passwords, that generally should not be specified either. However, these have 
been left for ease of use given the nature of the project and the typical deployment environments
and teams.

Implementations can be made to have the configuration generated by init containers with passwords
populated at that point from secrets, vaults, etc. as desired.

---

### SSL_CERT

This should be the path to the SSL certificate to be used for the server *if* the server mode is set
to TLS. Many deployment environments put the IDS behind a Kubernetes Ingress Controller that can handle
TLS termination, where this may not be needed.

---

### SSL_KEY

This should be the path to the SSL key to be used for the server *if* the server mode is set to TLS.
