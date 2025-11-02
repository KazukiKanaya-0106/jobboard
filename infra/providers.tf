provider "aws" {
  region  = var.region
  profile = var.aws_profile

  default_tags {
    tags = {
      Application = var.app_name
      Stage       = var.stage
    }
  }
}
