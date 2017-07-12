pipeline {
  agent none

  options {
    timeout(time:45, unit:'MINUTES')
    buildDiscarder(logRotator(numToKeepStr:'20'))
  }

  stages {
    stage('Cluster Tests') {
      steps {
        parallel (
          etcd3: {
            node('fedora && bare-metal') {
              timeout(time:5, unit:'MINUTES') {
                checkout scm
                sh '''#!/bin/bash -e
                export ASSETS_DIR=~/assets; ./tests/smoke/etcd3
                '''
                deleteDir()
              }
            }
          },
          bootkube: {
            node('fedora && bare-metal') {
              timeout(time:60, unit:'MINUTES') {
                checkout scm          
                sh '''#!/bin/bash -e
                chmod 600 ./tests/smoke/fake_rsa
                export ASSETS_DIR=~/assets; ./tests/smoke/bootkube
                '''
                deleteDir()
              }
            }
          },
          "etcd3-terraform": {
            node('fedora && bare-metal') {
              timeout(time:10, unit:'MINUTES') {
                checkout scm
                sh '''#!/bin/bash -e
                export ASSETS_DIR=~/assets; export CONFIG_DIR=~/matchbox/examples/etc/matchbox; ./tests/smoke/etcd3-terraform
                '''
                deleteDir()
              }
            }
          },
          "bootkube-terraform": {
            node('fedora && bare-metal') {
              timeout(time:60, unit:'MINUTES') {
                checkout scm          
                sh '''#!/bin/bash -e
                chmod 600 ./tests/smoke/fake_rsa
                export ASSETS_DIR=~/assets; export CONFIG_DIR=~/matchbox/examples/etc/matchbox; ./tests/smoke/bootkube-terraform
                '''
                deleteDir()
              }
            }
          },
        )
      }      
    }
  }
}
