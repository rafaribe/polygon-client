
output "alb_url" {
  value = "http://${aws_alb.polygon_client.dns_name}"
}

output "domain_validations" {
  value = aws_acm_certificate.polygon_client.domain_validation_options
}
