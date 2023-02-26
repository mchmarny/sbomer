# Description: This file contains the file resources for the deployment

# data.template_file.version.rendered
data "template_file" "version" {
  template = file("../.version")
}

# data.template_file.schema_vul.rendered
data "template_file" "schema_vul" {
  template = file("schema/vul.json")
}

# data.template_file.schema_pkg.rendered
data "template_file" "schema_pkg" {
  template = file("schema/pkg.json")
}
