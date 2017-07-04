<html>
    <head>
    <title></title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
    </head>
    <body>
    <div class="container">
      <h3>Listing database entries.</h3>
      <table class="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Message</th>
          </tr>
        </thead>
        <tbody>
            {{ range . }}
                <tr>
                    <td>"{{ .Name.Value }}"</td>
                    <td>"{{ .Email.Value }}"</td>
                    <td>"{{ .Message.Value }}"</td>
                </tr>
            {{ end }}
          </tr>
        </tbody>
      </table>
    </div>
    </body>
</html>
