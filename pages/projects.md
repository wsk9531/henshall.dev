---
title: Projects
description: Ben Henshall's Programming Projects
url: projects
---
# Projects
* [q3query](#quake-3-server-status-api)
* [SSG](#static-site-generator)
* [rmo-meta-util](#medical-metadata-pipeline)

## Quake 3 Server Status API
 *q3query* | **WIP**

Quake 3 servers use UDP for network communication. 
I'm writing a web service to make it easier to access a server's status (player counts, map, live score, and other game settings) via a REST API. 

This will be complemented with a simple web frontend, and a Discord bot to serve timely updates to a channel/on demand with a slash command.

There are plenty of opportunities to leverage Go's concurrency features in this project, and I'm excited to continue building with it here.

Looking forward to releasing this soon!

***

## Static Site Generator 
*Powering henshall.dev* | [Source on GitHub.](https://github.com/wsk9531/henshall.dev)

My little slice of the internet is powered by a custom Go application that turns hand-written markdown into static HTML and CSS.

It runs on Cloudflare Pages. I use GitHub Actions workflows to continuously integrate changes and deploy updates. `golangci-lint` is used in parallel to catch bugs and improve consistency. 

I largely stuck to a test-driven-development (TDD) approach for this project, enjoyed the iteration loop, and felt that my code was easy to work with throughout the process. 

### Good things: 
- I know how everything works.
- I use a handful of robust libraries like [goldmark](https://github.com/yuin/goldmark) and [go-yaml](https://gopkg.in/yaml.v3); rather than depend on a whole framework.
- Everything is really quick: content authoring, CI/CD pipelines, page load times, etc.
- Core SEO features are built in from the get-go.
- If I want to do something, I get to build it `=]`

### *However...*
- If I want to do something, I *have* to build it `=[`

***

## Medical Metadata Pipeline
*rmo-meta-util* 

Part of a year-long R&D co-op project with Deloitte and Auckland District Health Board to:

>Explore how innovative AI & cognitive services could be used to enhance clinician workflows for Resident Medical Officers (RMOs). 

RMOs use a quick-reference handbook application when meeting patients. Sessions with RMOs and ADHB staff established the handbook as an area ripe for improvement.

My focus was on *rmo-meta-util*, a stripped-back demonstration of a content pipeline. 

### What is the value?
- Almost every clinican interviewed complained about search. 
    
    - Search queries to in our demo app are now filterable and facetable, resulting in a direct improvement to a clinician's workflow.

    - We demonstrated how querying Kendra is a lot more flexible than exact-match string searches, which came across positively in a feedback session.

    - Metadata improves the relevance of Kendra's search results.

- RMO Handbook maintainers complained about a slow, yearly cycle to update content, supplementing this with paper printouts as information changed.

    - A frontend to our bucket would allow maintainers to get information out as it develops, while also keeping knowledge in one central place.

### How does it work?
1.  Uploading a page of the RMO handbook to an `S3` bucket triggers a Lambda function.
2.  The Lambda passes content through to `Comprehend Medical`, an NLP service for entity extraction.
3.  Entities are extracted based on the following categories: ```"ANATOMY", "MEDICAL_CONDITION", "MEDICATION", "TEST_TREATMENT_PROCEDURE".```
4. A number of entities over a relevance threshold are structured into custom attributes for use in `Kendra`.


| [![Digital Human Project Poster](/static/img/projects/digital-human.jpg)](/static/img/projects/digital-human.jpg) |
|:---:| 
| *Winner of the 2021 Fisher and Paykel Healthcare Excellence Award for Best R&D Project Poster.*  |
