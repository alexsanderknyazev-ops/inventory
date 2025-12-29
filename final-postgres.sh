#!/bin/bash

echo "🚀 FINAL PostgreSQL Setup for Kubernetes"

# 1. Очистка
echo "🗑️  Cleaning up old deployment..."
kubectl delete namespace market --ignore-not-found 2>/dev/null
sleep 3

# 2. Создание namespace
echo "📁 Creating namespace..."
kubectl create namespace market

# 3. Развертывание PostgreSQL
echo "📦 Deploying PostgreSQL..."
cat <<YAML | kubectl apply -n market -f -
apiVersion: v1
kind: Service
metadata:
  name: postgres
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15
        env:
        - name: POSTGRES_USER
          value: "admin"
        - name: POSTGRES_PASSWORD
          value: "admin123"
        - name: POSTGRES_DB
          value: "marketdb"
        ports:
        - containerPort: 5432
        readinessProbe:
          exec:
            command:
            - sh
            - -c
            - "pg_isready -U admin -d marketdb"
          initialDelaySeconds: 15
          periodSeconds: 5
YAML

# 4. Ожидание
echo "⏳ Waiting 70 seconds for PostgreSQL to fully initialize..."
echo "   (This can take a while on the first run)"
for i in {1..35}; do
  echo -n "."
  sleep 2
done
echo ""

# 5. Проверка статуса
echo ""
echo "📊 CURRENT STATUS:"
kubectl get pods,svc -n market

# 6. Проверка логов PostgreSQL
echo ""
echo "📋 PostgreSQL logs (last 10 lines):"
kubectl logs -n market -l app=postgres --tail=10 2>/dev/null || echo "Logs not available yet"

# 7. Тестирование подключения ВНУТРИ Kubernetes
echo ""
echo "🔍 TEST 1: Connecting from another pod in Kubernetes..."
echo "   This test runs INSIDE the cluster, no port forwarding needed"
echo "   -------------------------------------------------------------"

kubectl run -n market postgres-test --rm -i --restart=Never --image=postgres:15 -- \
  bash -c "
    echo 'Waiting 5 seconds...'
    sleep 5
    echo 'Testing connection to postgres service...'
    PGPASSWORD=admin123 psql -h postgres -U admin -d marketdb -c 'SELECT version(), current_timestamp;'
    EXIT_CODE=\$?
    if [ \$EXIT_CODE -eq 0 ]; then
      echo '✅ SUCCESS! PostgreSQL is working inside Kubernetes!'
    else
      echo '❌ Connection failed inside Kubernetes'
    fi
    exit \$EXIT_CODE
  " 2>&1

TEST_RESULT=$?

# 8. Альтернативный тест
if [ $TEST_RESULT -ne 0 ]; then
  echo ""
  echo "🔍 TEST 2: Alternative test with simpler command..."
  timeout 15 kubectl run -n market test2 --rm --restart=Never --image=postgres:15 -- \
    sh -c "PGPASSWORD=admin123 psql -h postgres -U admin -d marketdb -c 'SELECT 1;'" 2>&1
fi

# 9. Вывод информации
echo ""
echo "================================================"
echo "✅ DEPLOYMENT COMPLETE!"
echo "================================================"
echo ""
echo "📌 Connection details INSIDE Kubernetes:"
echo "   Service:    postgres.market.svc.cluster.local"
echo "   Port:       5432"
echo "   Database:   marketdb"
echo "   Username:   admin"
echo "   Password:   admin123"
echo ""
echo "📌 To use PostgreSQL from YOUR APPLICATIONS in Kubernetes:"
echo "   Use this connection string in your app config:"
echo "   postgresql://admin:admin123@postgres.market:5432/marketdb"
echo ""
echo "�� Troubleshooting commands:"
echo "   kubectl logs -n market -l app=postgres"
echo "   kubectl describe pod -n market -l app=postgres"
echo "   kubectl get endpoints -n market postgres"
echo ""
echo "🌐 If you need external access (outside Kubernetes):"
echo "   Option A: Stop local PostgreSQL and use port forwarding:"
echo "     brew services stop postgresql"
echo "     kubectl port-forward -n market svc/postgres 5432:5432"
echo ""
echo "   Option B: Use a different port:"
echo "     kubectl port-forward -n market svc/postgres 5433:5432"
echo "     Then connect to: localhost:5433"
echo ""
echo "   Option C: Use NodePort service type (not recommended for production)"
echo "================================================"
