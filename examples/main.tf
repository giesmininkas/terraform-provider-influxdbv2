terraform {
  required_providers {
    influxdbv2 = {
      source = "matrasas.dev/dev/influxdbv2"
      version = "0.1"
    }
  }
}

provider "influxdbv2" {
  host = "http://localhost:8086"
  token = "J0HGVO8RGAq-gNfcppdPBXkqxNkTiSR9k4Ph3ilZYuC4mhqneFSBiTcjl3VoQx8gHJp81FdGSZcvx-9A_QTTkg=="
}