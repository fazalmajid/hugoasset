# Hugo Asset Management

This utility helps manage the static assets under the `static/` subdirectory
of your [Hugo][1] installation:
* Find broken links (missing assets)
* Find assets that are not linked to anywhere in the content so you can reclaim the disk space

It searches for `<a href="...">`, `<img src="..." srcset="...">` and
`<script src="...">`

Just run it under your hugo directory, just as you would `hugo` itself. The
parameter `-prefix` tells it what the URL prefix is for the site URLs. It
assumes static assets are under `static/` and output files under `public/`

 [1]: https://gohugo.io/
