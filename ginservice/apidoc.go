package ginservice

const apiDoc = `<!DOCTYPE html>
<html>
<head lang="en">
  <meta charset="UTF-8">
  <title> API Doc </title>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
  <link href="http://fonts.googleapis.com/css?family=Roboto" rel="stylesheet" type="text/css">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.0/jquery.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
  <style type="text/css">
    body {
      font-family: 'Roboto', sans-serif;
    }
    strong {
      padding: 6px; margin: 6px;
    }
    pre {
      padding: 6px; margin: 6px;
      white-space: pre-wrap;
    }
    .prettyprint {
      border: 1px solid #ccc;
      margin-bottom: 0;
      padding: 9px;
    }
    .method {
      color: white;
      padding: 6px 15px;
      margin-right: 2px;
      border-radius: 3px;
      min-width: 80px;
      font-weight: 700;
    }
    .GET {
      background: #61affe;
    }
    .POST {
      background: #49cc90;
    }
    .PUT {
      background: #fca130;
    }
    .DELETE {
      background: #f93e3e;
    }
    .menu {
      max-height: 800px;
      overflow: scroll;
    }
    li {
      box-shadow: rgba(0, 0, 0, 0.19) 0px 0px 3px;
      margin: 5px 0px;
    }
  </style>
</head>
<body>
<nav class="navbar navbar-default navbar-fixed-top">
  <div class="container-fluid">
    <div class="navbar-header">
      <button type="button" class="navbar-toggle collapsed" data-toggle="collapse"
          data-target="#bs-example-navbar-collapse-1">
        <span class="sr-only">Toggle navigation</span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
      <a class="navbar-brand" href="#">{{.title}}</a>
    </div>
  </div>
</nav>
<div class="container-fluid" style="margin-top: 70px;margin-bottom: 20px;">
  <div class="container-fluid">
  <div class="menu col-md-4">
    <ul class="nav nav-pills nav-stacked" role="tablist">
      {{ range $key, $value := .apis }}
      <li role="presentation">
        <a href="#{{$key}}top" role="tab" data-toggle="tab">
          <span class="method {{$value.Method}}">{{$value.Method}}</span><span> <b>{{$value.URL}}</b></span>
        </a></li>
      {{ end }}
    </ul>
  </div>
  <div class="col-md-8 tab-content">
    {{ range $key, $value := .apis}}
    <div id="{{$key}}top"  role="tabpanel" class="tab-pane col-md-10">
      
      <p> <h4> Basic Info </h4> </p>
      <table class="table table-bordered table-striped">
        <tr>
          <th>Key</th>
          <th>Value</th>
        </tr>
        <tr>
          <td>Handler Name</td>
          <td> {{ $value.HandlerName }}</td>
        </tr>
        <tr>
          <td>Method</td>
          <td> {{ $value.Method }}</td>
        </tr>
        <tr>
          <td>Route</td>
          <td> {{ $value.URL }}</td>
        </tr>
        <tr>
          <td>Return Type</td>
          <td> {{ $value.ReturnType }}</td>
        </tr>
      </table>
      
      <p> <h4> Request Info </h4> </p>
      {{ if .Request }}
        <table class="table table-bordered table-striped">
          <tr>
            <th>Key</th>
            <th>Value</th>
          </tr>
          <tr>
            <td>Request Name</td>
            <td> {{ .Request.Name }} </td>
          </tr>
          <tr>
            <td>Package Path</td>
            <td> {{ .Request.PkgPath }} </td>
          </tr>
          <tr>
            <td>Example Url</td>
            <td> <pre>{{ .Request.CurlString }}</pre></td>
          </tr>
        </table>

        <p> <h5> Request Fields </h5> </p>

        <table class="table table-bordered table-striped">
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Tag</th>
            <th>Required</th>
          </tr>
          {{ range $key, $value := .Request.FieldInfos }}
          <tr>
            <td>{{ $value.Name }}</td>
            <td>{{ $value.Typ }}</td>
            <td>{{ $value.Tag }}</td>
            <td>{{ $value.Required }}</td>
          </tr>
          {{ end }}
        </table>
      {{ else }}
        <p> Not Defined </p>
      {{ end }}
      
      <p> <h4> Response Fields </h4> </p>
      {{ if .Response }}
        <p>Description: {{ .Response.Desc }} </p>
        {{ if .Response.FieldInfos }}
          <table class="table table-bordered table-striped">
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Tag</th>
            </tr>
            {{ range $key, $value := .Response.FieldInfos }}
            <tr>
              <td>{{ $value.Name }}</td>
              <td>{{ $value.Typ }}</td>
              <td>{{ $value.Tag }}</td>
            </tr>
            {{ end }}
          </table>
        {{ end }}
      {{ else }}
        <p> Not Defined </p>
      {{ end }}
      <hr>
    </div>
  {{ end }}
  </div>
  </div>
</div>
<hr>
</body>
</html>`
