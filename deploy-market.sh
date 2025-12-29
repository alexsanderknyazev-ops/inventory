#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${GREEN}▶${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

echo "========================================"
echo "🚀 MARKET APPLICATION DEPLOYMENT"
echo "========================================"

# 1. Проверка Minikube
print_step "Checking Minikube status..."
if ! minikube status | grep -q "Running"; then
    print_warning "Minikube is not running. Starting Minikube..."
    minikube start --memory=4096 --cpus=2
    sleep 10
else
    print_success "Minikube is already running"
fi

# 2. Настройка Docker
print_step "Setting up Docker environment..."
eval $(minikube docker-env)
print_success "Docker environment configured"

# 3. Сборка Docker образа
print_step "Building Docker image..."
if docker build -t market-app:latest .; then
    print_success "Docker image built successfully"
else
    print_error "Failed to build Docker image"
    exit 1
fi

# 4. Очистка старого приложения
print_step "Cleaning up old deployment..."
kubectl delete deployment -n market market-deployment --ignore-not-found 2>/dev/null
kubectl delete service -n market market-service --ignore-not-found 2>/dev/null
sleep 3
print_success "Old deployment cleaned"

# 5. Создание namespace если нет
print_step "Checking namespace..."
if ! kubectl get namespace market >/dev/null 2>&1; then
    kubectl create namespace market
    print_success "Namespace 'market' created"
else
    print_success "Namespace 'market' already exists"
fi

# 6. Проверка и настройка ConfigMap
print_step "Setting up ConfigMap..."
cat <<YAML | kubectl apply -n market -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: market-config
data:
  APP_PORT: "8070"
  APP_ENV: "development"
  LOG_LEVEL: "info"
  DB_HOST: "postgres.market.svc.cluster.local"
  DB_PORT: "5432"
  DB_NAME: "marketdb"
  DB_USER: "admin"
  DB_SSLMODE: "disable"
  GIN_MODE: "debug"
  READ_TIMEOUT: "30"
  WRITE_TIMEOUT: "30"
  DB_AUTO_MIGRATE: "true"
YAML
print_success "ConfigMap created"

# 7. Проверка и настройка Secret
print_step "Setting up Secret..."
cat <<YAML | kubectl apply -n market -f -
apiVersion: v1
kind: Secret
metadata:
  name: market-secret
type: Opaque
stringData:
  DB_PASSWORD: "admin123"
YAML
print_success "Secret created"

# 8. Проверка PostgreSQL
print_step "Checking PostgreSQL..."
if ! kubectl get pods -n market -l app=postgres 2>/dev/null | grep -q "postgres"; then
    print_warning "PostgreSQL not found. Deploying PostgreSQL..."
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
YAML
    echo "⏳ Waiting 60 seconds for PostgreSQL to start..."
    sleep 60
    print_success "PostgreSQL deployed"
else
    print_success "PostgreSQL is already running"
fi

# 9. Развертывание приложения
print_step "Deploying Market Application..."
cat <<YAML | kubectl apply -n market -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: market-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: market
  template:
    metadata:
      labels:
        app: market
    spec:
      containers:
      - name: market-app
        image: market-app:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8070
        env:
        - name: APP_PORT
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: APP_PORT
        - name: APP_ENV
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: APP_ENV
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: LOG_LEVEL
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: DB_HOST
        - name: DB_PORT
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: DB_PORT
        - name: DB_NAME
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: DB_NAME
        - name: DB_USER
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: DB_USER
        - name: DB_SSLMODE
          valueFrom:
            configMapKeyRef:
              name: market-config
              key: DB_SSLMODE
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: market-secret
              key: DB_PASSWORD
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        readinessProbe:
          httpGet:
            path: /inventory/health
            port: 8070
          initialDelaySeconds: 20
          periodSeconds: 5
          timeoutSeconds: 3
        livenessProbe:
          httpGet:
            path: /inventory/health
            port: 8070
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 3
---
apiVersion: v1
kind: Service
metadata:
  name: market-service
spec:
  selector:
    app: market
  ports:
  - name: http
    port: 8070
    targetPort: 8070
    protocol: TCP
  type: ClusterIP
YAML
print_success "Application deployed"

# 10. Ожидание запуска
print_step "Waiting for application to start..."
echo -n "⏳ "
for i in {1..30}; do
    echo -n "."
    sleep 2
done
echo ""

# 11. Проверка статуса
print_step "Checking deployment status..."
echo ""
echo "📊 PODS STATUS:"
kubectl get pods -n market -o wide

echo ""
echo "📊 SERVICES STATUS:"
kubectl get svc -n market

echo ""
echo "📊 DEPLOYMENTS STATUS:"
kubectl get deployments -n market

# 12. Проверка логов
print_step "Checking application logs..."
APP_POD=$(kubectl get pods -n market -l app=market -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
if [ -n "$APP_POD" ]; then
    echo ""
    echo "📋 LOGS (last 20 lines):"
    kubectl logs -n market $APP_POD --tail=20 2>/dev/null || print_warning "No logs yet"
else
    print_error "No application pods found"
    exit 1
fi

# 13. Тестирование подключения
print_step "Testing application connection..."
kubectl port-forward -n market svc/market-service 8070:8070 > /dev/null 2>&1 &
PORTFORWARD_PID=$!

echo "⏳ Waiting for port forward..."
sleep 8

echo ""
echo "🧪 Testing endpoints:"
echo "  1. Health check..."

if curl -s --max-time 10 http://localhost:8070/inventory/health > /dev/null 2>&1; then
    print_success "Health check passed!"
    echo "     Response:"
    curl -s http://localhost:8070/inventory/health | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8070/inventory/health
    echo ""
else
    print_error "Health check failed"
    echo "  2. Trying root endpoint..."
    curl -s --max-time 5 http://localhost:8070/ || echo "     No response"
fi

# Останавливаем port-forward
kill $PORTFORWARD_PID 2>/dev/null

# 14. Финальный вывод
echo ""
echo "========================================"
echo "✅ DEPLOYMENT COMPLETE!"
echo "========================================"
echo ""
echo "📌 Connection Information:"
echo "   Application URL inside cluster: http://market-service.market:8070"
echo "   PostgreSQL URL inside cluster:  postgres://admin:admin123@postgres.market:5432/marketdb"
echo ""
echo "🌐 Access from local machine:"
echo "   1. kubectl port-forward -n market svc/market-service 8070:8070"
echo "   2. Open: http://localhost:8070/inventory/health"
echo ""
echo "🔧 Useful Commands:"
echo "   kubectl logs -n market -l app=market -f              # Follow logs"
echo "   kubectl get all -n market                            # Show all resources"
echo "   kubectl describe pod -n market -l app=market         # Debug pod"
echo "   kubectl exec -n market -it $APP_POD -- sh           # Enter container"
echo ""
echo "🔄 To restart deployment:"
echo "   kubectl rollout restart deployment -n market market-deployment"
echo ""
echo "========================================"
