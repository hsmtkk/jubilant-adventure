// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import { Construct } from "constructs";
import { App, TerraformStack, CloudBackend, NamedCloudWorkspace } from "cdktf";
import * as google from '@cdktf/provider-google';

const project = 'jubilant-adventure';
const region = 'us-central1';
const repository = 'jubilant-adventure';

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new google.provider.GoogleProvider(this, 'google', {
      project,
      region,
    });

    new google.artifactRegistryRepository.ArtifactRegistryRepository(this, 'registry', {
      format: 'docker',
      location: region,
      repositoryId: 'registry',
    });

    new google.cloudbuildTrigger.CloudbuildTrigger(this, 'buildTrigger', {
      filename: 'cloudbuild.yaml',
      github: {
        owner: 'hsmtkk',
        name: repository,
        push: {
          branch: 'main',
        },
      },
    });

    const sumServiceRunner = new google.serviceAccount.ServiceAccount(this, 'sumServiceRunner', {
      accountId: 'sumServiceRunner',
    });

    new google.cloudRunV2Service.CloudRunV2Service(this, 'sumService', {
      name: 'sumService',
      template: {
        containers: [{
          image: 'us-docker.pkg.dev/cloudrun/container/hello',
        }],
        scaling: {
          minInstanceCount: 0,
          maxInstanceCount: 1,
        },
        serviceAccount: sumServiceRunner.email,
      },
    });

  }
}

const app = new App();
const stack = new MyStack(app, "jubilant-adventure");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("jubilant-adventure")
});
app.synth();
