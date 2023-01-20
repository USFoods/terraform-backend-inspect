terraform {
  required_version = "~> 1.3"
  backend "s3" {
    bucket = "aws.basic.bucket"
  }
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~>3.11.0"
    }
  }
}
