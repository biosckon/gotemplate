### gotemplate uses go text/template package to apply templates to data

Example in the example folder:

```Shell
$ gotemplate index_template.html data.csv index.html
```

take `index_template.html` parse it, take `data.csv` parse it, finally apply template to data and spit out the index.html.

IF: index_template.html IS:

```html
CTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{.title}}</title>
</head>
<body>
    {{.body}}
</body>
</html>

```

AND: data.csv IS:

```csv
title,Intersting page
body,Very boring content
something,some other text is here
```

THEN index.html will be:

```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Interesting page</title>
    </head>
    <body>
        Very boring content
    </body>
</html>
```

This is go text/template at it's simplest.
