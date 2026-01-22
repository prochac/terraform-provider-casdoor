# Basic ACL model
resource "casdoor_model" "acl" {
  owner        = "my-organization"
  name         = "model-acl"
  display_name = "ACL Model"
  is_enabled   = true

  model_text = <<-EOT
    [request_definition]
    r = sub, obj, act

    [policy_definition]
    p = sub, obj, act

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
  EOT
}

# RBAC model with domains
resource "casdoor_model" "rbac_domain" {
  owner        = "my-organization"
  name         = "model-rbac-domain"
  display_name = "RBAC with Domains"
  is_enabled   = true

  model_text = <<-EOT
    [request_definition]
    r = sub, dom, obj, act

    [policy_definition]
    p = sub, dom, obj, act

    [role_definition]
    g = _, _, _

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
  EOT
}

# ABAC model
resource "casdoor_model" "abac" {
  owner        = "my-organization"
  name         = "model-abac"
  display_name = "ABAC Model"
  is_enabled   = true

  model_text = <<-EOT
    [request_definition]
    r = sub, obj, act

    [policy_definition]
    p = sub, obj, act

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = r.sub.Age > 18 && r.obj.Type == "document" && r.act == "read"
  EOT
}
