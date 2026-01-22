# AWS ECS Deployment Guide

## Prerequisites

- AWS Account with appropriate permissions
- Docker installed locally
- AWS CLI configured
- ECR repository created

## Step 1: Push Docker Image to ECR

### 1.1 Login to ECR
```bash
AWS_ACCOUNT_ID=<your-account-id>
AWS_REGION=us-east-1
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com
```

### 1.2 Create ECR Repository (if not exists)
```bash
aws ecr create-repository \
  --repository-name expense-tracker \
  --region $AWS_REGION
```

### 1.3 Build and Push Image
```bash
docker build -t expense-tracker:latest .
docker tag expense-tracker:latest $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/expense-tracker:latest
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/expense-tracker:latest
```

## Step 2: Create RDS PostgreSQL Instance

### 2.1 Create DB Subnet Group
```bash
aws rds create-db-subnet-group \
  --db-subnet-group-name expense-tracker-subnet \
  --db-subnet-group-description "Subnet group for Expense Tracker" \
  --subnet-ids subnet-xxxxx subnet-yyyyy \
  --region $AWS_REGION
```

### 2.2 Create RDS Instance
```bash
aws rds create-db-instance \
  --db-instance-identifier expense-tracker-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --engine-version 15.1 \
  --master-username postgres \
  --master-user-password "YourSecurePassword" \
  --allocated-storage 20 \
  --db-subnet-group-name expense-tracker-subnet \
  --publicly-accessible true \
  --region $AWS_REGION
```

Wait for instance to be available (check AWS console)

## Step 3: Create ECS Cluster

### 3.1 Create Cluster
```bash
aws ecs create-cluster \
  --cluster-name expense-tracker-cluster \
  --region $AWS_REGION
```

## Step 4: Create IAM Role for ECS Tasks

### 4.1 Create Trust Policy (trust-policy.json)
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

### 4.2 Create IAM Role
```bash
aws iam create-role \
  --role-name ecsTaskExecutionRole \
  --assume-role-policy-document file://trust-policy.json \
  --region $AWS_REGION

aws iam attach-role-policy \
  --role-name ecsTaskExecutionRole \
  --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
```

## Step 5: Create Task Definition

### 5.1 Create task-definition.json
```json
{
  "family": "expense-tracker",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "expense-tracker",
      "image": "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/expense-tracker:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "hostPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "PORT",
          "value": "8080"
        },
        {
          "name": "DB_HOST",
          "value": "expense-tracker-db.xxxxx.rds.amazonaws.com"
        },
        {
          "name": "DB_PORT",
          "value": "5432"
        },
        {
          "name": "DB_USER",
          "value": "postgres"
        },
        {
          "name": "DB_NAME",
          "value": "expense_tracker"
        }
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:$AWS_REGION:$AWS_ACCOUNT_ID:secret:db-password"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/expense-tracker",
          "awslogs-region": "$AWS_REGION",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ],
  "executionRoleArn": "arn:aws:iam::$AWS_ACCOUNT_ID:role/ecsTaskExecutionRole"
}
```

### 5.2 Register Task Definition
```bash
aws ecs register-task-definition \
  --cli-input-json file://task-definition.json \
  --region $AWS_REGION
```

## Step 6: Create ALB (Application Load Balancer)

### 6.1 Create Security Group
```bash
aws ec2 create-security-group \
  --group-name expense-tracker-alb-sg \
  --description "Security group for Expense Tracker ALB" \
  --vpc-id vpc-xxxxx

SG_ID=<security-group-id>

aws ec2 authorize-security-group-ingress \
  --group-id $SG_ID \
  --protocol tcp \
  --port 80 \
  --cidr 0.0.0.0/0
```

### 6.2 Create ALB
```bash
aws elbv2 create-load-balancer \
  --name expense-tracker-alb \
  --subnets subnet-xxxxx subnet-yyyyy \
  --security-groups $SG_ID \
  --region $AWS_REGION
```

### 6.3 Create Target Group
```bash
aws elbv2 create-target-group \
  --name expense-tracker-tg \
  --protocol HTTP \
  --port 8080 \
  --vpc-id vpc-xxxxx \
  --target-type ip \
  --health-check-path /health \
  --health-check-interval-seconds 30 \
  --health-check-timeout-seconds 5 \
  --healthy-threshold-count 2 \
  --unhealthy-threshold-count 2
```

