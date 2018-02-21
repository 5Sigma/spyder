package docgen

var IndexTemplate = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="theme-color" content="#000000">
		<link type="text/css" rel="stylesheet" href="doc.css" />
    <title>{{ .ProjectName }}</title>
  </head>
  <body>
    <div id="root"></div>
		<script>
			var apiData = {
					"projectName": "{{ .ProjectName }}",
					"baseUrl": "{{ .BaseURL}}",
					"endpoints": [
						{{ range $i, $ep := .Endpoints }}
							{{if $i}},{{end}}{
								"url": "{{ $ep.BaseURL }}",
								"name": "{{ $ep.Name }}",
								"path": "{{ $ep.Path }}",
								"method": "{{ $ep.Method}}",
								"params": [],
								"response": {{ $ep.ResponseExample}},
								"request": {{ $ep.RequestExample}},
								"documentation": "{{ $ep.Documentation }}",
								"headers": [
									{{ range $hi, $h := .Headers }}
										{{if $hi}},{{end}}
										{
											"key": "{{ $h.Key }}",
											"description": "{{ $h.Description }}",
											"example": "{{ $h.Example }}"
										}
									{{ end }}
								]
							}
						{{ end }}
					]
				}
			</script>
		<script src="doc.js"></script>
  </body>
</html>
`
