# ☁️ Cloud Deployment Guide

This guide provides comprehensive deployment strategies for the Transaction Tracker application on AWS and Google Cloud Platform.

## Platform Comparison

### AWS vs Google Cloud Overview

| Aspect                 | AWS Solution                      | Google Cloud Solution   |
| ---------------------- | --------------------------------- | ----------------------- |
| **Container Service**  | ECS Fargate                       | Cloud Run               |
| **Database**           | RDS MySQL                         | Cloud SQL MySQL         |
| **Cache**              | ElastiCache Redis                 | Memorystore Redis       |
| **Load Balancer**      | Application Load Balancer         | Cloud Load Balancing    |
| **Container Registry** | ECR                               | Artifact Registry       |
| **Network**            | VPC + Subnets                     | VPC + Subnets           |
| **Monitoring**         | CloudWatch                        | Cloud Monitoring        |
| **Cost Optimization**  | Spot instances, reserved capacity | Sustained use discounts |

## Architecture Diagrams

### AWS Architecture

```
Internet → ALB → ECS Fargate (3 services) → RDS MySQL
                              ↓
                         ElastiCache Redis
```

### Google Cloud Architecture

```
Internet → Load Balancer → Cloud Run (3 services) → Cloud SQL MySQL
                                    ↓
                              Memorystore Redis
```

## Detailed Component Specifications

### 1. Container Services

| Component         | AWS ECS Fargate              | Google Cloud Run              |
| ----------------- | ---------------------------- | ----------------------------- |
| **Frontend**      | React app in nginx container | React app in nginx container  |
| **Backend API**   | Go API service               | Go API service                |
| **Price Service** | Go microservice              | Go microservice               |
| **Auto Scaling**  | Target tracking scaling      | Automatic concurrency scaling |
| **CPU/Memory**    | 0.25 vCPU, 512MB each        | 1 vCPU, 512MB each            |
| **Networking**    | VPC with private subnets     | VPC connector                 |

### 2. Database & Cache

| Component    | AWS                                | Google Cloud                      |
| ------------ | ---------------------------------- | --------------------------------- |
| **Database** | RDS MySQL 8.0 (db.t3.micro)        | Cloud SQL MySQL 8.0 (db-f1-micro) |
| **Storage**  | 20GB GP2 SSD                       | 20GB SSD                          |
| **Cache**    | ElastiCache Redis (cache.t3.micro) | Memorystore Redis (1GB)           |
| **Backup**   | 7-day automated backup             | 7-day automated backup            |
| **Multi-AZ** | Single AZ (cost optimized)         | Single zone (cost optimized)      |

### 3. Networking & Security

| Component           | AWS                                    | Google Cloud              |
| ------------------- | -------------------------------------- | ------------------------- |
| **Load Balancer**   | Application Load Balancer              | Global Load Balancer      |
| **SSL Certificate** | AWS Certificate Manager (free)         | Google-managed SSL (free) |
| **Domain**          | Route 53 ($0.50/month)                 | Cloud DNS ($0.20/month)   |
| **VPC**             | Custom VPC with public/private subnets | Custom VPC with subnets   |
| **Security Groups** | ECS security groups                    | Firewall rules            |

## Deployment Steps

### AWS Deployment

#### 1. Infrastructure Setup

```bash
# Create VPC and networking
aws ec2 create-vpc --cidr-block 10.0.0.0/16
aws ec2 create-subnet --vpc-id <vpc-id> --cidr-block 10.0.1.0/24
aws ec2 create-internet-gateway
```

#### 2. Database Creation

```bash
# Create RDS MySQL instance
aws rds create-db-instance \
  --db-instance-identifier transaction-tracker-db \
  --db-instance-class db.t3.micro \
  --engine mysql \
  --master-username admin \
  --allocated-storage 20
```

#### 3. Cache Setup

```bash
# Create ElastiCache Redis
aws elasticache create-cache-cluster \
  --cache-cluster-id transaction-tracker-redis \
  --cache-node-type cache.t3.micro \
  --engine redis
```

#### 4. Container Registry

```bash
# Create ECR repositories
aws ecr create-repository --repository-name transaction-tracker/frontend
aws ecr create-repository --repository-name transaction-tracker/backend
aws ecr create-repository --repository-name transaction-tracker/price-service
```

#### 5. Build and Push Images

```bash
# Frontend
docker build -t transaction-tracker/frontend ./frontend
docker tag transaction-tracker/frontend:latest <account>.dkr.ecr.us-east-1.amazonaws.com/transaction-tracker/frontend:latest
docker push <account>.dkr.ecr.us-east-1.amazonaws.com/transaction-tracker/frontend:latest

# Backend
docker build -t transaction-tracker/backend ./backend
docker tag transaction-tracker/backend:latest <account>.dkr.ecr.us-east-1.amazonaws.com/transaction-tracker/backend:latest
docker push <account>.dkr.ecr.us-east-1.amazonaws.com/transaction-tracker/backend:latest

# Price Service
docker build -t transaction-tracker/price-service ./price_service
docker tag transaction-tracker/price-service:latest <account>.dkr.ecr.us-east-1.amazonaws.com/transaction-tracker/price-service:latest
docker push <account>.dkr.ecr.us-east-1.amazonaws.com/transaction-tracker/price-service:latest
```

