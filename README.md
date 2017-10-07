# Hugo Asset Management

## Introduction

This utility helps manage the static assets under the `static/` subdirectory
of your [Hugo][1] installation:
* Find broken links (missing assets)
* Find assets that are not linked to anywhere in the content so you can reclaim the disk space

It searches for `<a href="...">`, `<img src="..." srcset="...">` and
`<script src="...">`

## Building from source

You need [Go](https://golang.org) as a prerequisite (but then you need it for
Hugo as well).

Git clone or extract the source code from a tarball, then run `make`

## Usage

Just run it under your hugo directory, just as you would `hugo` itself. The
parameter `-prefix` tells it what the URL prefix is for the site URLs. It
assumes static assets are under `static/` and output files under `public/`

 [1]: https://gohugo.io/
