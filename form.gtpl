<html>
    <head>
    <title></title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
    </head>
    <body>
    <div class="container">
    <h3>Add some data to be saved to the database.</h3>
      <div class="row">
        <div class="col-sm-4">
            <div>
                {{ if .Errors }}
                    <h4>The following fields are invalid or required:</h4>
                    {{ range .Errors }}
                        <li>{{ . }}</li>
                    {{ end }}
                    <hr/>
                {{ end }}
            </div>
            <form method="POST" action="/save">
              <div class="form-group">
                <label for="pwd">Name:</label>
                <input type="text" class="form-control" name="name" id="name" value="{{ .Name }}" required>
              </div>
              <div class="form-group">
                <label for="email">Email address:</label>
                <input type="email" class="form-control" name="email" id="email" value="{{ .Email }}" required>
              </div>
              <div class="form-group">
                <label for="message">Message:</label>
                <input type="text" class="form-control" name="message" id="message" value="{{ .Message }}">
              </div>
              <button type="submit" class="btn btn-default">Submit</button>
            </form>
        </div>
      </div>
    </div>
    </body>
</html>
