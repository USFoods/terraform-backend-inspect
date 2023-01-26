terraform {
  required_version = "~> 1.3"
  backend "local" {
    path = "state.tfstate"
  }
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~>3.11.0"
    }
  }
}
