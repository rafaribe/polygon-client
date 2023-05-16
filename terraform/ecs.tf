locals {
  repository_url     = "ghcr.io/rafaribe/polygon-client"
  region             = "us-east-1"
  name               = "polygon-client"
  log_retention_days = 120
  tags = {
    Name       = local.name
    Repository = "https://github.com/rafaribe/polygon-client"
  }
}

# We need a cluster in which to put our service.
resource "aws_ecs_cluster" "app" {
  name = "app"
  tags = local.tags
}

# An ECR repository is a private alternative to Docker Hub, app could be pushed there, optional.
resource "aws_ecr_repository" "polygon_client" {
  name = "polygon-client"
  tags = local.tags
}

# Log groups hold logs from our app.
resource "aws_cloudwatch_log_group" "polygon_client" {
  name              = "/ecs/polygon-client"
  retention_in_days = local.log_retention_days
  tags              = local.tags
}

# The main service.
resource "aws_ecs_service" "polygon_client" {
  name            = "polygon-client"
  task_definition = aws_ecs_task_definition.polygon_client.arn
  cluster         = aws_ecs_cluster.app.id
  launch_type     = "FARGATE"

  desired_count = 1

  load_balancer {
    target_group_arn = aws_lb_target_group.polygon_client.arn
    container_name   = "polygon-client"
    container_port   = "3000"
  }

  network_configuration {
    assign_public_ip = false

    security_groups = [
      aws_security_group.egress_all.id,
      aws_security_group.ingress_api.id,
    ]

    subnets = [
      aws_subnet.private_d.id,
      aws_subnet.private_e.id,
    ]
  }
}

# The task definition for our app.
resource "aws_ecs_task_definition" "polygon_client" {
  family = "polygon-client"

  container_definitions = <<EOF
  [
    {
      "name": "polygon-client",
      "image": "${local.repository_url == "" ? aws_ecr_repository.polygon_client.repository_url : local.repository_url}:latest",
      
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-region": "${local.region}}",
          "awslogs-group": "/ecs/polygon-client",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]

EOF

  execution_role_arn = aws_iam_role.polygon_client_task_execution_role.arn

  # These are the minimum values for Fargate containers.
  cpu                      = 256
  memory                   = 512
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
}

# This is the role under which ECS will execute our task. This role becomes more important
# as we add integrations with other AWS services later on.

# The assume_role_policy field works with the following aws_iam_policy_document to allow
# ECS tasks to assume this role we're creating.
resource "aws_iam_role" "polygon_client_task_execution_role" {
  name               = "polygon-client-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume_role.json
}

data "aws_iam_policy_document" "ecs_task_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

# Normally we'd prefer not to hardcode an ARN in our Terraform, but since this is an AWS-managed
# policy, it's okay.
data "aws_iam_policy" "ecs_task_execution_role" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Attach the above policy to the execution role.
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.polygon_client_task_execution_role.name
  policy_arn = data.aws_iam_policy.ecs_task_execution_role.arn
}

resource "aws_lb_target_group" "polygon_client" {
  name        = "polygon-client"
  port        = 3000
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.app_vpc.id

  health_check {
    enabled = true
    path    = "/health"
  }

  depends_on = [aws_alb.polygon_client]
}

// From here on out, it is something that will not be needed in this specific application, however I want to showcase how one application would be exposed to the internet.

resource "aws_alb" "polygon_client" {
  name               = "polygon-client-lb"
  internal           = false
  load_balancer_type = "application"

  subnets = [
    aws_subnet.public.id,
  ]

  security_groups = [
    aws_security_group.http.id,
    aws_security_group.https.id,
    aws_security_group.egress_all.id,
  ]

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_alb_listener" "polygon_client_http" {
  load_balancer_arn = aws_alb.polygon_client.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_alb_listener" "polygon_client_https" {
  load_balancer_arn = aws_alb.polygon_client.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate.polygon_client.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.polygon_client.arn
  }
}

resource "aws_acm_certificate" "polygon_client" {
  domain_name       = "polygon-client.jimmysawczuk.net"
  validation_method = "DNS"
}
