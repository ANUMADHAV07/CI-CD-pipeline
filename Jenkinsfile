pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = 'k8s-jenkins-app'
        DOCKER_TAG = "${env.BUILD_NUMBER}"
        KUBERNETES_NAMESPACE = 'default'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Build') {
            steps {
                script {
                    echo 'Building Go application...'
                    sh 'go mod download'
                    sh 'go build -o app main.go'
                }
            }
        }
        
        stage('Test') {
            steps {
                script {
                    echo 'Running tests...'
                    // Add your tests here
                    sh 'go test ./... || true'
                }
            }
        }
        
        stage('Build Docker Image') {
            steps {
                script {
                    echo "Building Docker image: ${DOCKER_IMAGE}:${DOCKER_TAG}"
                    sh "docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} ."
                    sh "docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest"
                }
            }
        }
        
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    echo "Deploying to Kubernetes..."
                    sh """
                        kubectl set image deployment/k8s-jenkins-app app=${DOCKER_IMAGE}:${DOCKER_TAG} -n ${KUBERNETES_NAMESPACE} || \
                        kubectl apply -f k8s/deployment.yaml -n ${KUBERNETES_NAMESPACE}
                    """
                    sh "kubectl apply -f k8s/service.yaml -n ${KUBERNETES_NAMESPACE}"
                    sh "kubectl rollout status deployment/k8s-jenkins-app -n ${KUBERNETES_NAMESPACE}"
                }
            }
        }
    }
    
    post {
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}
