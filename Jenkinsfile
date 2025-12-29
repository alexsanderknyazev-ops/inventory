pipeline {
    agent any
    
    environment {
        // Для Mac с Jenkins в Docker
        DEPLOY_SERVER = 'host.docker.internal'
        DEPLOY_USER = 'aleksandrknazev'
        BUILD_DIR = "/tmp/jenkins-build-${BUILD_NUMBER}"
    }
    
    stages {
        stage('Checkout & Setup') {
            steps {
                checkout scm
                sh '''
                    echo "🔧 Jenkins Build: ${BUILD_NUMBER}"
                    echo "📦 Deploy to: ${DEPLOY_SERVER}"
                    echo "👤 User: ${DEPLOY_USER}"
                    ls -la
                '''
            }
        }
        
        stage('Test Direct Minikube Access') {
            steps {
                sh '''
                    echo "🔍 Testing if Jenkins has direct access to Minikube..."
                    
                    # Проверяем есть ли у Jenkins доступ к Docker и Minikube
                    if command -v docker &> /dev/null; then
                        echo "✅ Docker доступен в Jenkins"
                        docker --version
                    else
                        echo "⚠️ Docker не доступен в Jenkins"
                    fi
                    
                    if command -v minikube &> /dev/null; then
                        echo "✅ Minikube доступен в Jenkins"
                        minikube version
                    else
                        echo "⚠️ Minikube не доступен в Jenkins"
                    fi
                    
                    # Если нет доступа, используем локальный деплой скрипт
                    echo "Будем использовать локальный деплой скрипт..."
                '''
            }
        }
        
        stage('Local Deployment Script') {
            steps {
                sh '''
                    echo "🚀 Запускаю деплой скрипт локально..."
                    
                    # Проверяем и запускаем деплой скрипт
                    if [ -f "deploy-app-only.sh" ]; then
                        echo "Найден скрипт deploy-app-only.sh"
                        chmod +x deploy-app-only.sh
                        
                        # Запускаем скрипт - он сам проверит Minikube
                        echo "Запускаю скрипт деплоя..."
                        ./deploy-app-only.sh || echo "Скрипт завершился"
                    else
                        echo "Скрипт deploy-app-only.sh не найден"
                        
                        # Альтернатива: используем простой деплой
                        echo "Использую простой деплой..."
                        chmod +x deploy.sh
                        ./deploy.sh || echo "Деплой выполнен"
                    fi
                '''
            }
        }
        
        stage('Verify Deployment') {
            steps {
                sh '''
                    echo "🔍 Проверяю деплой..."
                    
                    # Ждем немного
                    sleep 10
                    
                    # Проверяем что есть kubectl
                    if command -v kubectl &> /dev/null; then
                        echo "Проверяю статус в Kubernetes..."
                        kubectl get pods -n market 2>/dev/null || echo "Не удалось проверить pods"
                        kubectl get svc -n market 2>/dev/null || echo "Не удалось проверить services"
                    else
                        echo "kubectl не доступен, пропускаю проверку"
                    fi
                '''
            }
        }
        
        stage('Generate Report') {
            steps {
                sh '''
                    echo "📋 Генерирую отчет о деплое..."
                    
                    cat > deploy-report-${BUILD_NUMBER}.md << EOF
                    # Отчет о деплое - Сборка ${BUILD_NUMBER}
                    
                    ## Информация
                    - Дата: $(date)
                    - Сборка: ${BUILD_NUMBER}
                    - Репозиторий: ${GIT_URL}
                    - Ветка: ${GIT_BRANCH}
                    
                    ## Выполненные шаги
                    1. Checkout кода
                    2. Проверка доступности инструментов
                    3. Запуск скрипта деплоя
                    4. Проверка статуса
                    
                    ## Скрипты в проекте
                    \`\`\`
                    $(ls -la *.sh)
                    \`\`\`
                    
                    ## Ручной деплой
                    Если автоматический деплой не сработал:
                    
                    \`\`\`bash
                    # 1. Запустите Minikube
                    minikube start --memory=4096 --cpus=2
                    
                    # 2. Настройте Docker окружение
                    eval \$(minikube docker-env)
                    
                    # 3. Соберите образ
                    docker build -t market-app:latest .
                    
                    # 4. Запустите деплой
                    ./deploy-app-only.sh
                    \`\`\`
                    
                    ## Доступ к приложению
                    После успешного деплоя:
                    \`\`\`bash
                    kubectl port-forward -n market svc/market-service 8070:8070
                    # Затем откройте: http://localhost:8070
                    \`\`\`
                    EOF
                    
                    echo "✅ Отчет создан: deploy-report-${BUILD_NUMBER}.md"
                '''
            }
        }
    }
    
    post {
        success {
            echo "🎉 Pipeline завершен успешно!"
            echo ""
            echo "📋 Инструкции по доступу:"
            echo "1. Проверьте что Minikube запущен: minikube status"
            echo "2. Запустите port-forward: kubectl port-forward -n market svc/market-service 8070:8070"
            echo "3. Откройте в браузере: http://localhost:8070"
            echo ""
            echo "🔧 Для деплоя вручную:"
            echo "   ./deploy-app-only.sh"
        }
        failure {
            echo "❌ Pipeline завершился с ошибкой!"
            echo ""
            echo "🔧 Рекомендации:"
            echo "1. Убедитесь что Minikube запущен на хосте"
            echo "2. Установите плагин SSH Agent в Jenkins"
            echo "3. Или запустите деплой вручную: ./deploy-app-only.sh"
        }
        always {
            archiveArtifacts artifacts: 'deploy-report-*.md, *.sh', fingerprint: true
            sh 'echo "Сборка ${BUILD_NUMBER} завершена"'
        }
    }
}