resource "aws_lb" "app" {
  name               = "f1-race-events-alb"
  load_balancer_type = "application"

  security_groups = [
    aws_security_group.alb.id
  ]

  subnets = data.aws_subnets.default.ids
}

resource "aws_lb_target_group" "app" {
  name        = "f1-race-events-tg"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"

  vpc_id = data.aws_vpc.default.id

  health_check {
    path = "/race-events"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.app.arn

  port     = 80
  protocol = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app.arn
  }
}