#### 6. ECS Cluster Setup

```bash
# Create ECS cluster
aws ecs create-cluster --cluster-name transaction-tracker-cluster

# Create task definitions for each service
aws ecs register-task-definition --cli-input-json file://frontend-task-definition.json
aws ecs register-task-definition --cli-input-json file://backend-task-definition.json
aws ecs register-task-definition --cli-input-json file://price-service-task-definition.json
```

#### 7. Load Balancer Setup

```bash
# Create Application Load Balancer
aws elbv2 create-load-balancer \
  --name transaction-tracker-alb \
  --subnets <subnet-1> <subnet-2> \
  --security-groups <security-group-id>
```

#### 8. Service Deployment

```bash
# Deploy services
aws ecs create-service \
  --cluster transaction-tracker-cluster \
  --service-name frontend-service \
  --task-definition frontend:1 \
  --desired-count 1 \
  --launch-type FARGATE
```

### Google Cloud Deployment

#### 1. Project Setup

```bash
# Set project and enable APIs
gcloud config set project <your-project-id>
gcloud services enable run.googleapis.com
gcloud services enable sql-component.googleapis.com
gcloud services enable redis.googleapis.com
```

#### 2. Database Creation

```bash
# Create Cloud SQL instance
gcloud sql instances create transaction-tracker-db \
  --database-version=MYSQL_8_0 \
  --tier=db-f1-micro \
  --region=us-central1 \
  --storage-size=20GB \
  --storage-type=SSD
```

#### 3. Cache Setup

```bash
# Create Memorystore Redis
gcloud redis instances create transaction-tracker-redis \
  --size=1 \
  --region=us-central1 \
  --redis-version=redis_6_x
```

#### 4. Container Registry

```bash
# Configure Artifact Registry
gcloud artifacts repositories create transaction-tracker \
  --repository-format=docker \
  --location=us-central1
```

#### 5. Build and Push Images

```bash
# Frontend
docker build -t us-central1-docker.pkg.dev/<project>/transaction-tracker/frontend ./frontend
docker push us-central1-docker.pkg.dev/<project>/transaction-tracker/frontend

# Backend
docker build -t us-central1-docker.pkg.dev/<project>/transaction-tracker/backend ./backend
docker push us-central1-docker.pkg.dev/<project>/transaction-tracker/backend

# Price Service
docker build -t us-central1-docker.pkg.dev/<project>/transaction-tracker/price-service ./price_service
docker push us-central1-docker.pkg.dev/<project>/transaction-tracker/price-service
```

#### 6. Cloud Run Deployment

```bash
# Deploy frontend
gcloud run deploy frontend \
  --image=us-central1-docker.pkg.dev/<project>/transaction-tracker/frontend \
  --platform=managed \
  --region=us-central1 \
  --allow-unauthenticated \
  --memory=512Mi \
  --cpu=1

# Deploy backend
gcloud run deploy backend \
  --image=us-central1-docker.pkg.dev/<project>/transaction-tracker/backend \
  --platform=managed \
  --region=us-central1 \
  --allow-unauthenticated \
  --memory=512Mi \
  --cpu=1 \
  --set-env-vars="DB_HOST=<cloud-sql-ip>,REDIS_HOST=<redis-ip>"

# Deploy price service
gcloud run deploy price-service \
  --image=us-central1-docker.pkg.dev/<project>/transaction-tracker/price-service \
  --platform=managed \
  --region=us-central1 \
  --allow-unauthenticated \
  --memory=512Mi \
  --cpu=1 \
  --set-env-vars="REDIS_HOST=<redis-ip>"
```

#### 7. Load Balancer Setup

```bash
# Create load balancer
gcloud compute backend-services create frontend-backend \
  --global \
  --load-balancing-scheme=EXTERNAL

gcloud compute url-maps create transaction-tracker-lb \
  --default-service=frontend-backend
```

## Cost Analysis

### AWS Pricing

#### Standard Configuration

| Service                       | Specification                               | Monthly Cost |
| ----------------------------- | ------------------------------------------- | ------------ |
| **ECS Fargate**               | 3 services × 0.25 vCPU × 512MB × 730 hours  | $13.14       |
| **RDS MySQL**                 | db.t3.micro (1 vCPU, 1GB RAM, 20GB storage) | $17.00       |
| **ElastiCache Redis**         | cache.t3.micro (2 vCPU, 0.5GB RAM)          | $11.50       |
| **Application Load Balancer** | 1 ALB + minimal data processing             | $16.20       |
| **Route 53**                  | Hosted zone                                 | $0.50        |
| **Data Transfer**             | ~10GB outbound                              | $1.00        |
| **CloudWatch Logs**           | Basic monitoring                            | $2.00        |
| **NAT Gateway**               | Single NAT for private subnets              | $32.40       |
| **Total**                     |                                             | **$93.74**   |

