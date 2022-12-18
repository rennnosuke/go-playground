terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region = "ap-northeast-1"
}

resource "aws_dynamodb_table" "test-dax-game-score" {
  name           = "test-dax-game-score"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_id"
  range_key      = "game_id"

  attribute {
    name = "user_id"
    type = "N"
  }

  attribute {
    name = "game_id"
    type = "N"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  tags = {
    Name        = "test-dax-game-score"
    Environment = "test"
  }
}