### 6.4 Create Listener
```bash
ALB_ARN=<load-balancer-arn>
TG_ARN=<target-group-arn>

aws elbv2 create-listener \
  --load-balancer-arn $ALB_ARN \
  --protocol HTTP \
  --port 80 \
  --default-actions Type=forward,TargetGroupArn=$TG_ARN
```

## Step 7: Create ECS Service

### 7.1 Create Service
```bash
CLUSTER_NAME=expense-tracker-cluster
SERVICE_NAME=expense-tracker-service

aws ecs create-service \
  --cluster $CLUSTER_NAME \
  --service-name $SERVICE_NAME \
  --task-definition expense-tracker \
  --desired-count 2 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxxxx,subnet-yyyyy],securityGroups=[sg-xxxxx],assignPublicIp=ENABLED}" \
  --load-balancers targetGroupArn=$TG_ARN,containerName=expense-tracker,containerPort=8080 \
  --region $AWS_REGION
```

## Step 8: Setup Auto Scaling

### 8.1 Create Auto Scaling Target
```bash
aws application-autoscaling register-scalable-target \
  --service-namespace ecs \
  --resource-id service/$CLUSTER_NAME/$SERVICE_NAME \
  --scalable-dimension ecs:service:DesiredCount \
  --min-capacity 2 \
  --max-capacity 10 \
  --region $AWS_REGION
```

### 8.2 Create Scaling Policy
```bash
aws application-autoscaling put-scaling-policy \
  --policy-name expense-tracker-scaling \
  --service-namespace ecs \
  --resource-id service/$CLUSTER_NAME/$SERVICE_NAME \
  --scalable-dimension ecs:service:DesiredCount \
  --policy-type TargetTrackingScaling \
  --target-tracking-scaling-policy-configuration file://scaling-policy.json
```

### scaling-policy.json
```json
{
  "TargetValue": 70.0,
  "PredefinedMetricSpecification": {
    "PredefinedMetricType": "ECSServiceAverageCPUUtilization"
  },
  "ScaleOutCooldown": 300,
  "ScaleInCooldown": 300
}
```

## Step 9: Setup CloudWatch Monitoring

### 9.1 Create Log Group
```bash
aws logs create-log-group \
  --log-group-name /ecs/expense-tracker \
  --region $AWS_REGION
```

### 9.2 Create CloudWatch Alarms
```bash
aws cloudwatch put-metric-alarm \
  --alarm-name expense-tracker-cpu-high \
  --alarm-description "Alert when CPU is high" \
  --metric-name CPUUtilization \
  --namespace AWS/ECS \
  --statistic Average \
  --period 300 \
  --threshold 80 \
  --comparison-operator GreaterThanThreshold \
  --dimensions Name=ServiceName,Value=$SERVICE_NAME Name=ClusterName,Value=$CLUSTER_NAME \
  --region $AWS_REGION
```

## Monitoring

### View Metrics
```bash
# Get service details
aws ecs describe-services \
  --cluster $CLUSTER_NAME \
  --services $SERVICE_NAME \
  --region $AWS_REGION

# View logs
aws logs tail /ecs/expense-tracker --follow
```

## Troubleshooting

### Check Task Status
```bash
aws ecs list-tasks \
  --cluster $CLUSTER_NAME \
  --region $AWS_REGION

aws ecs describe-tasks \
  --cluster $CLUSTER_NAME \
  --tasks <task-arn> \
  --region $AWS_REGION
```

### View Container Logs
```bash
aws logs get-log-events \
  --log-group-name /ecs/expense-tracker \
  --log-stream-name ecs/expense-tracker/<container-id> \
  --region $AWS_REGION
```

## Cleanup

To delete resources:
```bash
# Delete service
aws ecs delete-service \
  --cluster $CLUSTER_NAME \
  --service $SERVICE_NAME \
  --force

# Delete cluster
aws ecs delete-cluster \
  --cluster $CLUSTER_NAME

# Delete ALB
aws elbv2 delete-load-balancer --load-balancer-arn $ALB_ARN

# Delete target group
aws elbv2 delete-target-group --target-group-arn $TG_ARN

# Delete RDS instance
aws rds delete-db-instance \
  --db-instance-identifier expense-tracker-db \
  --skip-final-snapshot
```

## Monitoring Endpoints

Once deployed, access:
- API: `http://<alb-dns>/api/`
- Metrics: `http://<alb-dns>/metrics`
- Health: `http://<alb-dns>/health`
