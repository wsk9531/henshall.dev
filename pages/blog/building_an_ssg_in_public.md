---
title: Building a static site generator in public
description: 
url: building-a-static-site-generator-in-public
tags:
  - "programming"
  - "ssg"
published: 2023-11-04
updated: 2023-11-13
---

# Building a static site generator in public

I've been interested in the [build in public](https://buildinpublic.xyz/) movement for a while now. 

Developers, creatives, and entrepreneurs who share as they go, transparently make decisions, and stand out through public accountability is a level of authenticity that I can't help but resonate with.

I thought that even though I'm not trying to turn my website generator into a product, or even grow the userbase above 1, it could be useful to share my list of desired features publicly.

## Make it work, make it right, make it fast
Getting *anything* live was my biggest priority. You can visit my site on a real domain and click through pages I spent far too much time trying to style. This is the MVP, and I'm stoked to get it out there.

Now that things work, it is time to make things right: modularizing the generation process and refactoring out those shortcuts =]

## What I'm working on

- [ ] Code
  - [ ] Essential refactoring
    - [ ] Fix up error handling
    - [ ] Modularize and parallelize `cmd/cli` code (extremely messy right now)
  - [ ] `serve` and `generate` cmd improvements
    - [ ] Copying `static/img` and `static/css` folders in an automated + OS-agnostic way
    - [ ] Live updates & regeneration
  - [ ] New generators
    - [ ] Add CSS generator that appends `goldmark-highlighting` CSS to my stylesheet 
    instead of inlining with generated HTML.
    - [ ] Add RSS feed generator
    - [ ] Add sitemap.xml generator
- [ ] Content
  - [x] Blogs
    - [x] First post
    - [ ] Learning Journey for Go and TDD
  - [ ] Add a 404 page
  - [ ] Add `robots.txt`
- [ ] Styling
  - [x] Style `<blockquote>` elements
  - [ ] Revise questionable (ugly) templates / styles e.g. blogs show publish date before h1
  - [ ] Update links on the `#narrow-nav-content` view (mobile) to have 48x48px gaps (suggested by [Lighthouse audit](https://developer.chrome.com/docs/lighthouse/overview/))


