pipeline {
    agent any

    environment {
        imageName = "tyrm/go-activitypub-relay"
        image = ''
    }

    stages {
        stage("build") {
            steps {
                script {
                    image = docker.build(imageName, ".")
                }
            }
        }
        stage("push") {
            steps {
                script {
                    withDockerRegistry([credentialsId: 'docker-hub-credentials', url: '']) {
                        image.push("latest")
                    }
                }
            }
        }
    }
}