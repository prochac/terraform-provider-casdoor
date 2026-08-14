# Changelog

## 0.1.0

BREAKING CHANGES:

This release targets the Casdoor 3.x schema. Casdoor 3.0.0 removed a number of
fields that earlier provider versions exposed, and configurations that set any
of the attributes below will now fail with `Unsupported argument`. Remove them
from your configuration before upgrading:

* `casdoor_adapter`: `table_name_prefix`, `is_enabled`
* `casdoor_cert`: `authority_public_key`, `authority_root_public_key`
* `casdoor_enforcer`: `is_enabled`
* `casdoor_model`: `manager`, `contact_email`, `type`, `parent_id`,
  `is_top_model`, `is_enabled`, `updated_time`
* `casdoor_pricing`: `submitter`, `approver`, `approve_time`, `state`
* `casdoor_user`: `access_key`, `access_secret`

`casdoor_model` carried `casdoor_group`'s fields by mistake; the server's model
object only has `owner`, `name`, `created_time`, `display_name`, `description`
and `model_text`.

Requires Casdoor >= 3.0.0. Provider 0.0.13 and earlier target the pre-3.x
schema and are not compatible with a 3.x server.

ENHANCEMENTS:

* Dropped the `casdoor-go-sdk` fork in favour of upstream v1.53.0, which merged
  the fork's owner-handling and struct-sync changes.

BUG FIXES:

* `casdoor_organization`: `Create` populated state from the request rather than
  the server response, so server-assigned defaults never reached state and
  produced a diff on the next refresh.
* `casdoor_organization`: `account_items`, `country_codes`, `languages` and
  `password_options` are now `Optional`/`Computed`, so the defaults Casdoor
  populates no longer show as perpetual drift.
* `casdoor_user`: `avatar` and `ranking` no longer force empty defaults over the
  values Casdoor assigns on create, and `Update` now reads server state back.
