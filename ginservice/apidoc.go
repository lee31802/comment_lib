package ginservice

const apiDoc = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/default.min.css">
    <style>
       .api-content {
            display: none;
        }

       .title {
            color: #374151;
            font-weight: 700;
            font-size: 2.25rem;
            line-height: 2.5rem;
            margin-bottom: 1.5rem;
        }

       .subtitle {
            color: #6B7280;
            font-weight: 600;
            font-size: 1.5rem;
            line-height: 2rem;
            margin-bottom: 1rem;
        }

       .section-title {
            color: #374151;
            font-weight: 600;
            font-size: 1.25rem;
            line-height: 1.75rem;
            margin-bottom: 0.75rem;
        }

       .table {
            width: 100%;
            border-collapse: collapse;
        }

       .table th,
       .table td {
            border: 1px solid #E5E7EB;
            padding: 0.5rem 1rem;
        }

       .table th {
            background-color: #F3F4F6;
            text-align: left;
            font-weight: 600;
        }

       .json-sample {
            background-color: #F9FAFB;
            border-radius: 0.375rem;
            padding: 1rem;
            overflow-x: auto;
        }

       .fade-in {
            animation: fadeIn 0.5s ease-in-out;
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
            }
            to {
                opacity: 1;
            }
        }
    </style>
</head>

<body class="bg-gray-100 font-sans">
    <div class="container mx-auto p-8">
        <h1 class="title">{{.title}}</h1>
        <h2 class="subtitle">目录</h2>
        <ul class="list-disc pl-6 mb-6">
            {{range $index, $api :=.apis}}
            <li><a href="#" onclick="showAPIContent({{$index}}); return false;"
                    class="text-blue-600 hover:underline">{{$api.Method}} {{$api.URL}}</a></li>
            {{end}}
        </ul>
        {{range $index, $api :=.apis}}
        <div id="api-{{$index}}" class="api-content bg-white p-6 rounded shadow mb-6 fade-in">
            <h2 class="section-title">{{$api.Method}} {{$api.URL}}</h2>
            <p class="text-gray-600 mb-4">Handler: {{$api.HandlerName}}</p>
            {{if $api.Request}}
            <h3 class="section-title">Request</h3>
            <p class="text-gray-600 mb-2">Name: {{$api.Request.Name}}</p>
            <p class="text-gray-600 mb-2">PkgPath: {{$api.Request.PkgPath}}</p>
            <p class="text-gray-600 mb-4">Curl: <code>{{$api.Request.CurlString}}</code></p>
            <table class="table mb-6">
                <thead>
                    <tr>
                        <th class="text-left">Field Name</th>
                        <th class="text-left">Type</th>
                        <th class="text-left">Required</th>
                        <th class="text-left">Tag</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $api.Request.FieldInfos}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.Typ}}</td>
                        <td>{{.Required}}</td>
                        <td>{{.Tag}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <h4 class="section-title">Request JSON 样例</h4>
            <pre class="json-sample" id="request-json-sample-{{$index}}"></pre>
            <script>
                const curlString = '{{$api.Request.CurlString}}';
                const jsonStart = curlString.indexOf('{');
                const jsonEnd = curlString.lastIndexOf('}');
                if (jsonStart!== -1 && jsonEnd!== -1) {
                    const json = curlString.substring(jsonStart, jsonEnd + 1);
                    try {
                        const parsedJson = JSON.parse(json);
                        const formattedJson = JSON.stringify(parsedJson, null, 4);
                        document.getElementById('request-json-sample-{{$index}}').textContent = formattedJson;
                        hljs.highlightElement(document.getElementById('request-json-sample-{{$index}}'));
                    } catch (error) {
                        document.getElementById('request-json-sample-{{$index}}').textContent = json;
                        hljs.highlightElement(document.getElementById('request-json-sample-{{$index}}'));
                    }
                }
            </script>
            {{else}}
            <p class="text-gray-600 mb-4">无请求参数</p>
            {{end}}
            {{if $api.Response}}
            <h3 class="section-title">Response</h3>
            <table class="table mb-6">
                <thead>
                    <tr>
                        <th class="text-left">Field Name</th>
                        <th class="text-left">Type</th>
                        <th class="text-left">Tag</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $api.Response.FieldInfos}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.Typ}}</td>
                        <td>{{.Tag}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <h4 class="section-title">Response JSON 样例</h4>
            <pre class="json-sample" id="response-json-sample-{{$index}}"></pre>
            
            {{else}}
            <p class="text-gray-600 mb-4">无响应参数</p>
            {{end}}
        </div>
        {{end}}
    </div>
    <script>
        function showAPIContent(index) {
            const apiContents = document.querySelectorAll('.api-content');
            apiContents.forEach((content) => {
                content.style.display = 'none';
            });
            const targetContent = document.getElementById('api-' + index);
            if (targetContent) {
                targetContent.style.display = 'block';
                targetContent.classList.add('fade-in');
            }
        }
    </script>
</body>

</html>
`
