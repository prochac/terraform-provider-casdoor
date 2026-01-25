#!/bin/bash

import_list="resources_for_import.txt"
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

while read -r res_type res_name id_field id_value; do
  [[ -z "${res_type}" || "${res_type}" == \#* ]] && continue

  # Create temporary skeleton
  cat >${tmp_file} <<EOF
resource "${res_type}" "${res_name}" {
  $id_field = "${id_value}"
}
EOF

  ${TF_BIN} import "${res_type}.${res_name}" "${id_value}"

  ${TF_BIN} state show -show-sensitive "${res_type}.${res_name}" >>${builtin_file}
  echo >>${builtin_file}

done <${import_list}

# cleanup
[ -f ${tmp_file} ] && rm ${tmp_file}

# format
$TF_BIN fmt "${builtin_file}"

# remove read-only fields
ro_attrs=(
  "created_time"
  "client_id"
  "client_secret"
)
for ro_attr in "${ro_attrs[@]}"; do
  ${SED_CMD} -i'' -e "/$ro_attr/d" ${builtin_file}
done

# format after cleanup
$TF_BIN fmt "${builtin_file}"