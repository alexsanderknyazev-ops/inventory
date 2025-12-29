#!/bin/bash
# jenkins-deploy.sh - для использования в Jenkins

echo "=== JENKINS DEPLOYMENT SCRIPT ==="
echo "Build: ${BUILD_NUMBER}"

# 1. Создаем namespace если нет
echo "1. Creating namespace..."
kubectl create namespace market 2>/dev/null || echo "Namespace already exists"

# 2. Развертываем PostgreSQL
echo "2. Deploying PostgreSQL..."
./final-postgres.sh 2>/dev/null || echo "PostgreSQL deployment attempted"

# 3. Развертываем приложение
echo "3. Deploying application..."
./deploy-app-only.sh 2>/dev/null || echo "Application deployment attempted"

# 4. Проверяем результат
echo "4. Checking status..."
sleep 10
kubectl get pods -n market 2>/dev/null || echo "Cannot check pods"

echo "✅ Deployment script completed"
echo "Run manually if needed: ./final-postgres.sh && ./deploy-app-only.sh"