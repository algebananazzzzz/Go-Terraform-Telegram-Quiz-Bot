output "api_gateway_endpoint" {
  value = "${module.api_lambda_integration.api_endpoint}/latest/"
}
