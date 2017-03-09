properties([
    [$class: 'BuildDiscarderProperty', strategy: [$class: 'LogRotator', numToKeepStr: '20']],
    [$class: 'GithubProjectProperty', projectUrlStr: 'https://github.com/coreos/matchbox'],
    [$class: 'PipelineTriggersJobProperty', triggers: [
      [$class: 'GitHubPushTrigger'],
    ]]
])
parallel (
  etcd3: {
    node('fedora && bare-metal') {
      stage('etcd3') {
        timeout(time:3, unit:'MINUTES') {
          checkout scm
          sh '''#!/bin/bash -e
          cat /etc/os-release
          export ASSETS_DIR=~/assets; ./tests/smoke/etcd3
          '''
        }
      }
    }
  },
  k8s: {
    node('fedora && bare-metal') {
      stage('k8s') {
        timeout(time:8, unit:'MINUTES') {
          checkout scm          
          sh '''#!/bin/bash -e
          cat /etc/os-release
          export ASSETS_DIR=~/assets; ./tests/smoke/k8s
          '''
        }
      }
    }
  },
  bootkube: {
    node('fedora && bare-metal') {
      stage('bootkube') {
        timeout(time:10, unit:'MINUTES') {
          checkout scm          
          sh '''#!/bin/bash -e
          cat /etc/os-release
          chmod 600 ./tests/smoke/fake_rsa
          export ASSETS_DIR=~/assets; ./tests/smoke/bootkube
          '''
        }
      }
    }
  }
)
