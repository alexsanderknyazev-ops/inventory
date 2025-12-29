#!/bin/bash
# deploy-app-only.sh - Jenkins compatible version

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date '+%H:%M:%S')]${NC} $1"; }
error() { echo -e "${RED}[$(date '+%H:%M:%S')]${NC} ❌ $1"; exit 1; }

# Configuration
NAMESPACE="market"
BUILD_NUMBER="${BUILD_NUMBER:-$(date +%Y%m%d%H%M%S)}"
IMAGE_TAG="market-app:${BUILD_NUMBER}"
DEPLOYMENT_NAME="market-app-${BUILD_NUMBER}"

log "🚀 Deploying Market App (Build: ${BUILD_NUMBER})"

# Check prerequisites
log "1. Checking prerequisites..."
if ! kubectl get namespace ${NAMESPACE} >/dev/null 2>&1; then
    error "Namespace ${NAMESPACE} not found"
fi

if ! kubectl get deployment -n ${NAMESPACE} postgres >/dev/null 2>&1; then
    error "PostgreSQL not found in namespace ${NAMESPACE}"
fi

# Clean up old deployments
log "2. Cleaning old deployments..."
kubectl get deployments -n ${NAMESPACE} --no-headers 2>/dev/null | \
    awk '{print $1}' | grep "^market-app" | \
    while read DEPLOY; do
        log "   Deleting: $DEPLOY"
        kubectl delete deployment -n ${NAMESPACE} "$DEPLOY" --ignore-not-found
    done
sleep 2

# Apply ConfigMap and Secret
log "3. Applying configuration..."
kubectl apply -n ${NAMESPACE} -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: market-config
data:
  APP_PORT: "8070"
  DB_HOST: "postgres"
  DB_PORT: "5432"
  DB_NAME: "marketdb"
  DB_USER: "admin"
  DB_SSLMODE: "disable"
---
apiVersion: v1
kind: Secret
metadata:
  name: market-secret
type: Opaque
stringData:
  DB_PASSWORD: "admin123"
EOF

# Deploy new version
log "4. Deploying ${DEPLOYMENT_NAME}..."
kubectl apply -n ${NAMESPACE} -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${DEPLOYMENT_NAME}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: market-app
      version: "${BUILD_NUMBER}"
  template:
    metadata:
      labels:
        app: market-app
        version: "${BUILD_NUMBER}"
    spec:
      containers:
      - name: market-app
        image: ${IMAGE_TAG}
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
        livenessProbe:
          httpGet:
            path: /inventory/health
            port: 8070
          initialDelaySeconds: 20
          periodSeconds: 10
EOF

# Update Service
log "5. Updating service..."
kubectl apply -n ${NAMESPACE} -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: market-service
spec:
  selector:
    app: market-app
    version: "${BUILD_NUMBER}"
  ports:
  - port: 8070
    targetPort: 8070
EOF

# Wait for rollout
log "6. Waiting for rollout..."
kubectl rollout status deployment/${DEPLOYMENT_NAME} -n ${NAMESPACE} --timeout=120s

log "✅ Deployment ${DEPLOYMENT_NAME} completed successfully!"

# Show status
echo ""
echo "📊 Status:"
kubectl get pods,svc,deploy -n ${NAMESPACE} -l app=market-app