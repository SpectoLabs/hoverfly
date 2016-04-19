# Motivation

Building and testing interdependent applications and services is difficult. Whether you're building a mobile application that needs to talk to an legacy API that you don't have reliable access to, or whether you're building a microservice that relies on two other services that haven't been built yet - the problem is the same. How do you develop and test against external dependencies you can't control?

You can use mocking libraries as substitutes for external dependencies - but mocking libraries are intrusive. You have to modify your application code to use them. You can also use stubs - but for this you would need to configure your application to use different endpoints depending on the environment.

Then there is the problem of managing test data. Often, to write proper tests, you need fine-grained control over the data in your mocks or stubs. Managing test data across large projects and teams can become very complex, often introducing bottlenecks that can impact delivery times.

Integration testing "over the wire" is problematic too. When stubs and mocks are substituted for real services (in a continuous integration environment, for example) new variables are introduced, like network latency and random failures.

Hoverfly is a light-weight, flexible, easy-to-use tool that can be used in local development and in testing environments to "simulate" external dependencies. It is un-intrusive, resource-efficient and can be easily extended to manipulate data and simulate behaviours such as network failure, latency or rate-limits.

  
