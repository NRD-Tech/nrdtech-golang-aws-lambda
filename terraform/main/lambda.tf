# For Lambda deployed in a VPC.
#
# The function takes no inbound traffic, so this group has no ingress rules. Egress is
# scoped rather than allow-all: TCP 443 covers the AWS endpoints a function actually needs
# (Secrets Manager, SSM, S3, CloudWatch Logs) plus HTTPS APIs, and the VPC CIDR covers
# in-VPC dependencies like RDS. Widen it deliberately - add a rule for the specific port
# and destination - rather than reopening it to every port on every host.
#
# resource "aws_security_group" "lambda_sg" {
#   name   = "${var.app_ident}_lambda_sg"
#   vpc_id = data.aws_vpc.selected.id

#   egress {
#     description = "HTTPS to AWS service endpoints and external APIs"
#     from_port   = 443
#     to_port     = 443
#     protocol    = "tcp"
#     cidr_blocks = ["0.0.0.0/0"]
#   }

#   egress {
#     description = "All protocols to in-VPC dependencies"
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = [data.aws_vpc.selected.cidr_block]
#   }
# }

resource "aws_lambda_function" "lambda_function" {
  depends_on = [null_resource.push_image]
  function_name = var.APP_IDENT
  role          = aws_iam_role.lambda_exec.arn
  package_type  = "Image"
  image_uri = "${aws_ecr_repository.ecr_repository.repository_url}:${null_resource.push_image.triggers.code_hash}"
  timeout = var.APP_TIMEOUT
  memory_size = var.APP_MEMORY
  architectures = [var.CPU_ARCHITECTURE == "ARM64" ? "arm64" : "x86_64"]

  # For Lambda deployed in a VPC
  # Note: Use the vpc setting here if you need access to private resources like RDS
  #       instances but note that accessing general internet resources will not work
  #       unless your VPC has a NAT setup.
  # vpc_config {
  #   subnet_ids         = data.aws_subnets.private.ids
  #   security_group_ids = [aws_security_group.lambda_sg.id]
  # }

  environment {
    variables = {
      ENVIRONMENT = var.ENVIRONMENT,
      APP_IDENT = var.APP_IDENT
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name = "${var.APP_IDENT}_lambda_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Principal = {
          Service = "lambda.amazonaws.com",
        },
        Effect = "Allow",
      },
    ],
  })
}

resource "aws_iam_role_policy_attachment" "lambda_exec_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "lambda_policy" {
  name        = "${var.APP_IDENT}_Policy"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:CreateNetworkInterface",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DeleteNetworkInterface"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cloudwatch_put_metric_policy_attachment" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}
