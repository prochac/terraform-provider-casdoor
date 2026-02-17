# Casdoor Server Quirks

Documented behavior of the Casdoor server that affects this provider.
Based on reading the server source at
[casdoor/casdoor](https://github.com/casdoor/casdoor).

## Resource ID Format

All API endpoints use `id=owner/name` as the query parameter to identify a
single resource. The server parses this with `strings.Split(id, "/")` and
queries the database by **both** `owner` AND `name` as a composite key.

There is no unique constraint on `name` alone -- two different organizations
can have a resource with the same name.

## Owner Fallback Behavior

Some resources have implicit fallback logic on the server side: if a resource
is not found under the requested owner, the server silently retries with a
hardcoded fallback owner. This means a `GetX` call may return a resource you
didn't ask for.

| Resource | Fallback owner | Server code |
|----------|----------------|-------------|
| **Cert** | `"admin"` | `object/cert.go` -- `GetCert()` retries with `owner="admin"` if not found |
| **Model** | `"built-in"` | `object/model.go` -- `getModelEx()` retries with `owner="built-in"` if not found |

No other resources have this fallback behavior as of 2025.

### Impact on users

This affects anyone using `casdoor_cert` or `casdoor_model` resources:

- **Phantom state:** If your cert/model is deleted outside Terraform (or never
  existed), the server silently returns the `"admin"` / `"built-in"` fallback
  resource instead of a 404. Terraform sees a successful response, thinks your
  resource still exists, and shows no drift. You may not notice it's gone.
- **Wrong resource in state:** After an import or a recreate, Terraform could
  latch onto the fallback resource. Subsequent `terraform apply` would then
  attempt to **update the shared fallback resource**, potentially breaking
  other organizations that depend on it.
- **Silent name collisions:** If you create a cert with the same name as an
  admin-owned cert, `Read` will find it. But if yours is later deleted, `Read`
  seamlessly falls back to the admin's cert with the same name -- no error,
  no warning.

### Mitigation in the provider

The provider's `Read` methods for cert and model should verify that the
`owner` field in the API response matches the `owner` in state. If it doesn't,
the resource should be treated as deleted (remove from state) rather than
accepting the fallback.

### Cert listing also merges owners

`GetCerts(owner)` uses `WHERE owner = 'admin' OR owner = ?`, so listing certs
always includes admin-owned certs mixed with the requested organization's
certs.

## Admin-Owned Resources

Some resources are always owned by `"admin"` regardless of which organization
they logically belong to:

| Resource | Owner |
|----------|-------|
| Organization | `"admin"` |
| Application | `"admin"` |
| Token | `"admin"` |
| LDAP | `"admin"` |

For these resources, the `owner` field in the API is always `"admin"`. The
organizational association is expressed through other fields (e.g.,
`application.Organization`).

## Session: Non-Standard ID Format

Sessions use a 3-part ID format `owner/user/application` via the
`sessionPkId` query parameter (not the standard `id` parameter). This is
a completely different pattern from all other resources.

## LDAP: ID Is Not a Name

LDAP resources use a server-generated UUID as their primary key, not a
human-readable name. The API still expects `id=owner/ldap_id` format, but
the second part is the UUID `Id` field, not `Name`. The server extracts only
the UUID part and queries by `Id` alone (ignoring the owner in the query).
