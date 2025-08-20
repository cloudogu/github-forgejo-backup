@Library('github.com/cloudogu/gitops-build-lib@0.7.0')
import com.cloudogu.gitops.gitopsbuildlib.*
@Library('github.com/cloudogu/jenkins-deploy-helper-lib@main') _

node('sos') {
    properties([
            buildDiscarder(logRotator(numToKeepStr: '5')),
            disableConcurrentBuilds(),
            pipelineTriggers([cron('H H(2-4) * * *')])
    ])

    withCredentials([string(credentialsId: 'jenkins-pipeline-notifier-webhookurl', variable: 'WEBHOOK')]) {
        createTagAndDeploy(
                classname: 'github-forgejo-backup',
                webhook: "${WEBHOOK}",
                repositoryUrl: 'sos/gitops',
                filename: 'deployment.yaml',
                buildArgs: '',
                team: 'sos'
        )
    }
}