#### Cost Optimized Configuration

Remove NAT Gateway, use public subnets:
| Service | Monthly Cost |
| ----------------- | ------------ |
| **Core Services** | $61.34 |
| **Total** | **$61.34** |

### Google Cloud Pricing

#### Standard Configuration

| Service                  | Specification                                  | Monthly Cost |
| ------------------------ | ---------------------------------------------- | ------------ |
| **Cloud Run**            | 3 services × 1 vCPU × 512MB × minimal requests | $8.00        |
| **Cloud SQL MySQL**      | db-f1-micro (0.6GB RAM, 20GB SSD)              | $12.50       |
| **Memorystore Redis**    | 1GB standard tier                              | $26.00       |
| **Cloud Load Balancing** | Global LB + forwarding rules                   | $18.00       |
| **Cloud DNS**            | Hosted zone                                    | $0.20        |
| **Data Transfer**        | ~10GB egress                                   | $1.20        |
| **Cloud Monitoring**     | Basic monitoring                               | $0.50        |
| **Total**                |                                                | **$66.40**   |

#### Cost Optimized Configuration

Use Cloud Memorystore Basic tier:
| Service | Monthly Cost |
| --------------------- | ------------ |
| **Cloud Run** | $8.00 |
| **Cloud SQL** | $12.50 |
| **Memorystore Basic** | $15.00 |
| **Load Balancing** | $18.00 |
| **Other** | $2.00 |
| **Total** | **$55.50** |

## Ultra Cost-Optimized Solutions

### AWS Lightsail Alternative

| Service                 | Specification                     | Monthly Cost |
| ----------------------- | --------------------------------- | ------------ |
| **Lightsail Container** | 512MB RAM, 0.25 vCPU × 3 services | $21.00       |
| **Lightsail Database**  | 1GB RAM, 40GB storage             | $15.00       |
| **Redis**               | Self-hosted in container          | $0.00        |
| **Total**               |                                   | **$36.00**   |

### Google Cloud Minimal Setup

| Service            | Specification            | Monthly Cost |
| ------------------ | ------------------------ | ------------ |
| **Cloud Run**      | Minimal CPU allocation   | $5.00        |
| **Cloud SQL**      | Shared-core instance     | $9.50        |
| **Redis**          | Self-hosted in Cloud Run | $1.00        |
| **Load Balancing** | Basic forwarding         | $5.00        |
| **Total**          |                          | **$20.50**   |

## Deployment Recommendations

### For Personal Use (Cost Priority)

**Google Cloud Minimal Setup** - **$20.50/month**

- ✅ Serverless scaling with pay-per-use pricing
- ✅ Managed services reduce operational overhead
- ✅ Automatic scaling based on traffic
- ❌ Cold starts may affect performance
- ❌ Shared resources with potential limitations

### For Production (Performance Priority)

**AWS Cost-Optimized** - **$61.34/month**

- ✅ Dedicated resources for consistent performance
- ✅ Extensive monitoring and debugging tools
- ✅ Better control over infrastructure
- ❌ Higher operational costs
- ❌ More complex setup and maintenance

### Hybrid Approach

Start with Google Cloud minimal setup for development and testing, then migrate to AWS when scaling requirements and budget allow for production workloads.

## Infrastructure as Code

### Terraform Templates

Both AWS and Google Cloud deployments can be automated using Infrastructure as Code. Example Terraform configurations are available for:

- **AWS**: ECS Fargate, RDS, ElastiCache, ALB setup
- **Google Cloud**: Cloud Run, Cloud SQL, Memorystore, Load Balancer setup

Contact the repository maintainer for complete Terraform templates and deployment automation scripts.

## Monitoring and Maintenance

### AWS Monitoring

- CloudWatch for application logs and metrics
- X-Ray for distributed tracing
- Systems Manager for container management

### Google Cloud Monitoring

- Cloud Monitoring for metrics and alerting
- Cloud Logging for centralized log management
- Cloud Trace for application performance insights

## Security Considerations

### Network Security

- VPC with properly configured subnets and security groups
- SSL/TLS termination at load balancer level
- Database and cache access restricted to application services

### Application Security

- Environment variables for sensitive configuration
- Regular security updates for container images
- IAM roles with least privilege principles

## Support and Troubleshooting

For deployment issues, refer to:

1. Cloud provider documentation
2. Application logs in the monitoring systems
3. GitHub Issues for application-specific problems
4. Community forums for platform-specific questions
