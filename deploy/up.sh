#!/bin/bash

# Configuration Variables
CLUSTER_NAME="llm-inference-cluster"
REGION="ap-northeast-1"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
POLICY_NAME="AWSLoadBalancerControllerIAMPolicy"

echo "Checking Infrastructure..."

# 1. Create EKS Cluster
echo ">>> Step 1: Creating EKS Cluster (this may take 15-20 minutes)..."
eksctl create cluster -f infra/eks-cluster.yaml

# 2. Setup IAM Policy for ALB Controller
echo ">>> Step 2: Setting up IAM Policy..."
# Only create policy if it doesn't exist
POLICY_ARN=$(aws iam list-policies --query "Policies[?PolicyName=='$POLICY_NAME'].Arn" --output text)
if [ -z "$POLICY_ARN" ]; then
    echo "Creating new IAM Policy..."
    POLICY_ARN=$(aws iam create-policy --policy-name $POLICY_NAME --policy-document file://iam_policy.json --query 'Policy.Arn' --output text)
else
    echo "IAM Policy already exists: $POLICY_ARN"
fi

# 3. Create IAM Service Account (IRSA)
echo ">>> Step 3: Creating IAM Service Account (IRSA)..."
eksctl create iamserviceaccount \
  --cluster=$CLUSTER_NAME \
  --namespace=kube-system \
  --name=aws-load-balancer-controller \
  --role-name AmazonEKSLoadBalancerControllerRole \
  --attach-policy-arn=$POLICY_ARN \
  --approve --override-existing-serviceaccounts

# 4. Install AWS Load Balancer Controller via Helm
echo ">>> Step 4: Installing AWS Load Balancer Controller..."
helm repo add eks https://aws.github.io/eks-charts
helm repo update
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=$CLUSTER_NAME \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller

# 5. Deploy Applications
echo ">>> Step 5: Deploying Apps (vLLM, Gate Service, HPA, Ingress)..."
kubectl apply -f k8s/aws/fsx-pv.yaml
kubectl apply -f k8s/aws/fsx-pvc.yaml
kubectl apply -f k8s/aws/vllm-gpu.yaml
kubectl apply -f k8s/base/gate-service.yaml
kubectl apply -f k8s/aws/gate-hpa.yaml
kubectl apply -f k8s/aws/ingress.yaml

echo "🎉 Deployment Sync Completed!"