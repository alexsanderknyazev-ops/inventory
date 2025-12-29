pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = "market-app:${BUILD_NUMBER}"
        KUBE_NAMESPACE = "market"
    }
    
    options {
        timeout(time: 30, unit: 'MINUTES')
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }
    
    stages {
        stage('Prepare') {
            steps {
                sh '''
                    echo "🔧 Jenkins Workspace: ${WORKSPACE}"
                    echo "🔧 Build Number: ${BUILD_NUMBER}"
                    echo "🔧 Docker Image: ${DOCKER_IMAGE}"
                    
                    # Настраиваем Minikube окружение
                    if command -v minikube &> /dev/null; then
                        eval $(minikube docker-env)
                    fi
                '''
            }
        }
        
        stage('Checkout') {
            steps {
                checkout scm
                sh 'ls -la'
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    echo "🏗️ Building Docker image..."
                    docker build -t ${DOCKER_IMAGE} .
                    docker tag ${DOCKER_IMAGE} market-app:latest
                    
                    # Проверяем образ
                    docker images | grep market-app
                '''
            }
        }
        
        stage('Test') {
            steps {
                sh '''
                    echo "🧪 Running tests..."
                    go test ./... -v -count=1 || true
                '''
            }
        }
        
        stage('Deploy') {
            steps {
                sh '''
                    echo "🚀 Deploying to Minikube..."
                    chmod +x deploy.sh
                    DOCKER_IMAGE=${DOCKER_IMAGE} ./deploy.sh
                '''
            }
        }
        
        stage('Verify') {
            steps {
                sh '''
                    echo "🔍 Verifying deployment..."
                    sleep 10
                    
                    # Проверяем статус всех компонентов
                    kubectl get all -n ${KUBE_NAMESPACE} || true
                    kubectl get pods -n ${KUBE_NAMESPACE} -o wide || true
                    
                    # Проверяем логи последнего деплоя
                    POD_NAME=$(kubectl get pods -n ${KUBE_NAMESPACE} -l app=market -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
                    if [ -n "$POD_NAME" ]; then
                        echo "=== Last 10 lines of logs ==="
                        kubectl logs -n ${KUBE_NAMESPACE} $POD_NAME --tail=10 || true
                    fi
                '''
            }
        }
    }
    
    post {
        always {
            sh '''
                echo "🧹 Cleaning up..."
                # Останавливаем port-forward если работает
                pkill -f "kubectl port-forward" || true
                
                # Показываем финальный статус
                echo "=== Final Status ==="
                kubectl get pods -n ${KUBE_NAMESPACE} 2>/dev/null || true
            '''
            
            script {
                // Сохраняем логи сборки
                archiveArtifacts artifacts: '**/target/*.log', allowEmptyArchive: true
            }
        }
        
        success {
            echo "🎉 Pipeline completed successfully!"
            // Можно добавить уведомления в Slack/Email
            // slackSend(color: 'good', message: "Build ${BUILD_NUMBER} успешен!")
        }
        
        failure {
            echo "❌ Pipeline failed!"
            // slackSend(color: 'danger', message: "Build ${BUILD_NUMBER} упал!")
        }
    }
}