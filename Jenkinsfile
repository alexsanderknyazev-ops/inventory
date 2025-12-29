pipeline {
    agent any
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
                sh 'ls -la'
            }
        }
        
        stage('Run Deployment Script') {
            steps {
                sh '''
                    echo "🚀 Запускаю скрипт деплоя..."
                    echo "Доступные скрипты:"
                    ls -la *.sh
                    
                    # Делаем скрипты исполняемыми
                    chmod +x *.sh 2>/dev/null || true
                    
                    # Проверяем какой скрипт есть
                    if [ -f "deploy-app-only.sh" ]; then
                        echo "Запускаю: ./deploy-app-only.sh"
                        ./deploy-app-only.sh
                    elif [ -f "deploy.sh" ]; then
                        echo "Запускаю: ./deploy.sh"
                        ./deploy.sh
                    else
                        echo "❌ Скрипт деплоя не найден!"
                        echo "Создайте deploy-app-only.sh или deploy.sh"
                        exit 1
                    fi
                '''
            }
        }
        
        stage('Check Status') {
            steps {
                sh '''
                    echo "📊 Проверяю статус..."
                    sleep 5
                    echo "Деплой скрипт выполнен."
                    echo ""
                    echo "Для проверки статуса выполните на хосте с Minikube:"
                    echo "  kubectl get pods -n market"
                    echo "  kubectl logs -n market -l app=market-app"
                '''
            }
        }
    }
    
    post {
        success {
            echo "✅ Pipeline завершен успешно!"
            echo "Приложение должно быть развернуто в Minikube."
            echo ""
            echo "🌐 Для доступа к приложению:"
            echo "   kubectl port-forward -n market svc/market-service 8070:8070"
            echo "   Затем откройте: http://localhost:8070"
        }
        failure {
            echo "❌ Pipeline завершился с ошибкой!"
            echo ""
            echo "🔧 Возможные причины:"
            echo "   1. Minikube не запущен"
            echo "   2. Нет доступа к Docker из Jenkins"
            echo "   3. Ошибка в скрипте деплоя"
            echo ""
            echo "💡 Решение:"
            echo "   Запустите деплой вручную на машине с Minikube:"
            echo "   ./deploy-app-only.sh"
        }
    }
}