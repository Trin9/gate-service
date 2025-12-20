#!/bin/bash
# EKS Infrastructure Termination Script (Cost Optimization)

CLUSTER_NAME="llm-inference-cluster"

echo "⚠️ WARNING: Starting complete teardown of $CLUSTER_NAME..."

# 1. Delete Ingress first to trigger ALB deletion
# This is crucial because ALB is managed by the controller, not eksctl
echo "🛑 Step 1: Deleting Ingress (ALB)..."
kubectl delete ingress gate-ingress --ignore-not-found

# 2. Wait a moment for ALB deletion to initiate
sleep 30

# 3. Destroy the entire EKS cluster
# This removes NodeGroups, Fargate profiles, and the Control Plane
echo "💣 Step 2: Destroying EKS Cluster (This takes ~15 mins)..."
eksctl delete cluster --name $CLUSTER_NAME --force

# 4. Final Cleanup Reminders
echo "----------------------------------------------------"
echo "📢 Automated teardown submitted. Please manual check:"
echo "1. EC2 Volumes: Delete any 'available' volumes (AWS Console -> EC2 -> Volumes)."
echo "2. FSx Filesystem: Ensure FSx is deleted (AWS Console -> FSx)."
echo "3. CloudFormation: Confirm stacks are 'DELETE_COMPLETE'."
echo "4. IAM Policy: AWSLoadBalancerControllerIAMPolicy can be kept for next time."
echo "----------------------------------------------------"
echo "✅ Teardown process initiated successfully."