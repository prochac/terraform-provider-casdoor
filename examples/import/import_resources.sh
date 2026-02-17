#!/bin/bash

import_list=${IMPORT_LIST:-"resources_for_import.txt"}
tmp_file="tmp-built-in.tf"
builtin_file="built-in.tf"

SED_CMD="sed"
if [[ "$OSTYPE" == "darwin"* ]]; then
  # Require gnu-sed.
  if ! [ -x "$(command -v gsed)" ]; then
    echo "Error: 'gsed' is not istalled." >&2
    echo "If you are using Homebrew, install with 'brew install gnu-sed'." >&2
    echo "Or remove the read-only attributes manually." >&2
    exit 1
  fi
  SED_CMD="gsed"
fi

detect_terraform() {
  # Allow override via environment
  if [[ -n "$TF_BIN" ]]; then
    if command -v "$TF_BIN" >/dev/null 2>&1; then
      echo "$TF_BIN"
      return
    else
      echo "Error: TF_BIN set to '$TF_BIN' but not found in PATH" >&2
      exit 1
    fi
  fi

  local has_terraform has_tofu

  has_terraform=$(command -v terraform >/dev/null 2>&1 && echo 1 || echo 0)
  has_tofu=$(command -v tofu >/dev/null 2>&1 && echo 1 || echo 0)

  if [[ $has_terraform -eq 1 && $has_tofu -eq 1 ]]; then
    echo "Both 'terraform' and 'tofu' are installed." >&2
    echo "Which one would you like to use? (terraform/tofu): " >&2
    read -r choice
    case "$choice" in
    terraform | tofu)
      echo "$choice"
      ;;
    *)
      echo "Invalid choice. Defaulting to terraform." >&2
      echo "terraform"
      ;;
    esac
  elif [[ $has_terraform -eq 1 ]]; then
    echo "terraform"
  elif [[ $has_tofu -eq 1 ]]; then
    echo "tofu"
  else
    echo "Error: Neither 'terraform' nor 'tofu' found in PATH" >&2
    exit 1
  fi
}

TF_BIN=$(detect_terraform)
echo "Using: $TF_BIN" >&2

# truncate the file
>${builtin_file}

while read -r res_type res_name fields; do
  [[ -z "${res_type}" || "${res_type}" == \#* ]] && continue

  # Parse key-value pairs from remaining fields
  skeleton=""
  import_id=""
  read -ra tokens <<<"${fields}"
  for ((i = 0; i < ${#tokens[@]}; i += 2)); do
    key="${tokens[i]}"
    val="${tokens[i + 1]}"
    skeleton+="  ${key} = \"${val}\""$'\n'

    # Build import ID: owner/name for most resources, bare id for LDAP
    case "${key}" in
    owner) import_id="${val}/${import_id}" ;;
    name) import_id="${import_id}${val}" ;;
    id) import_id="${val}" ;;
    esac
  done

  # Create temporary skeleton
  cat >${tmp_file} <<EOF
resource "${res_type}" "${res_name}" {
${skeleton}}
EOF

  ${TF_BIN} import "${res_type}.${res_name}" "${import_id}"

  ${TF_BIN} state show -show-sensitive "${res_type}.${res_name}" >>${builtin_file}
  echo >>${builtin_file}

done <${import_list}

# cleanup
[ -f ${tmp_file} ] && rm ${tmp_file}

# format
$TF_BIN fmt "${builtin_file}"

# remove read-only (Computed-only) fields
ro_attrs=(
  # common
  "created_time"
  "updated_time"
  # application
  "client_id"
  "client_secret"
  # user
  "is_default_avatar"
  "is_online"
  "hash"
  "pre_hash"
  "created_ip"
  "last_signin_time"
  "last_signin_ip"
  "last_change_password_time"
  "last_signin_wrong_time"
  "signin_wrong_times"
  # ldap
  "last_sync"
)
for ro_attr in "${ro_attrs[@]}"; do
  ${SED_CMD} -i'' -e "/^[[:space:]]*${ro_attr}[[:space:]]*=/d" ${builtin_file}
done

# Remove computed composite id (owner/name) but keep LDAP's bare id
${SED_CMD} -i'' -e '/^[[:space:]]*id[[:space:]]*=.*\//d' ${builtin_file}

# Strip trailing whitespace (fixes indented blank lines inside heredocs)
${SED_CMD} -i'' -e 's/[[:space:]]*$//' ${builtin_file}

# format after cleanup
$TF_BIN fmt "${builtin_file}"
