terraform {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
  profile = "iam_admin"
}

resource "aws_iam_policy" "ec2_ecr_rds_full_access" {
  name        = "EC2ECRRDSFullAccess"
  description = "Provides full access to EC2, ECR, and RDS"
  policy      = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid    = "FullEC2Access",
        Effect = "Allow",
        Action = "ec2:*",
        Resource = "*"
      },
      {
        Sid    = "FullECRAccess",
        Effect = "Allow",
        Action = [
          "ecr:*",
          "ecr-public:*"
        ],
        Resource = "*"
      },
      {
        Sid    = "FullRDSAccess",
        Effect = "Allow",
        Action = "rds:*",
        Resource = "*"
      }
    ]
  })
}