pipeline {
  agent any
  environment {
    GIT_REPO = 'git@github.com:nibroos/elearning-go.git'
    SSH_CREDENTIALS_ID = 'vps-ssh-credentials-elearning-27'
    VPS_USER = credentials('vps-user-27')
    VPS_HOST = credentials('vps-host-27')
    VPS_DEPLOY_DIR = credentials('vps-deploy-dir-elearningbe-27')

    POSTGRES_USER = credentials('vps-postgres-username-elearningbe-27')
    POSTGRES_PASSWORD = credentials('vps-postgres-password-elearningbe-27')
    POSTGRES_DB = credentials('vps-postgres-elearningbe-27')
    POSTGRES_PORT = credentials('vps-postgres-port-elearningbe-27')
    POSTGRES_HOST = credentials('vps-postgres-host-elearningbe-27')

    POSTGRES_USER_TEST = credentials('vps-postgres-username-test-elearningbe-27')
    POSTGRES_PASSWORD_TEST = credentials('vps-postgres-password-test-elearningbe-27')
    POSTGRES_DB_TEST = credentials('vps-postgres-test-elearningbe-27')
    POSTGRES_PORT_TEST = credentials('vps-postgres-port-test-elearningbe-27')
    POSTGRES_HOST_TEST = credentials('vps-postgres-host-test-elearningbe-27')

    GATEWAY_PORT = credentials('vps-gateway-elearningbe-27')
    SERVICE_GRPC_PORT = credentials('vps-service-grpc-elearningbe-27')
    SERVICE_REST_PORT = credentials('vps-service-rest-elearningbe-27')
    MASTER_SERVICE_GRPC_PORT = credentials('vps-master-service-grpc-elearningbe-27')
    MASTER_SERVICE_REST_PORT = credentials('vps-master-service-rest-elearningbe-27')
    ACTIVITIES_SERVICE_GRPC_PORT = credentials('vps-activities-service-grpc-elearningbe-27')
    ACTIVITIES_SERVICE_REST_PORT = credentials('vps-activities-service-rest-elearningbe-27')

    REDIS_HOST = credentials('vps-redis-host-elearningbe-27')
    REDIS_PORT = credentials('vps-redis-port-elearningbe-27')
    REDIS_PASSWORD = credentials('vps-redis-password-elearningbe-27')
    REDIS_DB = credentials('vps-redis-db-elearningbe-27')

    REDIS_HOST_TEST = credentials('vps-redis-host-test-elearningbe-27')
    REDIS_PORT_TEST = credentials('vps-redis-port-test-elearningbe-27')
    REDIS_PASSWORD_TEST = credentials('vps-redis-password-test-elearningbe-27')
    REDIS_DB_TEST = credentials('vps-redis-db-test-elearningbe-27')

    JWT_SECRET = credentials('vps-jwt-secret-elearning-27')
    BUILD_DIR = "build-${env.BUILD_ID}"
  }

  stages {
    stage('Clone Repository on VPS') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              # Add known hosts for GitHub
              ssh-keyscan -H github.com >> ~/.ssh/known_hosts
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'ssh-keyscan -H github.com >> ~/.ssh/known_hosts'
              
              # Test SSH connection first
              echo "Testing SSH connection..."
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'source ~/.bashrc; echo "SSH connection successful!"'
              
              # Clone the repository
              echo "Cloning repository..."
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'rm -rf ${VPS_DEPLOY_DIR} &&
              git clone -b rebase-test ${GIT_REPO} ${VPS_DEPLOY_DIR}'
            """
          }
        }
      }
    }

    stage('Create .env File') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                echo "POSTGRES_USER=${POSTGRES_USER}" > ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_DB=${POSTGRES_DB}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_DB_TEST=${POSTGRES_DB_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_PORT=${POSTGRES_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_HOST=${POSTGRES_HOST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_HOST_TEST=${POSTGRES_HOST_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_USER_TEST=${POSTGRES_USER_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_PASSWORD_TEST=${POSTGRES_PASSWORD_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "POSTGRES_PORT_TEST=${POSTGRES_PORT_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "GATEWAY_PORT=${GATEWAY_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "SERVICE_GRPC_PORT=${SERVICE_GRPC_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "SERVICE_REST_PORT=${SERVICE_REST_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "MASTER_SERVICE_GRPC_PORT=${MASTER_SERVICE_GRPC_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "MASTER_SERVICE_REST_PORT=${MASTER_SERVICE_REST_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "ACTIVITIES_SERVICE_GRPC_PORT=${ACTIVITIES_SERVICE_GRPC_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "ACTIVITIES_SERVICE_REST_PORT=${ACTIVITIES_SERVICE_REST_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "JWT_SECRET=${JWT_SECRET}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_HOST=${REDIS_HOST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_PORT=${REDIS_PORT}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_PASSWORD=${REDIS_PASSWORD}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_DB=${REDIS_DB}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_HOST_TEST=${REDIS_HOST_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_PORT_TEST=${REDIS_PORT_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_PASSWORD_TEST=${REDIS_PASSWORD_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                echo "REDIS_DB_TEST=${REDIS_DB_TEST}" >> ${VPS_DEPLOY_DIR}/docker/.env &&
                cp ${VPS_DEPLOY_DIR}/docker/.env ${VPS_DEPLOY_DIR}/service/.env &&
                cp ${VPS_DEPLOY_DIR}/docker/.env ${VPS_DEPLOY_DIR}/gateway/.env
              '
            """
          }
        }
      }
    }

    stage('Running Tests') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                cd ${VPS_DEPLOY_DIR} &&
                docker compose -f docker/docker-compose-test.yml down --remove-orphans &&
                docker compose -f docker/docker-compose-test.yml up --build -d &&
                sleep 5 # Wait for containers to start
              '

              echo "Running tests.."

              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                docker exec service-test go test -v /app/internal/tests/... > test_output.log 2>&1 &&
                cat test_output.log
              '
              echo "Tests completed."
            """
          }
        }
      }
    }

    stage('Run Migrations on Test DB') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                cd ${VPS_DEPLOY_DIR}/service &&

                echo "Running test migrations on test database..." &&
                make migrate-test-up &&

                echo "Migrations completed."
                echo "Downing test migrations..."

                echo "Removing test containers..."
                docker compose -f ${VPS_DEPLOY_DIR}/docker/docker-compose-test.yml down --remove-orphans
              '
            """
          }
        }
      }
    }

    stage('Build & Deploy') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                cd ${VPS_DEPLOY_DIR} &&
                docker compose -f docker/docker-compose.yml down --remove-orphans &&
                docker compose -f docker/docker-compose.yml up --build -d > build_output.log 2>&1
              '
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} 'cat ${VPS_DEPLOY_DIR}/build_output.log'
            """
          }
        }
      }
    }

    stage('Run Migrations on Prod DB') {
      steps {
        script {
          sshagent(credentials: [SSH_CREDENTIALS_ID]) {
            sh """
              ssh -A -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_HOST} '
                docker exec \$(docker ps --filter "name=service" --format "{{.ID}}" | head -n 1) /usr/local/bin/migrate -path /apps/internal/database/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable up > migrate_output.log 2>&1 &&
                cat migrate_output.log
              '
            """
          }
        }
      }
    }
  }

  post {
    always {
      cleanWs()
    }

    failure {
      script {
        echo 'Build failed. Keeping the previous build up and running.'
      }
    }
  }
}