terraform {
  backend "local" {
    path = "integration.tfstate"
  }
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~>3.11.0"
    }
  }
}