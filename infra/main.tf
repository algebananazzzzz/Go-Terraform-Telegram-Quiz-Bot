provider "aws" {
  region  = var.aws_region
  profile = var.profile
}

locals {
  function_name = "${var.env}-func-goquizbot"
  environment_variables = {
    TOKEN         = ""
    REDIS_ADDR    = "xyz.cache.amazonaws.com:6379"
    REDIS_KEY     = "foo"
    START_MESSAGE = "Welcome to GoQuizBot!"
    ERROR_MESSAGE = "An unexpected error occured."
    PASSPHRASE    = "" # Optional, omit to disable passphrase authentication
  }
  api_gateway_name    = "${var.env}-apigw-goquizbot"
  execution_role_name = "${var.env}-role-goquizbot"
  execution_role_policy_document = {
    name = "${var.env}-policy-goquizbot"
    statements = {
      allowCreateNetworkInterface = {
        effect = "Allow"
        actions = [
          "ec2:CreateNetworkInterface",
          "ec2:DescribeNetworkInterfaces",
          "ec2:DeleteNetworkInterface",
          "ec2:AssignPrivateIpAddresses",
          "ec2:UnassignPrivateIpAddresses"
        ]
        resources = ["*"]
      }
    }
  }
}

module "lambda_function" {
  source                         = "github.com/algebananazzzzz/terraform_modules/modules/lambda_function"
  function_name                  = local.function_name
  runtime                        = "provided.al2"
  handler                        = "bootstrap"
  environment_variables          = local.environment_variables
  execution_role_name            = local.execution_role_name
  execution_role_policy_document = local.execution_role_policy_document
  deployment_package = {
    local_path = "./build"
  }

  vpc_config = {
    subnet_ids         = [data.aws_subnet.private.ids]
    security_group_ids = data.aws_security_groups.allow_nat.ids
  }
}


module "api_lambda_integration" {
  source           = "github.com/algebananazzzzz/terraform_modules/modules/api_lambda_integration"
  api_gateway_name = local.api_gateway_name
  lambda_integrations = {
    lambda = {
      function_name   = module.lambda_function.function_name
      path            = "latest"
      integration_uri = module.lambda_function.function_invoke_arn
    }
  }
}
