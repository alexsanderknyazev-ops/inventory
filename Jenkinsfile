pipeline {
    agent any
    
    stages {
        stage('Checkout') {
            steps { checkout scm }
        }
        
        stage('Deploy') {
            steps {
                sh '''
                    echo "🚀 Jenkins Deployment"
                    chmod +x *.sh
                    ./final-postgres.sh || true
                    ./deploy-app-only.sh || true
                    echo "✅ Если деплой не сработал, запустите вручную на машине с Minikube"
                '''
            }
        }
    }
}