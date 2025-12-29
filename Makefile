.PHONY: build test docker deploy clean

# Переменные
IMAGE_NAME = market-app
TAG = latest
NAMESPACE = market

# Сборка приложения
build:
	go build -o market-app main.go

# Тесты
test:
	go test ./... -v

# Сборка Docker образа
docker:
	docker build -t $(IMAGE_NAME):$(TAG) .

# Запуск локально
run:
	go run main.go

# Развертывание в Minikube
deploy:
	./deploy.sh

# Запуск полного цикла CI
ci: test docker deploy

# Очистка
clean:
	kubectl delete -n $(NAMESPACE) --all all 2>/dev/null || true
	docker rmi $(IMAGE_NAME):$(TAG) 2>/dev/null || true

# Просмотр логов
logs:
	@POD=$$(kubectl get pods -n $(NAMESPACE) -l app=market -o jsonpath='{.items[0].metadata.name}' 2>/dev/null); \
	if [ -n "$$POD" ]; then \
		kubectl logs -n $(NAMESPACE) -f $$POD; \
	else \
		echo "No pods found"; \
	fi

# Port-forward для доступа
port-forward:
	kubectl port-forward -n $(NAMESPACE) svc/market-service 8070:8070

# Проверка состояния
status:
	kubectl get all -n $(NAMESPACE)