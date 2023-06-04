!!!THIS REPO WAS MOVED TO github.com/evozon/future

Future is a library that facilitates the upgrade process of PHP projects.

LOGO HERE

# Setup
### Prerequisites
* Github or Gitlab CI
* composer

### Gitlab
The project is meant to be plugged into a Gitlab CI pipeline.

Steps:
1. <code>composer require --dev calinpristavu/future-proofing</code>
2. <code>cp vendor/calinpristavu/future-proofing/future-proofing.yaml .</code>
3. open `.future-proofing.yml` and change the `future-config:` section to your needs
4. plug the pipeline CI file and the stage into your project's CI file. Example:
    ```yaml
    # .gitlab-ci.yml
    include: '/future-proofing.yaml'
    
    stages:
      ...
      - future-proofing
    ```
5. run the pipeline and check the output of the `future-proofing` job

### Github
Soon(tm)

# What it does
The pipeline job will:
* upgrade the php version of the project in the docker container. See #Upgrading PHP(link here)
* upgrade composer dependencies to the latest versions
* run rectors against the codebase. See #Rector(link here)
* run the tests

# What it doesn't do
* have any purpose if the app has no tests
* write tests for you :)
* commit changes to the codebase

# Output examples
### Everything is prepared for the upgrade
(output from the pipeline here)
### This that should be covered before the upgrade
(output from the pipeline here)

# Upgrading PHP
Future is run in the pipeline, meaning it needs a container to run in. The container is defined in the `.future-proofing.yml` file. The default image is `php:fpm-alpine`. You can change it to any other container that has the php version you want to upgrade to. The container is used to run the tests, so make sure it has all the dependencies needed to run the tests.

Our recommendation is to use the same container as the one used in production, but change the php version to the one you want to upgrade to. This way you can be sure that the tests will run in the same environment as the production one.

Example of a container that uses php 7.4 but tries to upgrade to 8.2:
```yaml
# .gitlab-ci.yml
...
include: '/future-proofing.yaml'

stages:
  ...
  - test
  ...
  - future-proofing

phpunit:
  stage: test
  image: php:7.4-cli
  script:
    - composer install
    - vendor/bin/phpunit
```
```yaml
# .future-proofing.yml
.future-config:
   #replace this with the base setup of your project
   #if your project has a Dockerfile, use it and change the "FROM" section to php:fpm-alpine
   image: php:8.2-cli
```

This way the future proofing job will run all the upgrades and validation steps using PHP8.2 instead of PHP7.4.


# Rector
Big thanks to [Rector](https://github.com/rectorphp/rector) for providing the tool that does the heavy lifting. Future uses Rector to upgrade the codebase to the latest standards of coding.

The Rector rules are defined in the `.future-proofing.yml` file, in the future-proofing job.
Feel free to add/change/remove rectors as you see fit for the context of your application.

Rectors are manipulated using a binary provided by this library called `future`. You can run it locally to see what changes it will make to your codebase. Example:
```bash
$ vendor/bin/future add-rule \\Rector\\CodeQuality\\Rector\\Identical\\GetClassToInstanceOfRector
```
This command will add the `GetClassToInstanceOfRector` to the list of rectors in `rector.php` that will be run against the codebase.
However, we recommend that you configure rectors in the `.future-proofing.yml` file, since it will be configured specifically for the future proofing step in the CI pipeline.

# Recommendations
Future can be used to validate if you can upgrade *EVERYTHING* at once: 
* PHP version to the latest one
* composer dependencies to the latest versions
* the codebase to the latest standards. 

This is not recommended, as actually upgrading everything at once can lead to massive amounts of changes that are hard to review and test.

We recommend splitting the upgrade process into stages. First do the PHP version upgrades, one minor version at a time. For example if we are on PHP7.4, first use Future to upgrade to PHP8.0, then 8.1 and so on...

Then start working on dependencies and coding style(see #Rector(link here)). This way you can have smaller PRs that are easier to review and test.
