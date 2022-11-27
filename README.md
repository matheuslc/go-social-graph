<p align="center">
  <img src="https://user-images.githubusercontent.com/4161171/204121270-37e488f1-9a51-4442-a007-de13fbb91f0a.png" width="250px">
</p>

[![Coverage Status](https://coveralls.io/repos/github/matheuslc/go-social-graph/badge.svg?branch=main)](https://coveralls.io/github/matheuslc/go-social-graph?branch=main) ![Vercel](https://vercelbadge.vercel.app/api/matheuslc/go-social-graph)

# Social Graph
Social graph is a dead simple micro blogging platform. It remembers the old and good medium platform, combined with some fast interactions twitter has.

Check the **[ADRs](https://github.com/matheuslc/go-social-graph/tree/main/adrs)** folder to better understand some decisions taken on this project. We track new features and changes through issues. Take a look at the [CONTRIBUTING.md](https://github.com/matheuslc/go-social-graph/blob/main/CONTRIBUTING.md) file to better understand how to contribute.

## Stack
* Golang 1.18
* Neo4j

## Interface
[We also have an interface!](https://github.com/matheuslc/social-graph)


## Deploy
Right now we are deploying on Vercel using Serverless function. It's free! [Take a look](https://vercel.com/docs/concepts/functions/serverless-functions/supported-languages#go).

It's all automated, so after your PR got merged to the main branch, we automatically ship your code to production. During the development, we also generate a preview version of your code. 

## Having fun
Although we have a production version running, this is just a toy project that I use to stress some software engineer concepts. Everything can change.

That said, feel free to purpose architectural improvements and new ideas. 
