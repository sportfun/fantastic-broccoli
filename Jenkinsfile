pipeline {
  agent any
  stages {
    stage('Build gakisitor') {
      agent {
        docker {
          image 'golang:alpine3.8'
        }

      }
      steps {
        sh 'apk add --no-cache git'
        dir(path: '/opt') {
          git(url: 'https://github.com/sportfun/gakisitor', branch: 'master', changelog: true)
          sh '''cd gakisitor
go build -o gakisitor .'''
        }

      }
    }
  }
}