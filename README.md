# henshall.dev

Source code and content for my static website generator. 

See the end result at [henshall.dev](https://henshall.dev)

## Usage
This project is designed to work within a CI/CD pipeline. 
1. Clone the repo and checkout a new branch
2. Modify `.github/workflows/publish.yml` / GitHub Secrets with your Cloudflare Pages `apiToken, accountId, projectName`
3. Overwrite `/pages/` with your own site content
4. Modify templates, if required
5. Make a PR to run code audit, push to `master` to deploy.


### Local Usage
WIP, but run the following to build the executable, generate pages, and serve files locally.
#### Build
    $ go build -v -o bin/ssg cmd/cli/main.go
#### Run
    $ bin/ssg generate
    $ bin/ssg serve

Run `bin/ssg COMMAND -help` for available flags.


## Content Guide

A web page consists of frontmatter, written in YAML, and content, written in markdown. 

Frontmatter MUST open and close with a seperator `---` on a new line. 

| Fields        | Page          | Blog           |
| ------------- | ------------- | -------------  |
| title         | ✔️            | ✔️            |
| description   | ✔️            | ✔️            |
| url           | ✔️            | ✔️            |
| tags          | ❌            | ✔️            |
| published     | ❌            | ✔️            |
| updated       | ❌            | ✔️            |

**MANDATORY FIELDS:**
- title
- description
- url

All other fields are optional.

--- 

### Page Example

```yaml
---
title: A Lovely Web Page
description: This is my description. It populates our meta description tag!
slug: lovely-web-page
---

Content Begins here!
```

### Blog Example

```yaml
---
title: A Lovely Blog Post
description: This is my description. It populates our meta description tag!
slug: lovely-blog
tags: 
  - "Example"
  - "Another one"
published: 2022-12-31
updated: 2023-01-01
---

Content Begins here!
```
The YAML package I'm using can parse the following date formats:

    canonical: 2001-12-15T02:59:43.1Z
    iso8601: 2001-12-14t21:59:43.10-05:00
    date: 2002-12-14
