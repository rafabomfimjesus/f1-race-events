resource "aws_ecs_cluster" "main" {
  name = "f1-race-events-cluster"
}

resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/f1-race-events"
  retention_in_days = 7
}

resource "aws_ecs_task_definition" "app" {
  family                   = "f1-race-events"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  cpu    = 256
  memory = 512

  execution_role_arn = aws_iam_role.ecs_execution_role.arn
  task_role_arn      = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name = "f1-race-events"

      image = "${aws_ecr_repository.app.repository_url}:latest"

      essential = true

      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"

        options = {
          awslogs-group         = aws_cloudwatch_log_group.ecs.name
          awslogs-region        = "us-east-2"
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}