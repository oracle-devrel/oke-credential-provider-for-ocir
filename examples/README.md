# Provider Configuration

## TOKEN_VALIDATION

To utilize token validation, set to `enabled`.

When the first image is pulled, the CredentialProvider will obtain a token from OCI which will be cached for subsequent image pulls up to the `cacheDuration` (normally equal to the tokens lifetime: 1hr).  When the cached token expires by reaching the `cacheDuration`, the next image pull will request a new one.

If IAM policies are not in place at the time the token is cached then the token is essentially unauthorized for the entire duration of the `cacheDuration`.  The `TOKEN_VALIDATION` configuration setting will set an initial short lived `cacheDuration` until it is determined that the cached token is authorized to pull images.  Once it is determined that the token is authorized, it will set the `cacheDuration` to the tokens lifetime.

It is important to note that a new token is _not_ requested when the cached token expires.  It is only requested when the cached token expires **and** an image pull occurs.