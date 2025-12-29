pipeline {
  agent any

  environment {
    APP_NAME   = "my-hertz-app"
    IMAGE_NAME = "registry:5000/${APP_NAME}:${BUILD_NUMBER}"
    GO111MODULE = "on"
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    stage('Deps') {
      steps {
        sh 'go mod download'
      }
    }

    stage('Test') {
      steps {
        sh 'go test ./... -v'
      }
    }

    stage('Build Docker') {
      steps {
        sh "docker build -t ${IMAGE_NAME} ."
      }
    }

    stage('Push Docker') {
      steps {
        sh "docker push ${IMAGE_NAME}"
      }
    }

    stage('Run Integration Test') {
      steps {
        sh '''
        docker run -d --rm -p 8080:8080 --name hertz-test ${IMAGE_NAME}
        sleep 5
        curl -f http://localhost:8080/ || (docker logs hertz-test && exit 1)
        docker stop hertz-test
        '''
      }
    }
  }

  post {
    success {
      echo "Build and tests succeeded! Image pushed: ${IMAGE_NAME}"
    }
    failure {
      echo "Pipeline failed!"
    }
  }
}
