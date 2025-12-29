#!/bin/bash
set -e

NAMESPACE=${1:-market}
TIMEOUT=60
INTERVAL=5

echo "🩺 Starting health check for namespace: $NAMESPACE"

for i in $(seq 1 $((TIMEOUT/INTERVAL))); do
    # Проверяем готовность всех подов
    READY_PODS=$(kubectl get pods -n $NAMESPACE -l app=market \
        -o jsonpath='{.items[*].status.containerStatuses[?(@.ready==true)].ready}' | wc -w)
    
    TOTAL_PODS=$(kubectl get pods -n $NAMESPACE -l app=market \
        -o jsonpath='{.items[*].metadata.name}' | wc -w)
    
    if [ "$READY_PODS" -eq "$TOTAL_PODS" ] && [ "$TOTAL_PODS" -gt 0 ]; then
        echo "✅ All pods are ready!"
        exit 0
    fi
    
    echo "⏳ Waiting for pods... ($READY_PODS/$TOTAL_PODS ready)"
    sleep $INTERVAL
done

echo "❌ Health check timeout!"
exit 1