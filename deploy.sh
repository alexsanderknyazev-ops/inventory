#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_step() { echo -e "${GREEN}▶${NC} $1"; }
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }

echo "========================================"
echo "🚀 MARKET APP DEPLOYMENT (CLEAN UPDATE)"
echo "========================================"

# 1. Проверка Minikube
print_step "Checking Minikube..."
if ! minikube status | grep -q "Running"; then
    print_error "Minikube is not running."
    exit 1
fi

# 2. Настройка Docker
print_step "Setting up Docker..."
eval $(minikube docker-env)
print_success "Docker configured"

# 3. Сборка образа
print_step "Building Docker image..."
if docker build -t market-app:latest .; then
    print_success "Image built"
else
    print_error "Build failed"
    exit 1
fi

# 4. Проверка namespace
print_step "Checking namespace..."
if ! kubectl get namespace market >/dev/null 2>&1; then
    print_error "Namespace 'market' not found"
    exit 1
fi

# 5. Удаляем ВСЕ deployment с префиксом "market-app"
print_step "Cleaning up OLD market-app deployments..."
# Находим все deployment с префиксом market-app
MARKET_DEPLOYMENTS=$(kubectl get deployments -n market --no-headers 2>/dev/null | awk '{print $1}' | grep "^market-app")

if [ -n "$MARKET_DEPLOYMENTS" ]; then
    echo "Found market-app deployments to delete:"
    for DEPLOY in $MARKET_DEPLOYMENTS; do
        echo "  - $DEPLOY"
        kubectl delete deployment -n market "$DEPLOY" --ignore-not-found
    done
    print_success "Old market-app deployments deleted"
    
    # Ждем пока старые поды удалятся
    echo "Waiting for old pods to terminate..."
    sleep 5
    
    # Проверяем что старые поды удалены
    OLD_PODS=$(kubectl get pods -n market --no-headers 2>/dev/null | grep -E "(market-app|market-deployment)" | wc -l)
    if [ "$OLD_PODS" -gt 0 ]; then
        echo "Force deleting remaining old pods..."
        kubectl delete pods -n market -l 'app in (market, market-app)' --ignore-not-found
        sleep 2
    fi
else
    print_success "No old market-app deployments found"
fi

# 6. Создаём новый deployment с именем market-app-v[timestamp]
TIMESTAMP=$(date +%Y%m%d%H%M%S)
NEW_DEPLOYMENT_NAME="market-app-v${TIMESTAMP}"

print_step "Creating NEW deployment: ${NEW_DEPLOYMENT_NAME}..."
cat <<YAML | kubectl apply -n market -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${NEW_DEPLOYMENT_NAME}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: market-app
      version: "v${TIMESTAMP}"
  template:
    metadata:
      labels:
        app: market-app
        version: "v${TIMESTAMP}"
    spec:
      containers:
      - name: market-app
        image: market-app:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 8070
        env:
        - name: APP_PORT
          value: "8070"
        - name: DB_HOST
          value: "postgres"
        - name: DB_PORT
          value: "5432"
        - name: DB_NAME
          value: "marketdb"
        - name: DB_USER
          value: "admin"
        - name: DB_PASSWORD
          value: "admin123"
        readinessProbe:
          httpGet:
            path: /inventory/health
            port: 8070
          initialDelaySeconds: 10
          periodSeconds: 5
YAML
print_success "Deployment ${NEW_DEPLOYMENT_NAME} created"

# 7. Service (создаем или обновляем)
print_step "Creating/Updating service..."
cat <<YAML | kubectl apply -n market -f -
apiVersion: v1
kind: Service
metadata:
  name: market-service
spec:
  selector:
    app: market-app
    version: "v${TIMESTAMP}"
  ports:
  - port: 8070
    targetPort: 8070
YAML
print_success "Service ready (now pointing to version v${TIMESTAMP})"

# 8. Ждём запуска новой поды
print_step "Waiting for NEW pod..."
MAX_WAIT=60
POD_READY=false
for i in $(seq 1 $MAX_WAIT); do
    POD_NAME=$(kubectl get pods -n market -l version=v${TIMESTAMP} -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    
    if [ -n "$POD_NAME" ]; then
        POD_STATUS=$(kubectl get pod -n market "$POD_NAME" -o jsonpath='{.status.phase}' 2>/dev/null)
        POD_READY_STATE=$(kubectl get pod -n market "$POD_NAME" -o jsonpath='{.status.containerStatuses[0].ready}' 2>/dev/null)
        
        if [[ "$POD_STATUS" == "Running" ]] && [[ "$POD_READY_STATE" == "true" ]]; then
            print_success "✅ New pod $POD_NAME is running and ready!"
            POD_READY=true
            break
        fi
    fi
    
    if [ $i -eq $MAX_WAIT ]; then
        print_error "❌ Timeout waiting for pod"
        echo "Current pods:"
        kubectl get pods -n market
        echo ""
        echo "Checking deployment status:"
        kubectl describe deployment -n market ${NEW_DEPLOYMENT_NAME} | tail -20
        exit 1
    fi
    
    echo -n "."
    sleep 1
done
echo ""

# 9. Проверка статуса
print_step "Checking status..."
echo ""
echo "📊 CURRENT PODS:"
kubectl get pods -n market -o wide
echo ""
echo "📊 CURRENT DEPLOYMENTS:"
kubectl get deployments -n market

# 10. Тест
print_step "Testing application..."
kubectl port-forward -n market svc/market-service 8070:8070 > /dev/null 2>&1 &
PF_PID=$!
sleep 3

if curl -s --max-time 5 http://localhost:8070/inventory/health > /dev/null 2>&1; then
    print_success "✅ App is working!"
    echo "   Health response:"
    curl -s http://localhost:8070/inventory/health | head -c 100
    echo "..."
else
    print_error "❌ App not responding"
    echo "Checking logs..."
    kubectl logs -n market -l version=v${TIMESTAMP} --tail=10
fi

kill $PF_PID 2>/dev/null

# 11. Удаляем старые deployment с префиксом market-app (кроме текущего)
print_step "Final cleanup of other market-app deployments..."
OTHER_DEPLOYMENTS=$(kubectl get deployments -n market --no-headers 2>/dev/null | awk '{print $1}' | grep "^market-app" | grep -v "^${NEW_DEPLOYMENT_NAME}$")

if [ -n "$OTHER_DEPLOYMENTS" ]; then
    echo "Found other market-app deployments to clean up:"
    for DEPLOY in $OTHER_DEPLOYMENTS; do
        echo "  - $DEPLOY"
        kubectl delete deployment -n market "$DEPLOY" --ignore-not-found
    done
    print_success "Other market-app deployments cleaned"
else
    print_success "No other market-app deployments found"
fi

echo ""
echo "========================================"
echo "✅ DEPLOYMENT COMPLETE!"
echo "========================================"
echo ""
echo "📌 Summary:"
echo "   • New Deployment: ${NEW_DEPLOYMENT_NAME}"
echo "   • Version: v${TIMESTAMP}"
echo "   • Image: market-app:latest"
echo "   • Old market-app deployments: Removed"
echo ""
echo "🌐 Access from Postman:"
echo "   1. kubectl port-forward -n market svc/market-service 8070:8070"
echo "   2. Use: http://localhost:8070"
echo ""
echo "📊 Current resources:"
echo "   kubectl get pods -n market"
echo "   kubectl get deployments -n market"
echo ""
echo "🔄 To update again:"
echo "   Just run this script!"
echo ""