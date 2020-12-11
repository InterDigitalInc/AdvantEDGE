# How to contribute to AdvantEDGE

_So you've decided that you would like to contribute to AdvantEDGE project... what next?<br>_

It's great to hear that you have interest in AdvantEDGE & we'd love to accept your bug fix or other contributions;
there is however few small guidelines that you need to follow before getting started.

_[... But I just have a question & I don't want to read this whole thing!!!](#if-you-have-questions)_

There are three main ways of contributing to the project: by reporting an issue, by suggesting an enhancement or by contributing content.<br> As such, we cover these use cases below
- [Reporting an issue (bug/documentation)](#reporting-issues)
- [Suggesting an enhancement](#suggesting-enhancements)
- [Contributing content](#contributing-content)
  - [Contributor License Agreement](#contributor-license-agreement)
  - [Your First Content Contribution](#your-first-content-contribution)
  - [Pull Requests](#pull-request)
  - [But... What can I contribute on?](#what-can-i-contribute-on)

## If You Have Questions
Please don't open a GitHub Issue to ask a question; you'll get faster result by using the resources below.

#### Resource #1 - AdvantEDGE Wiki
We put time & efforts keeping the wiki up to date, so we recommend to look there first.<br>
- Project questions -- [FAQ](https://github.com/InterDigitalInc/AdvantEDGE/wiki/faq) & [Roadmap](https://github.com/InterDigitalInc/AdvantEDGE/wiki/roadmap)
- Concepts questions -- [Platform concepts](https://github.com/InterDigitalInc/AdvantEDGE/wiki/platform-concepts), [platform APIs](https://github.com/InterDigitalInc/AdvantEDGE/wiki/API-Documentation), [Edge App. Models](https://github.com/InterDigitalInc/AdvantEDGE/wiki/edge-app-models), [Edge App. Types](https://github.com/InterDigitalInc/AdvantEDGE/wiki/edge-app-types) & [Frontend concepts](https://github.com/InterDigitalInc/AdvantEDGE/wiki/frontend-concepts)
- Services questions -- [platform APIs](https://github.com/InterDigitalInc/AdvantEDGE/wiki/API-Documentation), [Location service](https://github.com/InterDigitalInc/AdvantEDGE/wiki/location-service), [Application State Transfer service](https://github.com/InterDigitalInc/AdvantEDGE/wiki/state-transfer)
- Setup questions -- [Hardware](https://github.com/InterDigitalInc/AdvantEDGE/wiki/hw-configuration), [Runtime environment](https://github.com/InterDigitalInc/AdvantEDGE/wiki/runtime-environment), [Development environment](https://github.com/InterDigitalInc/AdvantEDGE/wiki/development-environment)
- Deployment questions -- [Deployment cheat-sheet](https://github.com/InterDigitalInc/AdvantEDGE/wiki/deployment-details), [meepctl CLI tool](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/meepctl/meepctl.md), [Build](https://github.com/InterDigitalInc/AdvantEDGE/wiki/build-advantedge), [Deploy](https://github.com/InterDigitalInc/AdvantEDGE/wiki/deploy-advantedge), [Upgrade](https://github.com/InterDigitalInc/AdvantEDGE/wiki/upgrade-advantedge), [Test](https://github.com/InterDigitalInc/AdvantEDGE/wiki/test-advantedge)
- Usage questions -- [GUI](https://github.com/InterDigitalInc/AdvantEDGE/wiki/gui-overview), [basic operation](https://github.com/InterDigitalInc/AdvantEDGE/wiki/basic-operation), [creating a first scenario](https://github.com/InterDigitalInc/AdvantEDGE/wiki/first-scenario), [scenario monitoring](https://github.com/InterDigitalInc/AdvantEDGE/wiki/scenario-monitoring), [using external nodes](https://github.com/InterDigitalInc/AdvantEDGE/wiki/external-nodes), [pod placement](https://github.com/InterDigitalInc/AdvantEDGE/wiki/pod-placement)

#### Resource #2 - Start a Discussion
Of course - because we cannot document everything - if you still have a question, you can reach out to us.

We use GitHub's new feature called [Discussions](https://github.com/InterDigitalInc/AdvantEDGE/discussions) - simply start a discussion and you will have direct access to the development team!

## Reporting issues
Whether it's a bug found while using the platform or simply a typo noticed while browsing the documentation, we appreciate that you open a GitHub issue ([here](https://github.com/InterDigitalInc/AdvantEDGE/issues))

#### Bug
When reporting a bug, try to be as concise & specific as possible so we can reproduce the problem.

A Bug Report template is provided in GitHub to help documenting the problem.
#### Documentation
Use the Custom Issue template to provide a link to the page, a copy the problematic text and an indication of what the problem is.

## Suggesting enhancements
We appreciate feature enhancements and as such, we collect feature ideas via GitHub issues ([here](https://github.com/InterDigitalInc/AdvantEDGE/issues))

If you are not sure about the proposed enhancement, it's always a good idea to communicate with us by starting a [Discussion](https://github.com/InterDigitalInc/AdvantEDGE/discussions) about the feature beforehand.

## Contributing content
Contributing content is more involving than submitting an issue or an enhancement request as it requires a CLA and learning how to operate with the project team.

The following sub-sections cover these aspects.

#### Contributor License Agreement
In order for us to accept content contributions, a Contributor License Agreement (CLA) is required.

- _[CLA for individual contributor](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/cla/interdigital-individual-cla-v1.pdf)_
- _[CLA for corporate contributors](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/cla/interdigital-corporate-cla-v1.pdf)_

Fill-in the requested information and send it to AdvantEDGE@InterDigital.com using subject `CLA`

> Note: it is important to include your GitHub id(s) in the CLA so we know who you are when submitting a pull-request.

#### Your First Content Contribution
In order to prepare your contribution, you will need to create your fork of the AdvantEDGE repo. To do so, simply click the Fork button in the [repo](https://github.com/InterDigitalInc/AdvantEDGE)

Your fork is your own copy of the repo where you can modify things without impacting others; your fork's GitHub path should be `https://github.com/<your-gh-username/AdvantEDGE.git`

After creating your fork, you can clone it locally and start making modifications. We also require making modifications on a branch originating from the `develop` branch (preferred) or alternatively the `master` branch.
```
git clone https://github.com/<your-gh-username/AdvantEDGE.git
git checkout develop
git checkout -b <your-branch-name>
```

During development, keep your commit message concise and precise so we understand what the change is about.

We accept a single contribution per branch, so if you are fixing two bugs - please create two different branches each originating from `develop` & each containing only the necessary changes for the bug they fix.

Finally, for code contributions, we require that you run the linter (e.g. `meepctl lint all`), [Unit Tests](https://github.com/InterDigitalInc/AdvantEDGE/wiki/Test-AdvantEDGE#run-unit-tests) and [Cypress tests](https://github.com/InterDigitalInc/AdvantEDGE/wiki/Test-AdvantEDGE#run-cypress-tests) on your branch.
> Note: Keep results to include in your pull request.

#### Pull Request
With implementation complete and linter/tests passing<br>
You are ready to make your pull request, [here](https://github.com/InterDigitalInc/AdvantEDGE/pulls).

First, make sure that all necessary code is committed to your branch and that you have pushed your branch back to your fork.
```
git add <your-modified-files>
git commit -m 'what-has-changed'
git push
```

If unsure, it's a good idea to double check that your branch has been pushed back to the repo using the GitHub browser client. Your fork can be accessed from [here](https://github.com/InterDigitalInc/AdvantEDGE/network/members), your branch should show up in the branch drop down of your fork and your changes should show up once the branch is selected.

In the pull request, indicate the branch containing your changes, the nature of the changes you made, a reference to the related GitHub issue and the test results.

_What to expect next?_<br>
We do peer reviews of all internal contributions - so as an external contributor you can expect that someone from the core team will review your PR.

This is normal and is part of the process. Be patient as we may be busy addressing other issues - we will eventually get back to you with a status & possibly change requests. Please note that we reserve the right to accept or not a contribution. For internal reasons, from time to time, we may decide not to include your proposition, this has nothing to do with your skills or the value of your contribution.

#### What can I contribute on
Making first contributions can be intimidating - after all it's often difficult to figure out where to start and learn how to interact with a new team (that's us :) ). This is normal and we are here to help.

We recommend the following approach:
1. Start small - documentation contributions are an easy way to start as they are simpler to perform. If you used the project before (or not) and find that documentation is confusing / missing - then you may have an opportunity. Alternatively, look for [issues](https://github.com/InterDigitalInc/AdvantEDGE/issues) tagged `documentation` & `good first issue`
2. Look for the [issues](https://github.com/InterDigitalInc/AdvantEDGE/issues) tagged `good first issue` - as these are simpler and often isolated / self-contained
3. Once you are accustomed with the project, you can start tackling other larger existing issues.
4. The holly-grail of contributing would be opening an issue / getting an enhancement request approved and submitting the PR implementing it

Anyhow, when picking an issue, it's always a good idea to comment on the issue to let others know you are looking at it; you can also email us at AdvantEDGE@InterDigital.com

We hope this is helpful - and...<br>
... looking forward to hear from you!
