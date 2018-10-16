pipeline {
  agent any
  stages {
    stage('Build gakisitor - AMD64') {
      agent {
        docker {
          image 'golang:alpine3.8'
        }

      }
      steps {
        dir(path: '/opt/sportfun/go') {
          git(url: 'https://github.com/sportfun/gakisitor', branch: 'master', changelog: true)
        }

        sh 'go build'
      }
    }
  }
}