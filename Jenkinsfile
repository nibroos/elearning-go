pipeline {
  agent any
  environment {
    GIT_REPO = 'git@github.com:nibroos/elearning-go.git'
    SSH_CREDENTIALS_ID = '2dcfa3e4-fa4d-4702-a362-4ace13f87646'
    VPS_USER = credentials('vps-user-27')
    VPS_HOST = credentials('vps-host-27')
    VPS_DEPLOY_DIR = credentials('vps-deploy-dir-elearningbe-27')

    POSTGRES_USER = credentials('vps-postgres-username-elearningbe-27')
    POSTGRES_PASSWORD = credentials('vps-postgres-password-elearningbe-27')
    POSTGRES_DB = credentials('vps-postgres-elearningbe-27')
    POSTGRES_PORT = credentials('vps-postgres-port-elearningbe-27')
    POSTGRES_HOST = credentials('vps-postgres-host-elearningbe-27')

    GATEWAY_PORT = credentials('vps-gateway-elearningbe-27')
    SERVICE_GRPC_PORT = credentials('vps-service-grpc-elearningbe-27')
    SERVICE_REST_PORT = credentials('vps-service-rest-elearningbe-27')
    MASTER_SERVICE_GRPC_PORT = credentials('vps-master-service-grpc-elearningbe-27')
    MASTER_SERVICE_REST_PORT = credentials('vps-master-service-rest-elearningbe-27')
    ACTIVITIES_SERVICE_GRPC_PORT = credentials('vps-activities-service-grpc-elearningbe-27')
    ACTIVITIES_SERVICE_REST_PORT = credentials('vps-activities-service-rest-elearningbe-27')

    JWT_SECRET = credentials('vps-jwt-secret-elearning-27')
  }

  stages {
    stage('Clone Repository on VPS') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              # Add known hosts for GitHub
              ssh-keyscan -H github.com >> ~/.ssh/known_hosts
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'ssh-keyscan -H github.com >> ~/.ssh/known_hosts'
              
              # Test SSH connection first
              echo "Testing SSH connection..."
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'source ~/.bashrc; echo "SSH connection successful!"'
              
              # Clone the repository
              echo "Cloning repository..."
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'git clone ${GIT_REPO} /var/www/e-learning'
            """
          }
        }
      }
    }

    stage('Build Docker Images on VPS') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                cd ${VPS_DEPLOY_DIR} &&
                docker-compose -f docker/docker-compose.yml build > build_output.log 2>&1
              '
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'cat ${VPS_DEPLOY_DIR}/build_output.log'
            """
          }
        }
      }
    }

    stage('Deploy') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                cd ${VPS_DEPLOY_DIR} &&
                docker-compose -f docker/docker-compose.yml up -d > deploy_output.log 2>&1
              '
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'cat ${VPS_DEPLOY_DIR}/deploy_output.log'
            """
          }
        }
      }
    }

    stage('Run Migrations') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                docker exec -it service-prod-learninggo /usr/local/bin/migrate -path /app/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable up > migrate_output.log 2>&1
              '
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'cat migrate_output.log'
            """
          }
        }
      }
    }

    stage('Check Logs') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                docker logs service-prod-learninggo > service_logs.log 2>&1
              '
              ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'cat service_logs.log'
            """
          }
        }
      }
    }
  }

  post {
    failure {
      script {
        echo 'Build failed. Keeping the previous build up and running.'
      }
    }
  }
}