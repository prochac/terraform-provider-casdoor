# Upload a text file
resource "casdoor_resource" "readme" {
  user           = "admin"
  tag            = "txt"
  file_name      = "/docs/readme.txt"
  content_base64 = base64encode("Hello, World!")
}

# Upload an image from disk
resource "casdoor_resource" "logo" {
  user           = "admin"
  tag            = "img"
  file_name      = "/images/logo.png"
  content_base64 = filebase64("${path.module}/assets/logo.png")
  description    = "Organization logo"
}

# Upload a document with a parent path
resource "casdoor_resource" "policy" {
  user           = "admin"
  tag            = "pdf"
  parent         = "policies"
  file_name      = "/policies/privacy-policy.pdf"
  content_base64 = filebase64("${path.module}/assets/privacy-policy.pdf")
  description    = "Privacy policy document"
}
